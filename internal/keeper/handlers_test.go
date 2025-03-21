package keeper_test

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/rombintu/GophKeeper/internal/config"
	"github.com/rombintu/GophKeeper/internal/keeper"
	kpb "github.com/rombintu/GophKeeper/internal/proto/keeper"
	mock_storage "github.com/rombintu/GophKeeper/internal/storage/mocks"
)

func TestKeeperService_Fetch(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mock_storage.NewMockSecretManager(ctrl)
	cfg := config.KeeperConfig{
		Config: config.Config{
			Secret:               "test-secret",
			Env:                  "test",
			HealthCheckDuration:  10 * time.Second,
			KeeperServiceAddress: ":3202",
		},
	}

	service := keeper.NewKeeperService(mockStore, cfg).(*keeper.KeeperService)

	t.Run("Success fetch", func(t *testing.T) {
		expectedSecrets := []*kpb.Secret{
			{Title: "test-secret", UserEmail: "test@example.com"},
		}

		mockStore.EXPECT().SecretList(gomock.Any(), "test@example.com").Return(expectedSecrets, nil)

		resp, err := service.Fetch(context.Background(), &kpb.FetchRequest{
			UserEmail: "test@example.com",
		})

		require.NoError(t, err)
		assert.Len(t, resp.Secrets, 1)
		assert.Equal(t, "test-secret", resp.Secrets[0].Title)
	})

	t.Run("Storage error", func(t *testing.T) {
		mockStore.EXPECT().SecretList(gomock.Any(), "error@example.com").Return(nil, assert.AnError)

		_, err := service.Fetch(context.Background(), &kpb.FetchRequest{
			UserEmail: "error@example.com",
		})

		assert.Error(t, err)
	})
}

func TestKeeperService_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mock_storage.NewMockSecretManager(ctrl)
	cfg := config.KeeperConfig{
		Config: config.Config{
			Secret:               "test-secret",
			KeeperServiceAddress: ":3202",
		},
	}

	service := keeper.NewKeeperService(mockStore, cfg).(*keeper.KeeperService)

	t.Run("Success create", func(t *testing.T) {
		secret := &kpb.Secret{
			Title:     "new-secret",
			UserEmail: "test@example.com",
		}

		mockStore.EXPECT().SecretCreate(gomock.Any(), secret).Return(nil)

		_, err := service.Create(context.Background(), &kpb.CreateRequest{
			UserEmail: "test@example.com",
			Secret:    secret,
		})

		assert.NoError(t, err)
	})

	t.Run("Create error", func(t *testing.T) {
		secret := &kpb.Secret{Title: "invalid-secret"}

		mockStore.EXPECT().SecretCreate(gomock.Any(), secret).Return(assert.AnError)

		_, err := service.Create(context.Background(), &kpb.CreateRequest{
			Secret: secret,
		})

		assert.Error(t, err)
	})
}

func TestKeeperService_CreateBatch(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mock_storage.NewMockSecretManager(ctrl)
	cfg := config.KeeperConfig{Config: config.Config{Secret: "test-secret"}}
	service := keeper.NewKeeperService(mockStore, cfg).(*keeper.KeeperService)

	t.Run("Success batch create", func(t *testing.T) {
		secrets := []*kpb.Secret{
			{Title: "secret1"},
			{Title: "secret2"},
		}

		mockStore.EXPECT().SecretCreateBatch(gomock.Any(), secrets).Return(nil)

		_, err := service.CreateBatch(context.Background(), &kpb.CreateBatchRequest{
			Secrets: secrets,
		})

		assert.NoError(t, err)
	})

	t.Run("Batch create error", func(t *testing.T) {
		mockStore.EXPECT().SecretCreateBatch(gomock.Any(), gomock.Any()).Return(assert.AnError)

		_, err := service.CreateBatch(context.Background(), &kpb.CreateBatchRequest{})
		assert.Error(t, err)
	})
}

func TestKeeperService_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mock_storage.NewMockSecretManager(ctrl)
	cfg := config.KeeperConfig{Config: config.Config{Secret: "test-secret"}}
	service := keeper.NewKeeperService(mockStore, cfg).(*keeper.KeeperService)

	t.Run("Success delete", func(t *testing.T) {
		secret := &kpb.Secret{Title: "to-delete"}

		mockStore.EXPECT().SecretPurge(gomock.Any(), secret).Return(nil)

		_, err := service.Delete(context.Background(), &kpb.DeleteRequest{
			Secret: secret,
		})

		assert.NoError(t, err)
	})

	t.Run("Delete error", func(t *testing.T) {
		secret := &kpb.Secret{Title: "invalid"}

		mockStore.EXPECT().SecretPurge(gomock.Any(), secret).Return(assert.AnError)

		_, err := service.Delete(context.Background(), &kpb.DeleteRequest{
			Secret: secret,
		})

		assert.Error(t, err)
	})
}
