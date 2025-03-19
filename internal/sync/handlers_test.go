package sync_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/rombintu/GophKeeper/internal/config"
	kpb "github.com/rombintu/GophKeeper/internal/proto/keeper"
	spb "github.com/rombintu/GophKeeper/internal/proto/sync"

	sync "github.com/rombintu/GophKeeper/internal/sync"
	mock_keeper "github.com/rombintu/GophKeeper/internal/sync/mocks"
)

func TestProcess_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockKeeper := mock_keeper.NewMockKeeperClient(ctrl)
	service := sync.NewSyncService(config.SyncConfig{})

	require.NoError(t, service.Configure())

	// Setup mocks
	req := &spb.SyncRequest{
		Email: "test@example.com",
		Secrets: []*kpb.Secret{
			{Title: "secret1", HashPayload: "hash1"},
		},
	}

	mockKeeper.EXPECT().Fetch(gomock.Any(), &kpb.FetchRequest{
		UserEmail: "test@example.com",
	}).Return(&kpb.FetchResponse{
		Secrets: []*kpb.Secret{
			{Title: "secret2", HashPayload: "hash2"},
		},
	}, nil)

	mockKeeper.EXPECT().CreateMany(
		gomock.Any(),
		gomock.AssignableToTypeOf(&kpb.CreateBatchRequest{}),
	).Return(&emptypb.Empty{}, nil)

	resp, err := service.(*sync.SyncService).Process(context.Background(), req)
	require.NoError(t, err)
	assert.True(t, resp.Success)
	assert.Len(t, resp.Secrets, 1)
}

func TestProcess_VersionConflict(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockKeeper := mock_keeper.NewMockKeeperClient(ctrl)
	service := sync.NewSyncService(config.SyncConfig{})
	require.NoError(t, service.Configure())

	req := &spb.SyncRequest{
		Email: "test@example.com",
		Secrets: []*kpb.Secret{
			{Title: "secret1", HashPayload: "hash1", Version: 1},
		},
	}

	mockKeeper.EXPECT().Fetch(gomock.Any(), gomock.Any()).Return(&kpb.FetchResponse{
		Secrets: []*kpb.Secret{
			{Title: "secret1", HashPayload: "old_hash", Version: 2},
		},
	}, nil)

	mockKeeper.EXPECT().CreateMany(gomock.Any(), gomock.Any()).DoAndReturn(
		func(_ context.Context, req *kpb.CreateBatchRequest) (*emptypb.Empty, error) {
			assert.Equal(t, int64(3), req.Secrets[0].Version)
			return nil, nil
		},
	)

	resp, err := service.(*sync.SyncService).Process(context.Background(), req)
	require.NoError(t, err)
	assert.True(t, resp.Success)
}

func TestProcess_ErrorCases(t *testing.T) {
	t.Run("Fetch error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockKeeper := mock_keeper.NewMockKeeperClient(ctrl)
		service := sync.NewSyncService(config.SyncConfig{})
		require.NoError(t, service.Configure())

		mockKeeper.EXPECT().Fetch(gomock.Any(), gomock.Any()).Return(
			nil, status.Error(codes.Internal, "server error"),
		)

		_, err := service.(*sync.SyncService).Process(context.Background(), &spb.SyncRequest{})
		require.Error(t, err)
		assert.Equal(t, codes.Internal, status.Code(err))
	})

	t.Run("CreateMany error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockKeeper := mock_keeper.NewMockKeeperClient(ctrl)
		service := sync.NewSyncService(config.SyncConfig{})
		require.NoError(t, service.Configure())

		mockKeeper.EXPECT().Fetch(gomock.Any(), gomock.Any()).Return(
			&kpb.FetchResponse{}, nil,
		)

		mockKeeper.EXPECT().CreateMany(gomock.Any(), gomock.Any()).Return(
			nil, status.Error(codes.Internal, "create error"),
		)

		_, err := service.(*sync.SyncService).Process(context.Background(), &spb.SyncRequest{
			Secrets: []*kpb.Secret{{}},
		})
		require.Error(t, err)
		assert.Equal(t, codes.Internal, status.Code(err))
	})
}
