package storage

import (
	"github.com/rombintu/GophKeeper/internal/models/auth"
	"github.com/rombintu/GophKeeper/internal/models/storage"
)

type Storage interface {
	UserGet(user auth.User) (auth.User, error)
	UserCreate(user auth.User) error

	SecretGet(userID int64) (storage.Secret, error)
	SecretCreate(secret storage.Secret) error
	SecretsGet(userID int64) ([]storage.Secret, error)
}

func NewStorage(temp bool) Storage {
	if temp {
		return NewTempStorage()
	}
	return nil
}
