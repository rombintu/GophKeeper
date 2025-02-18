package storage

import (
	"context"

	pb "github.com/rombintu/GophKeeper/internal/proto"
)

func (s *StorageService) UserGet(ctx context.Context, in *pb.UserRequest) (*pb.UserResponse, error) {
	user := pb.User{
		Email: in.User.Email,
	}
	r := pb.UserResponse{}
	err := s.driver.UserGet(&user)
	if err != nil {
		return nil, err
	}
	r.User = &user
	return &r, err
}

func (s *StorageService) UserCreate(ctx context.Context, in *pb.UserRequest) (*pb.UserResponse, error) {
	user := pb.User{
		Email: in.User.Email,
	}
	r := pb.UserResponse{}
	err := s.driver.UserCreate(&user)
	if err != nil {
		return nil, err
	}
	r.User = &user
	return &r, err
}
