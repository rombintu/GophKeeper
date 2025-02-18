package main

import (
	"flag"
	"log/slog"
	"os"

	"github.com/rombintu/GophKeeper/internal/client"
	"github.com/rombintu/GophKeeper/lib/crypto"
	"github.com/rombintu/GophKeeper/lib/logger"
)

func main() {
	logger.InitLogger("local")
	// TODO: сделать чтобы запоминал ввод ключа и адреса
	privateKeyPath := flag.String("key", "", "Path to master GPG key")
	address := flag.String("address", "localhost:3201", "Address to AuthService")
	action := flag.String("action", "profile", "Action")
	flag.Parse()

	master, err := crypto.LoadPrivateKey(*privateKeyPath)
	if err != nil {
		slog.Error("load master key", slog.String("error", err.Error()))
		os.Exit(0)
	}

	user, err := crypto.GetProfile(master)
	if err != nil {
		slog.Error("load get profile", slog.String("error", err.Error()))
		os.Exit(0)
	}

	cfg := client.Config{
		Address: *address,
		Profile: user,
	}
	cli := client.NewPublicClient(cfg)
	if err := cli.Connect(); err != nil {
		slog.Error("connect", slog.String("error", err.Error()))
		return
	}
	defer func() {
		if err := cli.Disconnect(); err != nil {
			slog.Error("disconnect", slog.String("error", err.Error()))
		}
	}()

	switch *action {
	case "profile", "p", "info":
		slog.Debug("client info", slog.String("email", user.GetEmail()), slog.String("fingerprint", string(user.GetHexKeys())))
	case "registration", "reg":
		if err := cli.Registration(); err != nil {
			slog.Error("registration", slog.String("error", err.Error()))
			return
		}
		slog.Debug("registration", slog.String("token", cli.GetToken()))
	case "login":
		if err := cli.Login(); err != nil {
			slog.Error("login", slog.String("error", err.Error()))
			return
		}
		slog.Debug("login", slog.String("token", cli.GetToken()))
	}
}
