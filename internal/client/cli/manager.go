package cli

import (
	"context"
	"fmt"

	"github.com/rombintu/GophKeeper/internal/client/models"
	kpb "github.com/rombintu/GophKeeper/internal/proto/keeper"
	"github.com/rombintu/GophKeeper/internal/storage"
)

type Manager struct {
	store   storage.ClientManager
	profile Profile
}

func NewManager(store storage.ClientManager) *Manager {
	return &Manager{
		store:   store,
		profile: Profile{},
	}
}

func (m *Manager) Configure() error {
	return m.store.Configure(context.Background())
}

func (m *Manager) SecretList(ctx context.Context) error {
	secrets, err := m.store.SecretList(ctx, m.profile.Email)
	if err != nil {
		return err
	}
	for _, s := range secrets {
		fmt.Printf("%+v \n", s)
	}
	return nil
}

func (m *Manager) SecretCreate(ctx context.Context, secret models.SecretAdapter) error {
	newSecret := &kpb.Secret{
		Title:      secret.Title(),
		SecretType: secret.Type(),
		UserEmail:  m.profile.Email,
		Payload:    secret.Encode(),
	}
	if err := m.store.SecretCreate(ctx, newSecret); err != nil {
		return err
	}

	return nil
}
