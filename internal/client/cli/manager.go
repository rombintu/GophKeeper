package cli

import (
	"context"
	"log/slog"

	"github.com/rombintu/GophKeeper/internal/client/models"
	kpb "github.com/rombintu/GophKeeper/internal/proto/keeper"
	"github.com/rombintu/GophKeeper/internal/storage"
)

type Manager struct {
	store   storage.ClientManager
	profile *Profile
}

func NewManager(profile *Profile, store storage.ClientManager) *Manager {
	return &Manager{
		store:   store,
		profile: profile,
	}
}

func (m *Manager) Configure() error {
	return m.store.Configure(context.Background())
}

func (m *Manager) SecretList(ctx context.Context) error {
	secrets, err := m.store.SecretList(ctx, m.profile.user.GetEmail())
	if err != nil {
		return err
	}
	for _, s := range secrets {
		slog.Info("secret",
			slog.String("title", s.GetTitle()),
			slog.String("type", s.GetSecretType().String()),
			slog.String("email", s.GetUserEmail()),
			slog.String("payload", string(s.GetPayload())),
		)
	}
	return nil
}

func (m *Manager) SecretCreate(ctx context.Context, secret models.SecretAdapter) error {
	newSecret := &kpb.Secret{
		Title:      secret.Title(),
		SecretType: secret.Type(),
		UserEmail:  m.profile.user.GetEmail(),
		Payload:    secret.Encode(),
	}
	if err := m.store.SecretCreate(ctx, newSecret); err != nil {
		return err
	}

	return nil
}
