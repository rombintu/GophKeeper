package common

import (
	"errors"
	"log/slog"
	"strings"

	apb "github.com/rombintu/GophKeeper/internal/proto/auth"
)

// TODO: Email validate
func UserValidate(user *apb.User) error {
	if user.GetEmail() == "" {
		return errors.New("email is required")
	}
	if user.GetKeyChecksum() == nil {
		return errors.New("fingerprint of keys is required")
	}
	return nil
}

func DotJoin(opts ...string) string {
	return strings.Join(opts, ".")
}

func ifEmptyOpt(opt string) string {
	if opt == "" {
		return "N/A"
	}
	return opt
}

func Version(buildDate, buildVersion, buildCommit, binary string) {
	slog.Info(
		"Init", slog.String("Binary", ifEmptyOpt(binary)),
		slog.String("Build version", ifEmptyOpt(buildVersion)),
		slog.String("Build date", ifEmptyOpt(buildDate)),
		slog.String("Build commit", ifEmptyOpt(buildCommit)),
	)
}
