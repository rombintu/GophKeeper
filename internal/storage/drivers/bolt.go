package drivers

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"log/slog"

	kpb "github.com/rombintu/GophKeeper/internal/proto/keeper"
	"github.com/rombintu/GophKeeper/lib/crypto"
	bolt "go.etcd.io/bbolt"
)

const (
	profileTable = "profile"
	secretsTable = "secrets"
)

type BoltDriver struct {
	driver *bolt.DB
	path   string
}

func NewBoltDriver(path string) *BoltDriver {
	return &BoltDriver{
		path: path,
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
		encoder := gob.NewEncoder(&buf)
		if err := encoder.Encode(secret); err != nil {
			return fmt.Errorf("encode failed: %s", err.Error())
		}
		hash := crypto.GetHash(buf.Bytes())
		if err := b.Put(
			[]byte(fmt.Sprintf("%s:::%s", secret.UserEmail, hash)),
			buf.Bytes()); err != nil {
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
				slog.Debug("secret founded", slog.String("secret", string(v)))
				secretsEncoded = append(secretsEncoded, v)
			}
			return nil
		})

		var buf bytes.Buffer
		for _, s := range secretsEncoded {
			buf.Reset()
			encoder := gob.NewDecoder(&buf)
			var secret *kpb.Secret
			if err := encoder.Decode(&s); err != nil {
				return fmt.Errorf("decode failed: %s", err.Error())
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
