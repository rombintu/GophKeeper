package client

import (
	apb "github.com/rombintu/GophKeeper/internal/proto/auth"
	kpb "github.com/rombintu/GophKeeper/internal/proto/keeper"
	"github.com/rombintu/GophKeeper/lib/connections"
)

type ClientPool struct {
	pool *connections.ConnPool
}

func NewClientPool(pool *connections.ConnPool) *ClientPool {
	return &ClientPool{
		pool: pool,
	}
}

func (c *ClientPool) NewAuthClient(addr string) (apb.AuthClient, error) {
	conn, err := c.pool.Get(addr)
	if err != nil {
		return nil, err
	}
	return apb.NewAuthClient(conn), nil
}

func (c *ClientPool) NewKeeperClient(addr string) (kpb.KeeperClient, error) {
	conn, err := c.pool.Get(addr)
	if err != nil {
		return nil, err
	}
	return kpb.NewKeeperClient(conn), nil
}
