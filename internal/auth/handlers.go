package auth

import (
	"context"
	"log/slog"
	"reflect"

	"github.com/rombintu/GophKeeper/internal/proto"
	"github.com/rombintu/GophKeeper/lib/common"
	"github.com/rombintu/GophKeeper/lib/jwt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *AuthService) Register(ctx context.Context, in *proto.RegisterRequest) (*proto.RegisterResponse, error) {
	if err := common.UserValidate(in.User); err != nil {
		slog.Debug("message", slog.String("func", "AuthService.Register"))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := s.store.UserCreate(in.User); err != nil {
		slog.Debug("message", slog.String("func", "AuthService.Register"))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	token, err := jwt.NewToken(in.User, s.config.Secret, s.config.TokenExp)
	if err != nil {
		return nil, err
	}
	r := proto.RegisterResponse{
		Token: token,
	}

	return &r, err
}

func (s *AuthService) Login(ctx context.Context, in *proto.LoginRequest) (*proto.LoginResponse, error) {
	if err := common.UserValidate(in.User); err != nil {
		slog.Debug("message", slog.String("func", "AuthService.Login"))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	userFounded, err := s.store.UserGet(in.User)
	if err != nil {
		slog.Debug("message", slog.String("func", "AuthService.Register"))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	r := proto.LoginResponse{}
	if reflect.DeepEqual(in.User.GetHexKeys(), userFounded.GetHexKeys()) {
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
