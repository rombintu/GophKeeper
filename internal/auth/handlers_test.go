package auth_test

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/rombintu/GophKeeper/internal/auth"
	"github.com/rombintu/GophKeeper/internal/config"
	apb "github.com/rombintu/GophKeeper/internal/proto/auth"
	mock_storage "github.com/rombintu/GophKeeper/internal/storage/mocks"
)

func TestAuthService_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mock_storage.NewMockUserManager(ctrl)
	cfg := config.AuthConfig{
		Config: config.Config{
			Secret: "test-secret",
			// Добавляем обязательные поля
			Env:                 "test",
			HealthCheckDuration: time.Second * 10,
			AuthServiceAddress:  ":3201",
		},
		TokenExp: time.Hour,
	}

	service := auth.NewAuthService(mockStore, cfg).(*auth.AuthService)

	t.Run("Success registration", func(t *testing.T) {
		req := &apb.RegisterRequest{
			User: &apb.User{
				Email:       "test@example.com",
				KeyChecksum: []byte("valid-checksum"),
			},
		}

		mockStore.EXPECT().UserCreate(gomock.Any(), req.User).Return(nil)

		resp, err := service.Register(context.Background(), req)
		require.NoError(t, err)
		assert.NotEmpty(t, resp.Token)
	})

	t.Run("Invalid user data", func(t *testing.T) {
		req := &apb.RegisterRequest{
			User: &apb.User{Email: "invalid-email"},
		}

		_, err := service.Register(context.Background(), req)
		assert.Error(t, err)
		assert.Equal(t, codes.InvalidArgument, status.Code(err))
	})

	t.Run("User creation error", func(t *testing.T) {
		req := &apb.RegisterRequest{
			User: &apb.User{
				Email:       "test@example.com",
				KeyChecksum: []byte("valid-checksum"),
			},
		}

		mockStore.EXPECT().UserCreate(gomock.Any(), req.User).Return(assert.AnError)

		_, err := service.Register(context.Background(), req)
		assert.Error(t, err)
		assert.Equal(t, codes.InvalidArgument, status.Code(err))
	})
}

func TestAuthService_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mock_storage.NewMockUserManager(ctrl)
	cfg := config.AuthConfig{
		Config: config.Config{
			Secret: "test-secret",
			// Добавляем обязательные поля
			Env:                 "test",
			HealthCheckDuration: time.Second * 10,
			AuthServiceAddress:  ":3201",
		},
		TokenExp: time.Hour,
	}

	service := auth.NewAuthService(mockStore, cfg).(*auth.AuthService)

	t.Run("Success login", func(t *testing.T) {
		req := &apb.LoginRequest{
			User: &apb.User{
				Email:       "test@example.com",
				KeyChecksum: []byte("valid-checksum"),
			},
		}

		storedUser := &apb.User{
			Email:       "test@example.com",
			KeyChecksum: []byte("valid-checksum"),
		}

		mockStore.EXPECT().UserGet(gomock.Any(), req.User).Return(storedUser, nil)

		resp, err := service.Login(context.Background(), req)
		require.NoError(t, err)
		assert.NotEmpty(t, resp.Token)
	})

	t.Run("Invalid credentials", func(t *testing.T) {
		req := &apb.LoginRequest{
			User: &apb.User{
				Email:       "test@example.com",
				KeyChecksum: []byte("invalid-checksum"),
			},
		}

		storedUser := &apb.User{
			Email:       "test@example.com",
			KeyChecksum: []byte("valid-checksum"),
		}

		mockStore.EXPECT().UserGet(gomock.Any(), req.User).Return(storedUser, nil)

		_, err := service.Login(context.Background(), req)
		assert.Error(t, err)
		assert.Equal(t, codes.InvalidArgument, status.Code(err))
	})

	t.Run("User not found", func(t *testing.T) {
		req := &apb.LoginRequest{
			User: &apb.User{
				Email:       "notfound@example.com",
				KeyChecksum: []byte("checksum"),
			},
		}

		mockStore.EXPECT().UserGet(gomock.Any(), req.User).Return(nil, assert.AnError)

		_, err := service.Login(context.Background(), req)
		assert.Error(t, err)
		assert.Equal(t, codes.InvalidArgument, status.Code(err))
	})
}
