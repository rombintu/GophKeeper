package cli

import (
	"log/slog"
	"os"

	proto "github.com/rombintu/GophKeeper/internal/proto/auth"
	"github.com/rombintu/GophKeeper/lib/crypto"
)

type Profile struct {
	user  *proto.User
	Token string
}

func NewProfile(keyPath string) *Profile {
	master, err := crypto.LoadPrivateKey(keyPath)
	if err != nil {
		slog.Error("load master key", slog.String("error", err.Error()))
		os.Exit(0)
	}

	user, err := crypto.GetProfile(master)
	if err != nil {
		slog.Error("load get profile", slog.String("error", err.Error()))
		os.Exit(0)
	}
	return &Profile{
		user: user,
	}
}
