package connections

import (
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ConnPool struct {
	mu    sync.Mutex
	conns map[string]*grpc.ClientConn
}

func NewConnPool() *ConnPool {
	return &ConnPool{
		conns: make(map[string]*grpc.ClientConn),
	}
}

func (p *ConnPool) Get(addr string) (*grpc.ClientConn, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if conn, ok := p.conns[addr]; ok {
		return conn, nil
	}

	conn, err := grpc.NewClient(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
	)
	if err != nil {
		return nil, err
	}
	p.conns[addr] = conn
	return conn, nil
}
