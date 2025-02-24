package client

import (
	"github.com/rombintu/GophKeeper/internal/proto"
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

func (c *ClientPool) NewAuthClient(addr string) (proto.AuthClient, error) {
	conn, err := c.pool.Get(addr)
	if err != nil {
		return nil, err
	}
	return proto.NewAuthClient(conn), nil
}
