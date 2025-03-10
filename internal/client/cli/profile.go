package cli

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/user"
	"path"

	"github.com/ProtonMail/go-crypto/openpgp"
	proto "github.com/rombintu/GophKeeper/internal/proto/auth"
	"github.com/rombintu/GophKeeper/lib/crypto"
)

type Profile struct {
	user   *proto.User
	Token  string
	key    openpgp.EntityList
	dbPath string
}

func NewProfile(keyPath string) *Profile {
	master, err := crypto.LoadPrivateKey(keyPath)
	if err != nil {
		slog.Info("load master key", slog.String("error", err.Error()))
		os.Exit(0)
	}

	user, err := crypto.GetProfile(master)
	if err != nil {
		slog.Info("load get profile", slog.String("error", err.Error()))
		os.Exit(0)
	}

	homeDirPath, err := getHomeDir()
	if err != nil {
		slog.Info("failed get homedir", slog.String("error", err.Error()))
		os.Exit(0)
	}

	projectDirPath := path.Join(homeDirPath, appDirName)
	if _, err := os.Stat(projectDirPath); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			if err := os.Mkdir(projectDirPath, os.ModePerm); err != nil {
				slog.Info("failed make project dir", slog.String("error", err.Error()))
				os.Exit(0)
			}
		} else {
			slog.Info(err.Error())
			os.Exit(0)
		}
	}
	dbPath := path.Join(projectDirPath, "bolt.db")
	if _, err := os.Stat(dbPath); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			slog.Warn("database not found... create once", slog.String("path", dbPath))
		} else {
			slog.Info("error get stat database", slog.String("error", err.Error()))
			os.Exit(0)
		}
	}

	return &Profile{
		user:   user,
		key:    master,
		dbPath: dbPath,
	}
}

func (p *Profile) GetKey() openpgp.EntityList {
	return p.key
}

func (p *Profile) GetDriverPath() string {
	return fmt.Sprintf("bolt://%s", p.dbPath)
}

func getHomeDir() (string, error) {
	user, err := user.Current()
	if err != nil {
		return "", err
	}
	return user.HomeDir, nil
}
