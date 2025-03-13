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

func NewManager() *Manager {
	return &Manager{}
}

func (m *Manager) SetStore(store storage.ClientManager) {
	m.store = store
}

func (m *Manager) SetProfile(profile *Profile) {
	m.profile = profile
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
		Title:      secret.GetTitle(),
		SecretType: secret.GetType(),
		UserEmail:  m.profile.user.GetEmail(),
		Payload:    secret.Payload(),
	}
	if err := m.store.SecretCreate(ctx, newSecret); err != nil {
		return err
	}

	return nil
}

func (m *Manager) ConfigSet(ctx context.Context, values map[string]string) error {
	for k, v := range values {
		if err := m.store.Set(ctx, []byte(k), []byte(v)); err != nil {
			return err
		}
	}
	return nil
}

func (m *Manager) ConfigGet(ctx context.Context, key string) (string, error) {
	data, err := m.store.Get(ctx, []byte(key))
	if err != nil {
		return "", err
	}
	return string(data), nil
}
