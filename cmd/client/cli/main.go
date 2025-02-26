package main

import (
	"context"
	"flag"
	"log/slog"
	"os"

	"github.com/rombintu/GophKeeper/internal/client"
	"github.com/rombintu/GophKeeper/internal/proto"
	"github.com/rombintu/GophKeeper/lib/connections"
	"github.com/rombintu/GophKeeper/lib/crypto"
	"github.com/rombintu/GophKeeper/lib/logger"
)

func main() {
	logger.InitLogger("local")
	// TODO: сделать чтобы запоминал ввод ключа и адреса
	privateKeyPath := flag.String("key", "", "Path to master GPG key")
	addressAuth := flag.String("address", "localhost:3201", "Address to AuthService")
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

	// Создание пула соединений, из него создаются клиенты
	clientPool := client.NewClientPool(connections.NewConnPool())
	authClient, err := clientPool.NewAuthClient(*addressAuth)
	if err != nil {
		slog.Error("get connection to auth service", slog.String("error", err.Error()))
		os.Exit(0)
	}
	ctx := context.Background()
	switch *action {
	case "profile", "p", "info":
		slog.Debug("client info", slog.String("email", user.GetEmail()), slog.String("fingerprint", string(user.GetHexKeys())))
	case "registration", "reg":
		reps, err := authClient.Register(ctx, &proto.RegisterRequest{User: user})
		if err != nil {
			slog.Error("registration", slog.String("error", err.Error()))
			return
		}
		slog.Debug("registration", slog.String("token", reps.GetToken()))
	case "login":
		reps, err := authClient.Login(ctx, &proto.LoginRequest{User: user})
		if err != nil {
			slog.Error("login", slog.String("error", err.Error()))
			return
		}
		slog.Debug("login", slog.String("token", reps.GetToken()))
	}
}
