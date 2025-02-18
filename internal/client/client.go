package client

import (
	"context"
	"errors"

	"github.com/rombintu/GophKeeper/internal/proto"
)

type PublicClient interface {
	BaseClient
	Registration() error
	Login() error
	GetToken() string
}

func NewPublicClient(cfg Config) PublicClient {
	return &SimpleClient{Config: cfg}
}

func (c *SimpleClient) CheckProfile() error {
	if c.Config.Profile.Email == "" {
		return errors.New("the Email address is not specified in the profile")
	}
	if c.Config.Profile.HexKeys == nil {
		return errors.New("the key was not found in the profile")
	}
	return nil
}

func (c *SimpleClient) Login() error {
	client := proto.NewAuthClient(c.Conn)
	resp, err := client.Login(context.Background(), &proto.UserRequest{User: c.Config.Profile})
	if err != nil {
		return err
	}
	c.Config.token = resp.Token
	return nil
}

func (c *SimpleClient) Registration() error {
	client := proto.NewAuthClient(c.Conn)
	resp, err := client.Register(context.Background(), &proto.UserRequest{User: c.Config.Profile})
	if err != nil {
		return err
	}
	c.Config.token = resp.Token
	return nil
}

func (c *SimpleClient) GetToken() string {
	return c.Config.token
}
