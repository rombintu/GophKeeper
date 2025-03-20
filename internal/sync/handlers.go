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

		_, keyExists := serverSecretsMap[key]
		sameTitleExists := false
		var maxVersion int64 = 0

		for _, serverSecret := range serverSecrets {
			if serverSecret.Title == clientSecret.Title {
				sameTitleExists = true
				if serverSecret.Version > maxVersion {
					maxVersion = serverSecret.Version
				}
			}
		}

		if !keyExists && sameTitleExists {
			newSecret := cloneSecret(clientSecret)
			newSecret.Version = maxVersion + 1
			secretsToCreate = append(secretsToCreate, newSecret)
		} else if !keyExists {
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

func cloneSecret(secret *kpb.Secret) *kpb.Secret {
	return &kpb.Secret{
		Title:       secret.Title,
		SecretType:  secret.SecretType,
		UserEmail:   secret.UserEmail,
		CreatedAt:   secret.CreatedAt,
		Version:     secret.Version,
		HashPayload: secret.HashPayload,
		Payload:     append([]byte(nil), secret.Payload...),
	}
}
