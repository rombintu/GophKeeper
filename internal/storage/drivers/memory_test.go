package drivers

import (
	"context"
	"testing"

	kpb "github.com/rombintu/GophKeeper/internal/proto/keeper"
)

func TestMemoryDriver_SecretCreate(t *testing.T) {
	md := &MemoryDriver{}
	ctx := context.Background()
	if err := md.Open(ctx); err != nil {
		t.Error("failed open memdb", err)
	}
	tests := []struct {
		name     string
		md       *MemoryDriver
		secret   *kpb.Secret
		wantErr  bool
		wantSize int
	}{
		{
			name:     "create_test_secret",
			md:       md,
			secret:   &kpb.Secret{Title: "test", UserEmail: "1"},
			wantErr:  false,
			wantSize: 1,
		},
		{
			name:     "create_test_secret_same",
			md:       md,
			secret:   &kpb.Secret{Title: "test", UserEmail: "1"},
			wantErr:  false,
			wantSize: 2,
		},
		{
			name:     "create_test_secret_more",
			md:       md,
			secret:   &kpb.Secret{Title: "test", UserEmail: "2"},
			wantErr:  false,
			wantSize: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := md.SecretCreate(ctx, tt.secret); (err != nil) != tt.wantErr {
				t.Errorf("MemoryDriver.SecretCreate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if len(md.Secrets) != tt.wantSize {
				t.Errorf("%+v : want size %d, got size %d", md.Secrets, tt.wantSize, len(md.Secrets))
			}
		})
	}
}

func TestMemoryDriver_SecretList(t *testing.T) {
	md := &MemoryDriver{}
	ctx := context.Background()
	if err := md.Open(ctx); err != nil {
		t.Error("failed open memdb", err)
	}

	s1 := &kpb.Secret{Title: "test", UserEmail: "1"}
	s2 := &kpb.Secret{Title: "test", UserEmail: "1"}
	s3 := &kpb.Secret{Title: "test", UserEmail: "2"}
	md.Secrets = append(md.Secrets,
		s1, s2, s3,
	)
	type args struct {
		UserEmail string
	}
	tests := []struct {
		name     string
		args     args
		wantSize int
		wantErr  bool
	}{
		{
			name:     "unknown_user_id",
			args:     args{UserEmail: "0"},
			wantSize: 0,
			wantErr:  false,
		},
		{
			name:     "1_user_id_all",
			args:     args{UserEmail: "1"},
			wantSize: 2,
			wantErr:  false,
		},
		{
			name:     "1_user_id_notfounded",
			args:     args{UserEmail: "3"},
			wantSize: 0,
			wantErr:  false,
		},
		{
			name:     "2_user_id_all",
			args:     args{UserEmail: "2"},
			wantSize: 1,
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := md.SecretList(ctx, tt.args.UserEmail)
			if (err != nil) != tt.wantErr {
				t.Errorf("MemoryDriver.SecretList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != tt.wantSize {
				t.Errorf("%+v : want size %d, got size %d", got, tt.wantSize, len(got))
			}
		})
	}
}
