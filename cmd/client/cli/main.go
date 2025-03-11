package main

import (
	"context"
	"log"
	"os"

	"github.com/rombintu/GophKeeper/internal/client/cli"
	"github.com/rombintu/GophKeeper/internal/storage"
	"github.com/rombintu/GophKeeper/lib/logger"
)

func main() {
	logger.InitLogger("local")

	man := cli.NewManager()
	// keyPath, err := man.ConfigGet(context.Background(), "key-path")
	// if err != nil {
	// 	slog.Warn("failed get key-path. use config set")
	// 	os.Exit(0)
	// }

	// Нужно переделать, чтобы подгружался ключ из базы
	// Изменить порядок
	profile := cli.NewProfile("./profiles/private-key.asc")
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
	man.SetProfile(profile)
	if err := man.Configure(); err != nil {
		panic(err)
	}
	app := cli.NewApp(man)

	if err := app.Cmd.Run(ctx, os.Args); err != nil {
		log.Fatal(err)
	}

}
