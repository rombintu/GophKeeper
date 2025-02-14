package storage

import (
	"github.com/rombintu/GophKeeper/internal/models/auth"
	"github.com/rombintu/GophKeeper/internal/models/storage"
)

type TempStorage struct {
	Users   []auth.User
	Secrets []storage.Secret
}

func NewTempStorage() Storage {
	return &TempStorage{}
}

func (ts *TempStorage) UserGet(user auth.User) (auth.User, error) {
	return auth.User{}, nil
}

func (ts *TempStorage) UserCreate(user auth.User) error {
	return nil
}

func (ts *TempStorage) SecretGet(userID int64) (storage.Secret, error) {
	return storage.Secret{}, nil
}

func (ts *TempStorage) SecretCreate(secret storage.Secret) error {
	return nil
}

func (ts *TempStorage) SecretsGet(userID int64) ([]storage.Secret, error) {
	return []storage.Secret{}, nil
}
