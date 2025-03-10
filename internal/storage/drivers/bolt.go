package drivers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/ProtonMail/go-crypto/openpgp"
	"github.com/google/uuid"
	kpb "github.com/rombintu/GophKeeper/internal/proto/keeper"
	"github.com/rombintu/GophKeeper/lib/crypto"
	bolt "go.etcd.io/bbolt"
)

const (
	secretsTable = "secrets"
	metaTable    = "meta"
)

type SecretMeta struct {
	Title      string
	SecretType kpb.Secret_SecretType
	Version    int64
	CreatedAt  int64
}

type SecretData struct {
	Payload []byte
}

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
		for _, table := range []string{metaTable, secretsTable} {
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
	tx, err := bd.driver.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	dataBucket := tx.Bucket([]byte(secretsTable))
	metaBucket := tx.Bucket([]byte(metaTable))

	meta := SecretMeta{
		Title:      secret.GetTitle(),
		SecretType: secret.GetSecretType(),
		Version:    secret.GetVersion(),
		CreatedAt:  secret.GetCreatedAt(),
	}

	data := SecretData{
		Payload: secret.GetPayload(),
	}

	metaBytes, err := json.Marshal(meta)
	if err != nil {
		return err
	}
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	hash := crypto.GetHash(dataBytes)
	uuid := uuid.New().String()
	key := fmt.Sprintf("%s:::%s:::%s", secret.GetUserEmail(), hash, uuid)

	dataBytesEncrypt, err := crypto.Encrypt(bd.cryptoKey, dataBytes)
	if err != nil {
		return fmt.Errorf("encrypt data failed: %s", err.Error())
	}

	if err := metaBucket.Put([]byte(key), metaBytes); err != nil {
		return err
	}
	if err := dataBucket.Put([]byte(key), dataBytesEncrypt); err != nil {
		return err
	}

	return tx.Commit()

}

// UNUSED FUNCTION
func (bd *BoltDriver) SecretCreateBatch(ctx context.Context, secrets []*kpb.Secret) error {
	return nil
}

func (bd *BoltDriver) SecretList(ctx context.Context, userEmail string) ([]*kpb.Secret, error) {
	var secrets []*kpb.Secret
	if userEmail == "" {
		return nil, errors.New("user email is empty")
	}
	if bd.cryptoKey == nil {
		return nil, errors.New("crypto key is not set")
	}
	if err := bd.driver.View(func(tx *bolt.Tx) error {
		dataBucket := tx.Bucket([]byte(secretsTable))
		metaBucket := tx.Bucket([]byte(metaTable))

		cursor := metaBucket.Cursor()
		prefix := []byte(userEmail + ":::")

		for k, v := cursor.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, v = cursor.Next() {
			// Декодирование метаданных
			var meta SecretMeta
			if err := json.Unmarshal(v, &meta); err != nil {
				return err
			}

			// Получение данных по тому же ключу
			dataEnc := dataBucket.Get(k)
			if dataEnc == nil {
				return fmt.Errorf("not found secret for key %s", k)
			}

			dataBytes, err := crypto.Decrypt(bd.cryptoKey, dataEnc)
			if err != nil {
				return err
			}
			var data SecretData
			if err := json.Unmarshal(dataBytes, &data); err != nil {
				return err
			}

			secrets = append(secrets, &kpb.Secret{
				Title:      meta.Title,
				SecretType: meta.SecretType,
				UserEmail:  userEmail,
				CreatedAt:  meta.CreatedAt,
				Version:    meta.Version,
				Payload:    data.Payload,
			})
		}

		return nil
	}); err != nil {
		return nil, err
	}
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
