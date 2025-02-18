package client

import (
	"github.com/rombintu/GophKeeper/internal/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type BaseClient interface {
	Connect() error
	Disconnect() error
}

type Config struct {
	Address string
	Profile *proto.User
	token   string
}

// Реализация простого клиента
type SimpleClient struct {
	Conn   *grpc.ClientConn
	Config Config
}

func (c *SimpleClient) Connect() error {
	conn, err := grpc.NewClient(c.Config.Address,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	c.Conn = conn
	return nil
}

func (c *SimpleClient) Disconnect() error {
	return c.Conn.Close()
}
