package sync_test

import (
	"context"
	"net"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
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

	// 1. Создаем мокированный keeper client
	mockKeeper := mock_keeper.NewMockKeeperClient(ctrl)

	// 2. Создаем тестовый сервер с bufconn
	lis := bufconn.Listen(1024 * 1024)
	bufDialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}

	// 3. Настраиваем сервис
	service := sync.NewSyncService(config.SyncConfig{}).(*sync.SyncService)
	service.TestClientConn = mockKeeper // Инжектим мок напрямую

	// 4. Настраиваем ожидания
	mockKeeper.EXPECT().Fetch(
		gomock.Any(),
		&kpb.FetchRequest{UserEmail: "test@example.com"},
	).Return(&kpb.FetchResponse{
		Secrets: []*kpb.Secret{
			{Title: "secret2", HashPayload: "hash2"},
		},
	}, nil)

	mockKeeper.EXPECT().CreateMany(
		gomock.Any(),
		gomock.AssignableToTypeOf(&kpb.CreateBatchRequest{}),
	).Return(&emptypb.Empty{}, nil)

	// 5. Запускаем gRPC сервер
	srv := grpc.NewServer()
	spb.RegisterSyncServer(srv, service)
	go func() {
		if err := srv.Serve(lis); err != nil {
			t.Logf("Server exited with error: %v", err)
		}
	}()
	defer srv.Stop()

	// 6. Создаем клиент для тестирования
	ctx := context.Background()
	conn, _ := grpc.DialContext(ctx, "bufnet",
		grpc.WithContextDialer(bufDialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer conn.Close()

	client := spb.NewSyncClient(conn)

	// 7. Выполняем тестовый запрос
	resp, err := client.Process(ctx, &spb.SyncRequest{
		Email: "test@example.com",
		Secrets: []*kpb.Secret{
			{Title: "secret1", HashPayload: "hash1"},
		},
	})

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
