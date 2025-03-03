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

	store := storage.NewClientManager("bolt:///tmp/bolt.db")
	ctx := context.Background()
	if err := store.Open(ctx); err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := store.Close(ctx); err != nil {
			log.Fatal(err)
		}
	}()
	man := cli.NewManager(store)
	man.Configure()
	app := cli.NewApp(man)

	if err := app.Cmd.Run(ctx, os.Args); err != nil {
		log.Fatal(err)
	}

}
