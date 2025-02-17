package drivers

import (
	"errors"

	pb "github.com/rombintu/GophKeeper/internal/proto"
)

type MemoryDriver struct {
	Users   []*pb.User
	Secrets []*pb.Secret
}

func (md *MemoryDriver) UserGet(user *pb.User) error {
	found := false
	for _, u := range md.Users {
		if user.GetEmail() == u.GetEmail() {
			user.HexKeys = u.GetHexKeys()
			found = true
		}
	}
	if !found {
		return errors.New("user not found")
	}
	return nil
}

func (md *MemoryDriver) UserCreate(user *pb.User) error {
	for _, u := range md.Users {
		if user.GetEmail() == u.GetEmail() {
			return errors.New("user already exists")
		}
	}
	md.Users = append(md.Users, user)
	return nil
}

func (md *MemoryDriver) SecretGet(userID int64) (*pb.Secret, error) {
	return &pb.Secret{}, nil
}

func (md *MemoryDriver) SecretCreate(secret *pb.Secret) error {
	return nil
}

func (md *MemoryDriver) SecretsGet(userID int64) ([]*pb.Secret, error) {
	return []*pb.Secret{}, nil
}

func (md *MemoryDriver) Ping() error {
	return nil
}
