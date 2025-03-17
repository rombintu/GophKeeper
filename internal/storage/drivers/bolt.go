package drivers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"

	"github.com/ProtonMail/go-crypto/openpgp"
	"github.com/google/uuid"
	kpb "github.com/rombintu/GophKeeper/internal/proto/keeper"
	"github.com/rombintu/GophKeeper/lib/crypto"
	bolt "go.etcd.io/bbolt"
)

const (
	profileTable = "profile"
	secretsTable = "secrets"
	metaTable    = "meta"
)

type SecretMeta struct {
	Title      string
	UserEmail  string
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
	return bd.driver.Update(func(tx *bolt.Tx) error {
		for _, table := range []string{metaTable, secretsTable, profileTable} {
			_, err := tx.CreateBucketIfNotExists([]byte(table))
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (bd *BoltDriver) SecretCreate(ctx context.Context, secret *kpb.Secret) error {
	tx, err := bd.driver.Begin(true)
	if err != nil {
		return err
	}
	defer func() {
		if err := tx.Rollback(); err != nil {
			slog.Error(err.Error())
		}
	}()

	dataBucket := tx.Bucket([]byte(secretsTable))
	metaBucket := tx.Bucket([]byte(metaTable))

	meta := SecretMeta{
		Title:      secret.GetTitle(),
		UserEmail:  secret.GetUserEmail(),
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
				slog.Warn("failed decrypt secret. skip...", slog.String("key", string(k)), slog.String("error", err.Error()))
				continue
			}
			var data SecretData
			if err := json.Unmarshal(dataBytes, &data); err != nil {
				return err
			}

			secrets = append(secrets, &kpb.Secret{
				Title:      meta.Title,
				SecretType: meta.SecretType,
				UserEmail:  meta.UserEmail,
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

func (bd *BoltDriver) Set(ctx context.Context, key, value []byte) error {
	return bd.driver.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(profileTable))
		if err := b.Put(key, value); err != nil {
			return err
		}
		return nil
	})
}

func (bd *BoltDriver) Get(ctx context.Context, key []byte) ([]byte, error) {
	var data []byte
	if err := bd.driver.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(profileTable))
		data = b.Get(key)
		if data == nil {
			return fmt.Errorf("not found %s", key)
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return data, nil
}

func (bd *BoltDriver) GetMap(ctx context.Context) (map[string]string, error) {
	data := make(map[string]string)
	if err := bd.driver.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(profileTable))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			data[string(k)] = string(v)
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return data, nil
}

// Дешифровки не происходит
func (bd *BoltDriver) SecretGetBatch(ctx context.Context) ([]*kpb.Secret, error) {
	var secrets []*kpb.Secret

	if bd.cryptoKey == nil {
		return nil, errors.New("crypto key is not set")
	}
	if err := bd.driver.View(func(tx *bolt.Tx) error {
		dataBucket := tx.Bucket([]byte(secretsTable))
		metaBucket := tx.Bucket([]byte(metaTable))

		cursor := metaBucket.Cursor()
		lastKey, _ := cursor.Last()

		for k, v := cursor.Seek(lastKey); k != nil; k, v = cursor.Prev() {
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

			secrets = append(secrets, &kpb.Secret{
				Title:      meta.Title,
				SecretType: meta.SecretType,
				UserEmail:  meta.UserEmail,
				CreatedAt:  meta.CreatedAt,
				Version:    meta.Version,
				Payload:    dataEnc,
			})
		}

		return nil
	}); err != nil {
		return nil, err
	}
	return secrets, nil
}

func (bd *BoltDriver) SecretCreateBatch(ctx context.Context, secrets []*kpb.Secret) error {
	tx, err := bd.driver.Begin(true)
	if err != nil {
		return err
	}
	defer func() {
		if err := tx.Rollback(); err != nil {
			slog.Error(err.Error())
		}
	}()

	dataBucket := tx.Bucket([]byte(secretsTable))
	metaBucket := tx.Bucket([]byte(metaTable))

	for _, s := range secrets {

		meta := SecretMeta{
			Title:      s.GetTitle(),
			UserEmail:  s.GetUserEmail(),
			SecretType: s.GetSecretType(),
			Version:    s.GetVersion(),
			CreatedAt:  s.GetCreatedAt(),
		}

		data := SecretData{
			Payload: s.GetPayload(),
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
		key := fmt.Sprintf("%s:::%s:::%s", s.GetUserEmail(), hash, uuid)

		if err := metaBucket.Put([]byte(key), metaBytes); err != nil {
			return err
		}
		if err := dataBucket.Put([]byte(key), dataBytes); err != nil {
			return err
		}

	}
	return tx.Commit()

}
