package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/gob"
	"fmt"
	"io"
)

// Encode структуры с шифрованием
func Encode(key []byte, data interface{}) ([]byte, error) {
	// Проверка длины ключа
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		return nil, fmt.Errorf("invalid key size: must be 16, 24 or 32 bytes")
	}

	// Сериализация данных
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(data); err != nil {
		return nil, fmt.Errorf("gob encode error: %w", err)
	}
	plaintext := buf.Bytes()

	// Создание cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("cipher creation error: %w", err)
	}

	// Создание GCM режима
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("GCM creation error: %w", err)
	}

	// Генерация nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("nonce generation error: %w", err)
	}

	// Шифрование данных
	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return ciphertext, nil
}

// Decode данных в структуру
func Decode(key []byte, ciphertext []byte, data interface{}) error {
	// Проверка длины ключа
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		return fmt.Errorf("invalid key size: must be 16, 24 or 32 bytes")
	}

	// Создание cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		return fmt.Errorf("cipher creation error: %w", err)
	}

	// Создание GCM режима
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("GCM creation error: %w", err)
	}

	// Извлечение nonce
	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	// Дешифровка
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return fmt.Errorf("decryption error: %w", err)
	}

	// Десериализация
	buf := bytes.NewBuffer(plaintext)
	dec := gob.NewDecoder(buf)
	if err := dec.Decode(data); err != nil {
		return fmt.Errorf("gob decode error: %w", err)
	}

	return nil
}
