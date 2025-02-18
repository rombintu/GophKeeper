package crypto

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/ProtonMail/go-crypto/openpgp"
	"github.com/ProtonMail/go-crypto/openpgp/armor"
	"github.com/ProtonMail/go-crypto/openpgp/packet"
	"github.com/rombintu/GophKeeper/internal/proto"
)

// LoadPublicKey загружает открытый ключ из файла
func LoadPublicKey(filename string) (openpgp.EntityList, error) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return nil, fmt.Errorf("файл %s не найден", filename)
	}
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("ошибка при открытии файла открытого ключа: %w", err)
	}
	defer file.Close()

	entityList, err := openpgp.ReadArmoredKeyRing(file)
	if err != nil {
		return nil, fmt.Errorf("ошибка при чтении открытого ключа: %w", err)
	}

	return entityList, nil
}

// LoadPrivateKey загружает закрытый ключ из файла
func LoadPrivateKey(filename string) (openpgp.EntityList, error) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return nil, fmt.Errorf("файл %s не найден. Укажите правильный путь до файла", filename)
	}
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("ошибка при открытии файла закрытого ключа: %w", err)
	}
	defer file.Close()

	entityList, err := openpgp.ReadArmoredKeyRing(file)
	if err != nil {
		return nil, fmt.Errorf("ошибка при чтении закрытого ключа: %w", err)
	}

	return entityList, nil
}

// Encrypt шифрует данные с использованием открытого ключа
func Encrypt(key openpgp.EntityList, message string) (string, error) {
	buf := new(bytes.Buffer)
	encryptedWriter, err := openpgp.Encrypt(buf, key, nil, nil, &packet.Config{})
	if err != nil {
		return "", fmt.Errorf("ошибка при создании шифровальщика: %w", err)
	}

	_, err = encryptedWriter.Write([]byte(message))
	if err != nil {
		return "", fmt.Errorf("ошибка при записи данных для шифрования: %w", err)
	}
	encryptedWriter.Close()

	return buf.String(), nil
}

// Decrypt расшифровывает данные с использованием закрытого ключа
func Decrypt(privateKey openpgp.EntityList, encryptedMessage string) (string, error) {
	decbuf := bytes.NewBuffer([]byte(encryptedMessage))
	block, err := armor.Decode(decbuf)
	if err != nil {
		return "", fmt.Errorf("ошибка при декодировании armored данных: %w", err)
	}

	md, err := openpgp.ReadMessage(block.Body, privateKey, nil, &packet.Config{})
	if err != nil {
		return "", fmt.Errorf("ошибка при чтении зашифрованного сообщения: %w", err)
	}

	decryptedData, err := io.ReadAll(md.UnverifiedBody)
	if err != nil {
		return "", fmt.Errorf("ошибка при чтении расшифрованных данных: %w", err)
	}

	return string(decryptedData), nil
}

// GetKeyHash возвращает отпечаток (хеш) ключа
func GetKeyHash(entity *openpgp.Entity) []byte {
	if entity.PrimaryKey == nil {
		return nil
	}
	return entity.PrimaryKey.Fingerprint
}

func GetProfile(privateKey openpgp.EntityList) (*proto.User, error) {
	user := &proto.User{}
	// Проходим по списку ключей и извлекаем profile
	for _, entity := range privateKey {
		user.HexKeys = entity.PrimaryKey.Fingerprint
		for _, identity := range entity.Identities {
			if identity.UserId.Email == "" {
				continue
			}
			user.Email = identity.UserId.Email
			return user, nil
		}
	}
	return nil, fmt.Errorf("у профиля не задан email")
}
