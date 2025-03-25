package main

import (
	"context"
	"log"
	"log/slog"
	"os"

	"github.com/rombintu/GophKeeper/internal/client/cli"
	"github.com/rombintu/GophKeeper/internal/config"
	"github.com/rombintu/GophKeeper/internal/storage"
	"github.com/rombintu/GophKeeper/lib/logger"
)

func main() {
	logger.InitLogger("local")

	man := cli.NewManager()

	profile := cli.NewProfile()
	conf, err := config.NewClientConfig(profile.GetConfFilePath())
	if err != nil {
		log.Fatal(err)
	}
	if err := conf.Save(profile.GetConfFilePath()); err != nil {
		log.Fatal(err)
	}
	if err := profile.LoadKey(conf.KeyPath); err != nil {
		log.Fatal(err)
	}
	store := storage.NewClientManager(storage.DriverOpts{
		ServiceName: "client",
		DriverPath:  profile.GetDriverPath(),
		CryptoKey:   profile.GetKey(),
	})
	ctx := context.Background()
	if err := store.Open(ctx); err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := store.Close(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	man.SetStore(store)
	if err := man.Configure(); err != nil {
		log.Fatal(err)
	}

	if err := profile.LoadKey(conf.KeyPath); err != nil {
		slog.Warn("failed load key", slog.String("error", err.Error()))
	}

	man.SetProfile(profile)
	app := cli.NewApp(man)

	if err := app.Cmd.Run(ctx, os.Args); err != nil {
		log.Fatal(err)
	}

}
