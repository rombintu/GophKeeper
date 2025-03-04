package drivers

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"log/slog"

	"github.com/ProtonMail/go-crypto/openpgp"
	kpb "github.com/rombintu/GophKeeper/internal/proto/keeper"
	"github.com/rombintu/GophKeeper/lib/crypto"
	bolt "go.etcd.io/bbolt"
)

const (
	profileTable = "profile"
	secretsTable = "secrets"
)

type BoltDriver struct {
	cryptoKey openpgp.EntityList
	driver    *bolt.DB
	path      string
}

func NewBoltDriver(path string, cryptoKey openpgp.EntityList) *BoltDriver {
	return &BoltDriver{
		cryptoKey: cryptoKey,
		path:      path,
	}
}

func (bd *BoltDriver) Open(ctx context.Context) (err error) {
	bd.driver, err = bolt.Open(bd.path, 0600, nil)
	if err != nil {
		return err
	}
	return nil
}

func (bd *BoltDriver) Close(ctx context.Context) error {
	return bd.driver.Close()
}

func (bd *BoltDriver) Ping(ctx context.Context, monitoring bool) error {
	return nil
}

// Create tables
func (bd *BoltDriver) Configure(ctx context.Context) error {
	bd.driver.Update(func(tx *bolt.Tx) error {
		for _, table := range []string{profileTable, secretsTable} {
			_, err := tx.CreateBucketIfNotExists([]byte(table))
			if err != nil {
				return err
			}
		}
		return nil
	})
	return nil
}

func (bd *BoltDriver) SecretCreate(ctx context.Context, secret *kpb.Secret) error {
	bd.driver.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(secretsTable))
		if b == nil {
			return fmt.Errorf("bucket %s does not exist", secretsTable)
		}

		var buf bytes.Buffer
		// ENCODE DATA
		encoder := gob.NewEncoder(&buf)
		if err := encoder.Encode(secret); err != nil {
			return fmt.Errorf("encode failed: %s", err.Error())
		}

		var data []byte
		var keyset bool
		var err error
		// ENCRYPT DATA
		if bd.cryptoKey != nil {
			// Если установлен ключ, производим шифрование
			data, err = crypto.Encrypt(bd.cryptoKey, buf.Bytes())
			if err != nil {
				return fmt.Errorf("encrypt failed: %s", err.Error())
			}
			// Устанавливаем флаг, что данные зашифрованны
			keyset = true
		} else {
			data = buf.Bytes()
			slog.Warn("key has not been installed, the cryptography is skipped")
		}

		hash := crypto.GetHash(buf.Bytes())
		if err := b.Put(
			fmt.Appendf(nil, "%s:::%s:::%t", secret.UserEmail, hash, keyset),
			data); err != nil {
			return err
		}
		return nil
	})
	return nil
}

func (bd *BoltDriver) SecretList(ctx context.Context, userEmail string) ([]*kpb.Secret, error) {
	var secrets []*kpb.Secret
	bd.driver.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(secretsTable))
		if b == nil {
			return fmt.Errorf("bucket %s does not exist", secretsTable)
		}

		var secretsEncoded [][]byte

		b.ForEach(func(k, v []byte) error {
			if bytes.HasPrefix(k, []byte(userEmail)) {
				secretsEncoded = append(secretsEncoded, v)
			}
			return nil
		})

		var buf bytes.Buffer
		for _, encoded := range secretsEncoded {
			// Записываем закодированные данные в буфер
			if _, err := buf.Write(encoded); err != nil {
				return fmt.Errorf("buffer write failed: %w", err)
			}

			// Создаём декодер для буфера
			decoder := gob.NewDecoder(&buf)
			var secret *kpb.Secret

			// Декодируем в переменную secret
			if err := decoder.Decode(&secret); err != nil {
				return fmt.Errorf("decode failed: %w", err)
			}

			secrets = append(secrets, secret)
		}
		return nil
	})
	return secrets, nil
}

func (bd *BoltDriver) SecretPurge(ctx context.Context, secret *kpb.Secret) error {
	return nil
}

func (bd *BoltDriver) SaveProfile() error {
	return nil
}

func (bd *BoltDriver) LoadProfile() error {
	return nil
}
