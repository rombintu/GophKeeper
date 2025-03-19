package drivers

import (
	"context"
	"os"
	"testing"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	apb "github.com/rombintu/GophKeeper/internal/proto/auth"
	kpb "github.com/rombintu/GophKeeper/internal/proto/keeper"
)

func TestPgxDriver_SecretCreateBatch(t *testing.T) {
	dbPath := os.Getenv("PGX_DB_PATH")
	type args struct {
		secrets []*kpb.Secret
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "add_secrets",
			args: args{
				secrets: []*kpb.Secret{
					{
						Title:       "test1",
						SecretType:  kpb.Secret_TEXT,
						UserEmail:   "email1",
						Version:     0,
						HashPayload: "123",
						Payload:     []byte("hello"),
					},
					{
						Title:       "test1",
						SecretType:  kpb.Secret_TEXT,
						UserEmail:   "email1",
						Version:     1,
						HashPayload: "123",
						Payload:     []byte("hello"),
					},
					{
						Title:       "test1",
						SecretType:  kpb.Secret_TEXT,
						UserEmail:   "email1",
						Version:     1,
						HashPayload: "1234",
						Payload:     []byte("hello1"),
					},
				},
			},
		},
	}
	for _, tt := range tests {
		ctx := context.Background()
		t.Run(tt.name, func(t *testing.T) {
			d := &PgxDriver{
				dbURL:       dbPath,
				serviceName: "text",
			}
			if err := d.Open(ctx); err != nil {
				t.Error(err)
			}
			if err := d.SecretCreateBatch(ctx, tt.args.secrets); (err != nil) != tt.wantErr {
				t.Errorf("PgxDriver.SecretCreateBatch() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPgxDriver_UserCreate(t *testing.T) {
	dbPath := os.Getenv("PGX_DB_PATH")
	type args struct {
		user *apb.User
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "user_create_one",
			args:    args{user: &apb.User{Email: "test.com"}},
			wantErr: false,
		},
		{
			name:    "user_create_two",
			args:    args{user: &apb.User{Email: "test2.com"}},
			wantErr: false,
		},
		{
			name:    "user_create_same",
			args:    args{user: &apb.User{Email: "test.com"}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		ctx := context.Background()
		t.Run(tt.name, func(t *testing.T) {
			d := &PgxDriver{
				dbURL:       dbPath,
				serviceName: "test",
			}
			if err := d.UserCreate(ctx, tt.args.user); (err != nil) != tt.wantErr {
				t.Errorf("PgxDriver.UserCreate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
