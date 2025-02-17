package main

import (
	"context"
	"log/slog"

	"github.com/rombintu/GophKeeper/internal/auth"
	"github.com/rombintu/GophKeeper/internal/proto"
	"github.com/rombintu/GophKeeper/lib/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	logger.InitLogger("env")
	conn, err := grpc.NewClient("localhost:3201", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		slog.Error("error dial to client",
			slog.String("from", "client"),
			slog.String("to", auth.ServiceName),
			slog.String("error", err.Error()),
		)
	}
	defer conn.Close()
	ctx := context.Background()
	client := proto.NewAuthClient(conn)
	userResponse, err := client.Login(ctx, &proto.UserRequest{User: &proto.User{Email: "email", HexKeys: []byte("123")}})
	if err != nil {
		slog.Error("message", slog.String("error", err.Error()))
	}
	slog.Warn("resp", slog.String("token", userResponse.GetToken()))
}
