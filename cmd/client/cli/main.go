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

	store := storage.NewSecretManager("mem://", "client")
	man := cli.NewManager(store)
	app := cli.NewApp(man)

	if err := app.Cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}

}
