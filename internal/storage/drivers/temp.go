package drivers

import (
	"errors"
	"log/slog"

	pb "github.com/rombintu/GophKeeper/internal/proto"
)

type MemoryDriver struct {
	Users   []*pb.User
	Secrets []*pb.Secret
}

func (md *MemoryDriver) UserGet(user *pb.User) (*pb.User, error) {
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

func (md *MemoryDriver) UserCreate(user *pb.User) error {
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

func (md *MemoryDriver) SecretGet(userID int64) (*pb.Secret, error) {
	return &pb.Secret{}, nil
}

func (md *MemoryDriver) SecretCreate(secret *pb.Secret) error {
	return nil
}

func (md *MemoryDriver) SecretList(userID int64) ([]*pb.Secret, error) {
	return []*pb.Secret{}, nil
}

func (md *MemoryDriver) Ping() error {
	return nil
}

func (md *MemoryDriver) Open() error {
	md.Users = []*pb.User{}
	md.Secrets = []*pb.Secret{}
	return nil
}

func (md *MemoryDriver) Close() error {
	md.Users = []*pb.User{}
	md.Secrets = []*pb.Secret{}
	return nil
}
