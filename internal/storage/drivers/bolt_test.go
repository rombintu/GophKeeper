package drivers

import (
	"bytes"
	"context"
	"crypto/rand"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"github.com/ProtonMail/go-crypto/openpgp"
	"github.com/rombintu/GophKeeper/internal/proto/keeper"
	"github.com/rombintu/GophKeeper/lib/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	bolt "go.etcd.io/bbolt"
)

const (
	testPrivateKeyPath = "../../../profiles/test.key"
	testUserEmail      = "dev@dev.com"
)

// Генерация тестовых ключей (выполнить один раз)
// go run github.com/ProtonMail/gopenpgp/v2/cmd/generate-key -n "Test Key" -e test@example.com -p 1234 -o profiles/test.key
// gpg --export test@example.com > profiles/test.key

func setupBoltDriver(t *testing.T) (*BoltDriver, string) {
	// Загрузка тестового ключа
	privKey, err := crypto.LoadPrivateKey(testPrivateKeyPath)
	require.NoError(t, err)

	// Создание временной БД
	dbPath := filepath.Join(t.TempDir(), "test.db")
	bd := NewBoltDriver(dbPath, privKey)

	// Инициализация драйвера
	require.NoError(t, bd.Open(context.Background()))
	t.Cleanup(func() {
		_ = bd.Close(context.Background()) //nolint:errcheck
	})
	require.NoError(t, bd.Configure(context.Background()))

	return bd, dbPath
}

func TestSecretCreateAndRetrieve(t *testing.T) {
	bd, _ := setupBoltDriver(t)

	// Создание тестового секрета
	secret := &keeper.Secret{
		Title:       "My Password",
		UserEmail:   testUserEmail,
		SecretType:  keeper.Secret_CRED,
		Payload:     []byte("secret_password_123"),
		HashPayload: "hash123",
		Version:     1,
		CreatedAt:   time.Now().Unix(),
	}

	// Сохранение секрета
	require.NoError(t, bd.SecretCreate(context.Background(), secret))

	// Получение списка секретов
	secrets, err := bd.SecretList(context.Background(), testUserEmail)
	require.NoError(t, err)
	require.Len(t, secrets, 1)

	// Проверка корректности данных
	retrieved := secrets[0]
	assert.Equal(t, secret.Title, retrieved.Title)
	assert.Equal(t, secret.SecretType, retrieved.SecretType)
	assert.Equal(t, secret.Payload, retrieved.Payload)
	assert.Equal(t, secret.HashPayload, retrieved.HashPayload)
}

func TestDataEncryption(t *testing.T) {
	bd, _ := setupBoltDriver(t)

	// Генерация случайных данных
	originalData := make([]byte, 128)
	_, _ = rand.Read(originalData) //nolint:errcheck

	secret := &keeper.Secret{
		UserEmail:   testUserEmail,
		Payload:     originalData,
		HashPayload: "random_data",
	}

	// Сохранение секрета
	require.NoError(t, bd.SecretCreate(context.Background(), secret))

	// Проверка зашифрованных данных в БД
	err := bd.driver.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(secretsTable))
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			// Проверка, что данные не совпадают с оригиналом
			assert.False(t, bytes.Equal(v, originalData),
				"Data should be encrypted")
		}
		return nil
	})
	require.NoError(t, err)
}

func TestSecretListForDifferentUsers(t *testing.T) {
	bd, _ := setupBoltDriver(t)

	// Секреты для разных пользователей
	secrets := []*keeper.Secret{
		{UserEmail: "user1@test.com", Payload: []byte("data1")},
		{UserEmail: "user2@test.com", Payload: []byte("data2")},
		{UserEmail: "user1@test.com", Payload: []byte("data3")},
	}

	for _, s := range secrets {
		require.NoError(t, bd.SecretCreate(context.Background(), s))
	}

	// Получение секретов для user1
	user1Secrets, err := bd.SecretList(context.Background(), "user1@test.com")
	require.NoError(t, err)
	assert.Len(t, user1Secrets, 2)

	// Получение секретов для несуществующего пользователя
	emptySecrets, err := bd.SecretList(context.Background(), "nonexistent@test.com")
	require.NoError(t, err)
	assert.Empty(t, emptySecrets)
}

func TestSecretBatchOperations(t *testing.T) {
	bd, _ := setupBoltDriver(t)

	// Генерация тестовых данных
	var batch []*keeper.Secret
	for i := 0; i < 10; i++ {
		batch = append(batch, &keeper.Secret{
			UserEmail:   testUserEmail,
			Payload:     []byte(fmt.Sprintf("batch_data_%d", i)),
			HashPayload: fmt.Sprintf("hash_%d", i),
		})
	}

	// Пакетное сохранение
	require.NoError(t, bd.SecretCreateBatch(context.Background(), batch))

	// Проверка количества записей
	err := bd.driver.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(metaTable))
		assert.Equal(t, 10, b.Stats().KeyN)
		return nil
	})
	require.NoError(t, err)
}

func TestErrorHandling(t *testing.T) {
	// Тест без криптографического ключа
	t.Run("NoCryptoKey", func(t *testing.T) {
		bd := NewBoltDriver("test.db", nil)
		err := bd.SecretCreate(context.Background(), &keeper.Secret{})
		require.ErrorContains(t, err, "crypto key is not set")
	})

	// Тест с неверным ключом
	t.Run("InvalidCryptoKey", func(t *testing.T) {
		badKey, _ := openpgp.NewEntity("test", "", "test@bad.key", nil)
		bd := NewBoltDriver("test.db", openpgp.EntityList{badKey})
		err := bd.SecretCreate(context.Background(), &keeper.Secret{
			UserEmail: testUserEmail,
			Payload:   []byte("test"),
		})
		require.ErrorContains(t, err, "encrypt data failed")
	})
}

func TestGetSetOperations(t *testing.T) {
	bd, _ := setupBoltDriver(t)

	// Тест обычной записи
	key := []byte("test_key")
	value := []byte("test_value")

	require.NoError(t, bd.Set(context.Background(), key, value))

	// Тест чтения
	retrieved, err := bd.Get(context.Background(), key)
	require.NoError(t, err)
	assert.Equal(t, value, retrieved)

	// Тест несуществующего ключа
	_, err = bd.Get(context.Background(), []byte("bad_key"))
	require.ErrorContains(t, err, "not found")
}
