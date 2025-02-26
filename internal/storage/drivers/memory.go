package drivers

import (
	"errors"
	"log/slog"
	"strings"

	apb "github.com/rombintu/GophKeeper/internal/proto/auth"
	kpb "github.com/rombintu/GophKeeper/internal/proto/keeper"
)

type MemoryDriver struct {
	Users   []*apb.User
	Secrets []*kpb.Secret
}

func (md *MemoryDriver) UserGet(user *apb.User) (*apb.User, error) {
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

func (md *MemoryDriver) UserCreate(user *apb.User) error {
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

func (md *MemoryDriver) SecretCreate(secret *kpb.Secret) error {
	for _, s := range md.Secrets {
		if s.GetTitle() == secret.GetTitle() && s.GetUserId() == secret.GetUserId() {
			secret.Version = s.GetVersion() + 1
			break
		}
	}
	md.Secrets = append(md.Secrets, secret)
	return nil
}

func (md *MemoryDriver) SecretList(userID int64, pattern string) ([]*kpb.Secret, error) {
	var founded []*kpb.Secret
	all := false
	if pattern == "" || pattern == "*" {
		all = true
	}
	for _, s := range md.Secrets {
		if s.GetUserId() == userID {
			if !all && !strings.Contains(s.Title, pattern) {
				continue
			}
			founded = append(founded, s)
		}
	}
	return founded, nil
}

func (md *MemoryDriver) Ping() error {
	return nil
}

func (md *MemoryDriver) Open() error {
	md.Users = []*apb.User{}
	md.Secrets = []*kpb.Secret{}
	return nil
}

func (md *MemoryDriver) Close() error {
	md.Users = []*apb.User{}
	md.Secrets = []*kpb.Secret{}
	return nil
}
