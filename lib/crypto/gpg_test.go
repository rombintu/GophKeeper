package crypto_test

import (
	"bytes"
	"crypto/rand"
	"os"
	"testing"

	"github.com/ProtonMail/go-crypto/openpgp"
	"github.com/ProtonMail/go-crypto/openpgp/armor"
	"github.com/rombintu/GophKeeper/lib/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testPrivateKeyPath = "../../profiles/test.key"
	testUserEmail      = "test@test.com"
)

// Генерация тестовых ключей (выполнить один раз)
// go run github.com/ProtonMail/gopenpgp/v2/cmd/generate-key -n "Test Key" -e test@example.com -p 1234 -o profiles/test.key
// gpg --export test@example.com > testdata/public.key

func setupTestKeys(t *testing.T) openpgp.EntityList {
	// Загрузка тестового ключа
	privKey, err := crypto.LoadPrivateKey(testPrivateKeyPath)
	require.NoError(t, err)
	return privKey
}

func TestLoadPrivateKey(t *testing.T) {
	t.Run("Valid key", func(t *testing.T) {
		privKey, err := crypto.LoadPrivateKey(testPrivateKeyPath)
		require.NoError(t, err)
		assert.NotEmpty(t, privKey, "Private key should be loaded")
	})

	t.Run("Invalid file path", func(t *testing.T) {
		_, err := crypto.LoadPrivateKey("nonexistent.key")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "не найден")
	})

	t.Run("Invalid key format", func(t *testing.T) {
		// Создаем файл с некорректным содержимым
		tmpFile, err := os.CreateTemp("", "invalid_key_*.key")
		require.NoError(t, err)
		defer func() {
			err := os.Remove(tmpFile.Name()) // nolint:errcheck
			if err != nil {
				t.Fatal(err)
			}
		}()

		_, err = tmpFile.Write([]byte("invalid key data"))
		require.NoError(t, err)
		if err := tmpFile.Close(); err != nil {
			t.Fatal(err)
		}

		_, err = crypto.LoadPrivateKey(tmpFile.Name())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "ошибка при чтении закрытого ключа")
	})
}

func TestEncryptDecrypt(t *testing.T) {
	privKey := setupTestKeys(t)

	// Генерация случайных данных
	originalData := make([]byte, 128)
	_, err := rand.Read(originalData)
	if err != nil {
		t.Fatal(err)
	}

	// Шифрование данных
	encryptedData, err := crypto.Encrypt(privKey, originalData)
	require.NoError(t, err)
	assert.NotEqual(t, originalData, encryptedData, "Encrypted data should not match original")

	// Дешифрование данных
	decryptedData, err := crypto.Decrypt(privKey, encryptedData)
	require.NoError(t, err)
	assert.Equal(t, originalData, decryptedData, "Decrypted data should match original")
}

func TestEncryptDecryptInvalidKey(t *testing.T) {
	// Генерация нового ключа, который не подходит для дешифровки
	invalidKey, err := openpgp.NewEntity("Test", "", "invalid@example.com", nil)
	require.NoError(t, err)

	// Шифрование данных с правильным ключом
	privKey := setupTestKeys(t)
	originalData := []byte("test data")
	encryptedData, err := crypto.Encrypt(privKey, originalData)
	require.NoError(t, err)

	// Попытка дешифровки с неправильным ключом
	_, err = crypto.Decrypt(openpgp.EntityList{invalidKey}, encryptedData)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "ошибка при чтении зашифрованного сообщения")
}

func TestGetKeyHash(t *testing.T) {
	privKey := setupTestKeys(t)

	// Получение хеша ключа
	hash := crypto.GetKeyHash(privKey[0])
	require.NotEmpty(t, hash, "Key hash should not be empty")

	// Проверка, что хеш соответствует ожидаемому формату (SHA-1 fingerprint)
	assert.Len(t, hash, 20, "Key hash should be 20 bytes long")
}

func TestGetUserFromKey(t *testing.T) {
	privKey := setupTestKeys(t)

	// Получение пользователя из ключа
	user, err := crypto.GetUserFromKey(privKey)
	require.NoError(t, err)
	assert.Equal(t, testUserEmail, user.Email, "User email should match key identity")

	// Проверка, что хеш ключа установлен
	assert.NotEmpty(t, user.KeyChecksum, "Key checksum should not be empty")
}

func TestGetUserFromKeyNoEmail(t *testing.T) {
	// Создание ключа без email
	entity, err := openpgp.NewEntity("Test", "", "", nil)
	require.NoError(t, err)

	// Попытка извлечь пользователя
	_, err = crypto.GetUserFromKey(openpgp.EntityList{entity})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "у профиля не задан email")
}

func TestArmorEncoding(t *testing.T) {
	// Проверка, что зашифрованные данные используют armor-кодировку
	privKey := setupTestKeys(t)
	originalData := []byte("test data")

	encryptedData, err := crypto.Encrypt(privKey, originalData)
	require.NoError(t, err)

	// Попытка декодировать armor-обертку
	block, err := armor.Decode(bytes.NewReader(encryptedData))
	require.NoError(t, err)
	assert.Equal(t, "PGP MESSAGE", block.Type, "Armor block type should be PGP MESSAGE")
}

func TestDecryptInvalidData(t *testing.T) {
	privKey := setupTestKeys(t)

	// Попытка дешифровать некорректные данные
	_, err := crypto.Decrypt(privKey, []byte("invalid data"))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "ошибка при декодировании armored данных")
}
