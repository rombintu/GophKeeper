package keeper

import (
	"context"
	"log/slog"

	kpb "github.com/rombintu/GophKeeper/internal/proto/keeper"
	"github.com/rombintu/GophKeeper/lib/common"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *KeeperService) Fetch(ctx context.Context, in *kpb.FetchRequest) (*kpb.FetchResponse, error) {
	secrets, err := s.store.SecretList(in.GetUserId(), in.GetPattern())
	if err != nil {
		slog.Debug("message", slog.String("func",
			common.DotJoin(ServiceName, "Fetch", "SecretList")), slog.String("error", err.Error()))
		return nil, err
	}
	r := kpb.FetchResponse{}
	r.Secrets = secrets
	return &r, nil
}

func (s *KeeperService) Create(ctx context.Context, in *kpb.CreateRequest) (*emptypb.Empty, error) {
	if err := s.store.SecretCreate(in.GetSecret()); err != nil {
		slog.Debug("message", slog.String("func",
			common.DotJoin(ServiceName, "Create", "SecretCreate")), slog.String("error", err.Error()))
		return &emptypb.Empty{}, err
	}

	return &emptypb.Empty{}, nil
}
