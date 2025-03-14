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
	user     *proto.User
	key      openpgp.EntityList
	dbPath   string
	confFile string
}

func checkFileExists(pathFile string) {
	if _, err := os.Stat(pathFile); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			if _, err := os.Create(pathFile); err != nil {
				slog.Info("failed make file", slog.String("error", err.Error()))
				os.Exit(0)
			}
		} else {
			slog.Info(err.Error())
			os.Exit(0)
		}
	}
}

func NewProfile() *Profile {

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

	confFile := path.Join(projectDirPath, "config.json")
	checkFileExists(confFile)

	dbPath := path.Join(projectDirPath, "bolt.db")
	checkFileExists(dbPath)

	return &Profile{
		dbPath:   dbPath,
		confFile: confFile,
	}
}

func (p *Profile) GetConfFilePath() string {
	return p.confFile
}

func (p *Profile) LoadKey(keyPath string) error {
	if keyPath == "" {
		return fmt.Errorf("key-path is not set in file %s", p.confFile)
	}
	master, err := crypto.LoadPrivateKey(keyPath)
	if err != nil {
		return err
	}

	user, err := crypto.GetUserFromKey(master)
	if err != nil {
		return err
	}

	p.user = user
	p.key = master
	return nil
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
