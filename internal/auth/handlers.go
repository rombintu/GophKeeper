package auth

import (
	"context"
	"log/slog"
	"reflect"

	pb "github.com/rombintu/GophKeeper/internal/proto"
	"github.com/rombintu/GophKeeper/internal/storage"
	"github.com/rombintu/GophKeeper/lib/common"
	"github.com/rombintu/GophKeeper/lib/jwt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

func (s *AuthService) Register(ctx context.Context, in *pb.UserRequest) (*pb.AuthResponse, error) {
	if err := common.UserValidate(in.User); err != nil {
		slog.Debug("message", slog.String("func", "AuthService.Register"))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// TODO: сделать отдельную функцию для унификации и уменьшения кода
	conn, err := grpc.NewClient(s.config.StorageServiceAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		slog.Error("error dial to client",
			slog.String("from", ServiceName),
			slog.String("to", storage.ServiceName),
			slog.String("error", err.Error()),
		)
	}
	defer conn.Close()
	storageClient := pb.NewStorageClient(conn)
	if _, err := storageClient.UserCreate(ctx, in); err != nil {
		slog.Error("message", slog.String("func", "AuthService.Register"), slog.String("error", err.Error()))
		return nil, err
	}

	token, err := jwt.NewToken(in.User, s.config.Secret, s.config.TokenExp)
	if err != nil {
		return nil, err
	}
	r := pb.AuthResponse{
		Token: token,
	}

	return &r, err
}

func (s *AuthService) Login(ctx context.Context, in *pb.UserRequest) (*pb.AuthResponse, error) {
	if err := common.UserValidate(in.User); err != nil {
		slog.Debug("message", slog.String("func", "AuthService.Login"))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// TODO: сделать отдельную функцию для унификации и уменьшения кода
	conn, err := grpc.NewClient(s.config.StorageServiceAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		slog.Error("error dial to client",
			slog.String("from", ServiceName),
			slog.String("to", storage.ServiceName),
			slog.String("error", err.Error()),
		)
	}
	defer conn.Close()
	storageClient := pb.NewStorageClient(conn)

	newUserRequest := pb.UserRequest{User: &pb.User{Email: in.User.GetEmail()}}
	if _, err := storageClient.UserGet(ctx, &newUserRequest); err != nil {
		slog.Error("message", slog.String("func", "AuthService.Login"), slog.String("error", err.Error()))
		return nil, err
	}

	r := pb.AuthResponse{}
	if reflect.DeepEqual(newUserRequest.User.GetHexKeys(), in.User.GetHexKeys()) {
		token, err := jwt.NewToken(in.User, s.config.Secret, s.config.TokenExp)
		if err != nil {
			slog.Error("message", slog.String("func", "AuthService.Login"), slog.String("error", err.Error()))
			return nil, err
		}
		r.Token = token
	} else {
		return nil, status.Error(codes.InvalidArgument, "user keys are not equal")
	}
	return &r, err
}
