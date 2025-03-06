package drivers

import (
	"context"
	"errors"
	"log/slog"

	apb "github.com/rombintu/GophKeeper/internal/proto/auth"
	kpb "github.com/rombintu/GophKeeper/internal/proto/keeper"
)

type MemoryDriver struct {
	Users   []*apb.User
	Secrets []*kpb.Secret
}

func (md *MemoryDriver) UserGet(ctx context.Context, user *apb.User) (*apb.User, error) {
	for _, u := range md.Users {
		if user.GetEmail() == u.GetEmail() {
			slog.Debug("founded user",
				slog.String("email", u.GetEmail()),
				slog.String("fingerprint", string(u.GetHexKeys())))
			return u, nil
		}
	}
	return nil, errors.New("user not found")

}

func (md *MemoryDriver) UserCreate(ctx context.Context, user *apb.User) error {
	for _, u := range md.Users {
		if user.GetEmail() == u.GetEmail() {
			return errors.New("user already exists")
		}
	}
	md.Users = append(md.Users, user)
	slog.Debug("user created",
		slog.String("email", user.GetEmail()),
		slog.String("fingerprint", string(user.GetHexKeys())))
	return nil
}

func (md *MemoryDriver) SecretCreate(ctx context.Context, secret *kpb.Secret) error {
	for _, s := range md.Secrets {
		if s.GetTitle() == secret.GetTitle() && s.GetUserEmail() == secret.GetUserEmail() {
			secret.Version = s.GetVersion() + 1
			break
		}
	}
	md.Secrets = append(md.Secrets, secret)
	return nil
}

func (md *MemoryDriver) SecretCreateBatch(ctx context.Context, secrets []*kpb.Secret) error {
	for _, s1 := range secrets {
		for _, s2 := range md.Secrets {
			if s1.GetTitle() == s2.GetTitle() && s1.GetUserEmail() == s2.GetUserEmail() {
				s1.Version = s2.GetVersion() + 1
				break
			}
		}
		md.Secrets = append(md.Secrets, s1)
	}
	return nil
}

func (md *MemoryDriver) SecretList(ctx context.Context, userEmail string) ([]*kpb.Secret, error) {
	var founded []*kpb.Secret
	for _, s := range md.Secrets {
		if s.GetUserEmail() == userEmail {
			founded = append(founded, s)
		}
	}
	return founded, nil
}

func (md *MemoryDriver) SecretPurge(ctx context.Context, secret *kpb.Secret) error {
	// А смысл заморачиваться?
	return nil
}

func (md *MemoryDriver) Ping(ctx context.Context, monitoring bool) error {
	return nil
}

func (md *MemoryDriver) Open(ctx context.Context) error {
	md.Users = []*apb.User{}
	md.Secrets = []*kpb.Secret{}
	return nil
}

func (md *MemoryDriver) Close(ctx context.Context) error {
	md.Users = []*apb.User{}
	md.Secrets = []*kpb.Secret{}
	return nil
}

func (md *MemoryDriver) Configure(ctx context.Context) error {
	return nil
}
