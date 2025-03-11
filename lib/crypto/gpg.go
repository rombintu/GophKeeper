package crypto

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/ProtonMail/go-crypto/openpgp"
	"github.com/ProtonMail/go-crypto/openpgp/armor"
	"github.com/ProtonMail/go-crypto/openpgp/packet"
	proto "github.com/rombintu/GophKeeper/internal/proto/auth"
)

type Key openpgp.EntityList

// LoadPublicKey загружает открытый ключ из файла
func LoadPublicKey(filename string) (openpgp.EntityList, error) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return nil, fmt.Errorf("файл %s не найден", filename)
	}
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("ошибка при открытии файла открытого ключа: %w", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			slog.Warn("failed close file", slog.String("error", err.Error()))
		}
	}()

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

	defer func() {
		if err := file.Close(); err != nil {
			slog.Warn("failed close file", slog.String("error", err.Error()))
		}
	}()

	entityList, err := openpgp.ReadArmoredKeyRing(file)
	if err != nil {
		return nil, fmt.Errorf("ошибка при чтении закрытого ключа: %w", err)
	}

	return entityList, nil
}

// Encrypt шифрует данные с использованием открытого ключа
func Encrypt(key openpgp.EntityList, message []byte) ([]byte, error) {
	buf := new(bytes.Buffer)

	// Создание armor-обертки
	armorWriter, err := armor.Encode(buf, "PGP MESSAGE", nil)
	if err != nil {
		return nil, err
	}

	encryptedWriter, err := openpgp.Encrypt(armorWriter, key, nil, nil, &packet.Config{})
	if err != nil {
		return nil, fmt.Errorf("ошибка при создании шифровальщика: %w", err)
	}

	_, err = encryptedWriter.Write(message)
	if err != nil {
		return nil, fmt.Errorf("ошибка при записи данных для шифрования: %w", err)
	}
	if err := encryptedWriter.Close(); err != nil {
		slog.Error("failed close writer", slog.String("error", err.Error()))
		return nil, err
	}
	if err := armorWriter.Close(); err != nil {
		slog.Error("failed close writer", slog.String("error", err.Error()))
		return nil, err
	}

	return buf.Bytes(), nil
}

// Decrypt расшифровывает данные с использованием закрытого ключа
func Decrypt(key openpgp.EntityList, encryptedMessage []byte) ([]byte, error) {
	decbuf := bytes.NewBuffer(encryptedMessage)
	block, err := armor.Decode(decbuf)
	if err != nil {
		return nil, fmt.Errorf("ошибка при декодировании armored данных: %w", err)
	}

	md, err := openpgp.ReadMessage(block.Body, key, nil, &packet.Config{})
	if err != nil {
		return nil, fmt.Errorf("ошибка при чтении зашифрованного сообщения: %w", err)
	}

	decryptedData, err := io.ReadAll(md.UnverifiedBody)
	if err != nil {
		return nil, fmt.Errorf("ошибка при чтении расшифрованных данных: %w", err)
	}

	return decryptedData, nil
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
		user.KeyChecksum = entity.PrimaryKey.Fingerprint
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
