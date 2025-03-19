package sync

import (
	"context"
	"fmt"
	"log/slog"

	kpb "github.com/rombintu/GophKeeper/internal/proto/keeper"
	spb "github.com/rombintu/GophKeeper/internal/proto/sync"
	"github.com/rombintu/GophKeeper/lib/common"
)

func (s *SyncService) Process(ctx context.Context, in *spb.SyncRequest) (*spb.SyncResponse, error) {
	serverSecrets, err := s.getServerSecrets(ctx, in.GetEmail())
	if err != nil {
		slog.Error("message", slog.String("func",
			common.DotJoin(ServiceName, "Process", "getServerSecrets")), slog.String("error", err.Error()))
		return nil, err
	}

	serverSecretsMap := make(map[string]*kpb.Secret)
	clientSecretsMap := make(map[string]*kpb.Secret)

	for _, secret := range serverSecrets {
		key := fmt.Sprintf("%s:%s", secret.GetTitle(), secret.GetHashPayload())
		serverSecretsMap[key] = secret
	}

	var secretsToCreate []*kpb.Secret
	for _, clientSecret := range in.Secrets {
		key := fmt.Sprintf("%s:%s", clientSecret.Title, clientSecret.HashPayload)
		clientSecretsMap[key] = clientSecret

		if serverSecret, ok := serverSecretsMap[key]; !ok {
			if serverSecret.Title == clientSecret.Title {
				clientSecret.Version += 1
			}
			secretsToCreate = append(secretsToCreate, clientSecret)
		}
	}

	if len(secretsToCreate) > 0 {
		_, err := s.keeper.CreateMany(ctx, &kpb.CreateBatchRequest{UserEmail: in.Email, Secrets: secretsToCreate})
		if err != nil {
			slog.Error("message", slog.String("func",
				common.DotJoin(ServiceName, "Process", "CreateMany")), slog.String("error", err.Error()))
			return nil, err
		}
	}

	var clientMissingSecrets []*kpb.Secret
	for _, serverSecret := range serverSecrets {
		key := fmt.Sprintf("%s:%s", serverSecret.Title, serverSecret.HashPayload)
		if _, ok := clientSecretsMap[key]; !ok {
			clientMissingSecrets = append(clientMissingSecrets, serverSecret)
		}
	}

	return &spb.SyncResponse{
		Email:   in.Email,
		Secrets: clientMissingSecrets,
		Success: true,
	}, nil
}

func (s *SyncService) getServerSecrets(ctx context.Context, email string) ([]*kpb.Secret, error) {
	resp, err := s.keeper.Fetch(ctx, &kpb.FetchRequest{UserEmail: email})
	if err != nil {
		slog.Error("message", slog.String("func",
			common.DotJoin(ServiceName, "getServerSecrets", "Fetch")), slog.String("error", err.Error()))
		return nil, err
	}
	return resp.GetSecrets(), nil
}
