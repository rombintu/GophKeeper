package auth

import (
	"context"
	"log/slog"
	"reflect"

	apb "github.com/rombintu/GophKeeper/internal/proto/auth"
	"github.com/rombintu/GophKeeper/lib/common"
	"github.com/rombintu/GophKeeper/lib/jwt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *AuthService) Register(ctx context.Context, in *apb.RegisterRequest) (*apb.RegisterResponse, error) {
	if err := common.UserValidate(in.GetUser()); err != nil {
		slog.Error("message", slog.String("func",
			common.DotJoin(ServiceName, "Register", "UserValidate")), slog.String("error", err.Error()))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := s.store.UserCreate(context.Background(), in.GetUser()); err != nil {
		slog.Error("message", slog.String("func",
			common.DotJoin(ServiceName, "Register", "UserCreate")), slog.String("error", err.Error()))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	token, err := jwt.NewToken(in.GetUser(), s.config.Secret, s.config.TokenExp)
	if err != nil {
		slog.Error("message", slog.String("func",
			common.DotJoin(ServiceName, "Register", "NewToken")), slog.String("error", err.Error()))
		return nil, status.Error(codes.Internal, "failed create token")
	}
	r := apb.RegisterResponse{
		Token: token,
	}

	return &r, nil
}

func (s *AuthService) Login(ctx context.Context, in *apb.LoginRequest) (*apb.LoginResponse, error) {
	if err := common.UserValidate(in.GetUser()); err != nil {
		slog.Error("message", slog.String("func",
			common.DotJoin(ServiceName, "Login", "UserValidate")), slog.String("error", err.Error()))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	userFounded, err := s.store.UserGet(context.Background(), in.GetUser())
	if err != nil {
		slog.Error("message", slog.String("func",
			common.DotJoin(ServiceName, "Login", "UserGet")), slog.String("error", err.Error()))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	r := apb.LoginResponse{}
	if reflect.DeepEqual(in.User.GetHexKeys(), userFounded.GetHexKeys()) {
		token, err := jwt.NewToken(in.GetUser(), s.config.Secret, s.config.TokenExp)
		if err != nil {
			slog.Error("message", slog.String("func",
				common.DotJoin(ServiceName, "Login", "NewToken")), slog.String("error", err.Error()))
			return nil, status.Error(codes.Internal, "failed create token")
		}
		r.Token = token
	} else {
		return nil, status.Error(codes.InvalidArgument, "user keys are not equal")
	}
	return &r, err
}
