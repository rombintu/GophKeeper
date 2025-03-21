package crypto_test

import (
	"crypto/rand"
	"testing"

	"github.com/rombintu/GophKeeper/lib/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestStruct используется для тестирования сериализации и десериализации
type TestStruct struct {
	Name  string
	Value int
	Data  []byte
}

func TestEncodeDecode(t *testing.T) {
	// Генерация случайного ключа
	key := make([]byte, 32)
	_, err := rand.Read(key)
	require.NoError(t, err)

	// Тестовые данные
	testData := TestStruct{
		Name:  "Test",
		Value: 42,
		Data:  []byte("random data"),
	}

	// Шифрование данных
	ciphertext, err := crypto.Encode(key, testData)
	require.NoError(t, err)
	assert.NotEmpty(t, ciphertext, "Ciphertext should not be empty")

	// Дешифрование данных
	var decodedData TestStruct
	err = crypto.Decode(key, ciphertext, &decodedData)
	require.NoError(t, err)

	// Проверка, что данные совпадают
	assert.Equal(t, testData.Name, decodedData.Name, "Decoded name should match original")
	assert.Equal(t, testData.Value, decodedData.Value, "Decoded value should match original")
	assert.Equal(t, testData.Data, decodedData.Data, "Decoded data should match original")
}

func TestEncodeInvalidKeySize(t *testing.T) {
	// Неправильный размер ключа
	key := make([]byte, 15)
	testData := TestStruct{Name: "Test"}

	_, err := crypto.Encode(key, testData)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid key size")
}

func TestDecodeInvalidKeySize(t *testing.T) {
	// Неправильный размер ключа
	key := make([]byte, 15)
	ciphertext := make([]byte, 32)
	var decodedData TestStruct

	err := crypto.Decode(key, ciphertext, &decodedData)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid key size")
}

func TestDecodeInvalidCiphertext(t *testing.T) {
	// Генерация случайного ключа
	key := make([]byte, 32)
	_, err := rand.Read(key)
	require.NoError(t, err)

	// Некорректный ciphertext (слишком короткий)
	ciphertext := make([]byte, 10)
	var decodedData TestStruct

	err = crypto.Decode(key, ciphertext, &decodedData)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "ciphertext too short")
}

func TestDecodeTamperedCiphertext(t *testing.T) {
	// Генерация случайного ключа
	key := make([]byte, 32)
	_, err := rand.Read(key)
	require.NoError(t, err)

	// Тестовые данные
	testData := TestStruct{Name: "Test"}

	// Шифрование данных
	ciphertext, err := crypto.Encode(key, testData)
	require.NoError(t, err)

	// Изменение ciphertext (подмена данных)
	ciphertext[10] ^= 0xFF

	// Попытка дешифровать измененные данные
	var decodedData TestStruct
	err = crypto.Decode(key, ciphertext, &decodedData)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "decryption error")
}

func TestEncodeDecodeEmptyData(t *testing.T) {
	// Генерация случайного ключа
	key := make([]byte, 32)
	_, err := rand.Read(key)
	require.NoError(t, err)

	// Пустые данные
	testData := TestStruct{}

	// Шифрование данных
	ciphertext, err := crypto.Encode(key, testData)
	require.NoError(t, err)

	// Дешифрование данных
	var decodedData TestStruct
	err = crypto.Decode(key, ciphertext, &decodedData)
	require.NoError(t, err)

	// Проверка, что данные совпадают
	assert.Equal(t, testData, decodedData, "Decoded data should match original")
}

func TestEncodeDecodeLargeData(t *testing.T) {
	// Генерация случайного ключа
	key := make([]byte, 32)
	_, err := rand.Read(key)
	require.NoError(t, err)

	// Большие данные
	largeData := make([]byte, 1024*1024) // 1 MB
	_, err = rand.Read(largeData)
	require.NoError(t, err)

	testData := TestStruct{
		Name:  "Large Data",
		Value: 100,
		Data:  largeData,
	}

	// Шифрование данных
	ciphertext, err := crypto.Encode(key, testData)
	require.NoError(t, err)

	// Дешифрование данных
	var decodedData TestStruct
	err = crypto.Decode(key, ciphertext, &decodedData)
	require.NoError(t, err)

	// Проверка, что данные совпадают
	assert.Equal(t, testData.Name, decodedData.Name, "Decoded name should match original")
	assert.Equal(t, testData.Value, decodedData.Value, "Decoded value should match original")
	assert.Equal(t, testData.Data, decodedData.Data, "Decoded data should match original")
}

func TestGobSerializationError(t *testing.T) {
	// Генерация случайного ключа
	key := make([]byte, 32)
	_, err := rand.Read(key)
	require.NoError(t, err)

	// Данные, которые не могут быть сериализованы (chan)
	invalidData := make(chan int)

	// Попытка шифрования
	_, err = crypto.Encode(key, invalidData)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gob encode error")
}

func TestGobDeserializationError(t *testing.T) {
	// Генерация случайного ключа
	key := make([]byte, 32)
	_, err := rand.Read(key)
	require.NoError(t, err)

	// Тестовые данные
	testData := TestStruct{Name: "Test"}

	// Шифрование данных
	ciphertext, err := crypto.Encode(key, testData)
	require.NoError(t, err)

	// Попытка дешифровать в неправильную структуру
	var wrongStruct struct {
		WrongField string
	}
	err = crypto.Decode(key, ciphertext, &wrongStruct)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gob decode error")
}
