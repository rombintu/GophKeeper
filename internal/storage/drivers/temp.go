package drivers

import (
	"github.com/rombintu/GophKeeper/internal/models/auth"
	models "github.com/rombintu/GophKeeper/internal/models/storage"
)

type MemoryDriver struct {
	Users   []auth.User
	Secrets []models.Secret
}

func (md *MemoryDriver) UserGet(user auth.User) (auth.User, error) {
	return auth.User{}, nil
}

func (md *MemoryDriver) UserCreate(user auth.User) error {
	return nil
}

func (md *MemoryDriver) SecretGet(userID int64) (models.Secret, error) {
	return models.Secret{}, nil
}

func (md *MemoryDriver) SecretCreate(secret models.Secret) error {
	return nil
}

func (md *MemoryDriver) SecretsGet(userID int64) ([]models.Secret, error) {
	return []models.Secret{}, nil
}

func (md *MemoryDriver) Ping() error {
	return nil
}
