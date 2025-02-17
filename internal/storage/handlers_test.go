package storage

import (
	"context"
	"testing"

	"github.com/rombintu/GophKeeper/internal/config"
	"github.com/rombintu/GophKeeper/internal/proto"
)

func TestStorageService_UserCreat(t *testing.T) {

	service := NewStorageService(config.StorageConfig{DriverPath: memDriver})

	tests := []struct {
		name    string
		newUser *proto.User
		wantErr bool
	}{
		{
			name:    "user_create_ok",
			newUser: &proto.User{Email: "email1"},
			wantErr: false,
		},
		{
			name:    "user_create_empty_email",
			newUser: &proto.User{Email: ""},
			wantErr: true,
		},
		{
			name:    "user_create_with_keys",
			newUser: &proto.User{Email: "email2", HexKeys: []byte("123")},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctx := context.TODO()
			_, err := service.UserCreate(ctx, &proto.UserRequest{User: tt.newUser})
			if err != nil && tt.wantErr == false {
				t.Fatalf("Ошибка при создании пользователя. ERROR: %+v", err)
			}

		})
	}

}
