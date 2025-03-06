package keeper

import (
	"context"
	"log/slog"

	kpb "github.com/rombintu/GophKeeper/internal/proto/keeper"
	"github.com/rombintu/GophKeeper/lib/common"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *KeeperService) Fetch(ctx context.Context, in *kpb.FetchRequest) (*kpb.FetchResponse, error) {
	secrets, err := s.store.SecretList(context.Background(), in.GetUserEmail())
	if err != nil {
		slog.Error("message", slog.String("func",
			common.DotJoin(ServiceName, "Fetch", "SecretList")), slog.String("error", err.Error()))
		return nil, err
	}
	r := kpb.FetchResponse{}
	r.Secrets = secrets
	return &r, nil
}

func (s *KeeperService) Create(ctx context.Context, in *kpb.CreateRequest) (*emptypb.Empty, error) {
	if err := s.store.SecretCreate(context.Background(), in.GetSecret()); err != nil {
		slog.Error("message", slog.String("func",
			common.DotJoin(ServiceName, "Create", "SecretCreate")), slog.String("error", err.Error()))
		return &emptypb.Empty{}, err
	}

	return &emptypb.Empty{}, nil
}

func (s *KeeperService) CreateBatch(ctx context.Context, in *kpb.CreateBatchRequest) (*emptypb.Empty, error) {
	if err := s.store.SecretCreateBatch(context.Background(), in.GetSecrets()); err != nil {
		slog.Error("message", slog.String("func",
			common.DotJoin(ServiceName, "CreateBatch", "SecretCreateBatch")), slog.String("error", err.Error()))
		return &emptypb.Empty{}, err
	}

	return &emptypb.Empty{}, nil
}

func (s *KeeperService) Delete(ctx context.Context, in *kpb.DeleteRequest) (*emptypb.Empty, error) {
	if err := s.store.SecretPurge(context.Background(), in.GetSecret()); err != nil {
		slog.Error("message", slog.String("func",
			common.DotJoin(ServiceName, "Delete", "SecretPurge")), slog.String("error", err.Error()))
		return &emptypb.Empty{}, err
	}

	return &emptypb.Empty{}, nil
}
