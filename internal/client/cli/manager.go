package cli

import (
	"context"
	"log/slog"

	"github.com/rombintu/GophKeeper/internal/client/models"
	apb "github.com/rombintu/GophKeeper/internal/proto/auth"
	kpb "github.com/rombintu/GophKeeper/internal/proto/keeper"
	spb "github.com/rombintu/GophKeeper/internal/proto/sync"
	"github.com/rombintu/GophKeeper/internal/storage"
	"github.com/rombintu/GophKeeper/lib/crypto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
	if len(secrets) == 0 {
		slog.Info("secret list is empty")
		return nil
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
	payload := secret.Payload()
	newSecret := &kpb.Secret{
		Title:       secret.GetTitle(),
		SecretType:  secret.GetType(),
		UserEmail:   m.profile.user.GetEmail(),
		HashPayload: crypto.GetHash(payload),
		Payload:     payload,
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

func (m *Manager) ConfigGetMap(ctx context.Context) (map[string]string, error) {
	return m.store.GetMap(ctx)

}

func (m *Manager) Login(ctx context.Context, serviceAddr string) error {
	conn, err := grpc.NewClient(serviceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	defer func() {
		if err := conn.Close(); err != nil {
			slog.Error("failed close connection", slog.String("error", err.Error()))
		}
	}()
	authClient := apb.NewAuthClient(conn)
	resp, err := authClient.Login(ctx, &apb.LoginRequest{User: m.profile.user})
	if err != nil {
		return err
	}
	slog.Debug("saved", slog.String("token", resp.GetToken()))
	return m.ConfigSet(ctx, map[string]string{
		"token": resp.GetToken(),
	})
}

func (m *Manager) Register(ctx context.Context, serviceAddr string) error {
	conn, err := grpc.NewClient(serviceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	defer func() {
		if err := conn.Close(); err != nil {
			slog.Error("failed close connection", slog.String("error", err.Error()))
		}
	}()
	authClient := apb.NewAuthClient(conn)
	resp, err := authClient.Register(ctx, &apb.RegisterRequest{User: m.profile.user})
	if err != nil {
		return err
	}
	slog.Debug("saved", slog.String("token", resp.GetToken()))
	return m.ConfigSet(ctx, map[string]string{
		"token": resp.GetToken(),
	})
}

func (m *Manager) Sync(ctx context.Context, serviceAddr string) error {
	conn, err := grpc.NewClient(serviceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	defer func() {
		if err := conn.Close(); err != nil {
			slog.Error("failed close connection", slog.String("error", err.Error()))
		}
	}()

	secretsForSync, err := m.store.SecretGetBatch(ctx)
	if err != nil {
		slog.Error("failed push secrets", slog.String("error", err.Error()))
		return err
	}

	syncClient := spb.NewSyncClient(conn)
	resp, err := syncClient.Process(
		ctx, &spb.SyncRequest{Email: m.profile.user.GetEmail(), Secrets: secretsForSync})
	if err != nil {
		return err
	}

	if err := m.store.SecretCreateBatch(ctx, resp.GetSecrets()); err != nil {
		slog.Error("failed pull secrets", slog.String("error", err.Error()))
		return err
	}

	slog.Debug("sync", slog.String("status", "OK"))
	return nil
}
