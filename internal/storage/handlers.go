package storage

import (
	"context"
	"log/slog"

	pb "github.com/rombintu/GophKeeper/internal/proto"
	"github.com/rombintu/GophKeeper/lib/common"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *StorageService) UserGet(ctx context.Context, in *pb.UserRequest) (*pb.UserResponse, error) {
	if err := common.UserValidate(in.User); err != nil {
		slog.Debug("message", slog.String("func", "AuthService.Register"))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
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
	if err := common.UserValidate(in.User); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
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
