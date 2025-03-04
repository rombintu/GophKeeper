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

	profile := cli.NewProfile("./profiles/private-key.asc")
	store := storage.NewClientManager(storage.DriverOpts{
		ServiceName: "client",
		DriverPath:  "bolt:///tmp/bolt.db",
		CryptoKey:   profile.GetKey(),
	})
	man := cli.NewManager(profile, store)
	man.Configure()
	ctx := context.Background()
	if err := store.Open(ctx); err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := store.Close(ctx); err != nil {
			log.Fatal(err)
		}
	}()
	app := cli.NewApp(man)

	if err := app.Cmd.Run(ctx, os.Args); err != nil {
		log.Fatal(err)
	}

}
