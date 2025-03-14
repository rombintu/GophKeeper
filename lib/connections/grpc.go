package connections

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
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

func (p *ConnPool) CleanUp() {
	p.mu.Lock()
	defer p.mu.Unlock()

	for addr, conn := range p.conns {
		if !p.checkHealth(conn) {
			if err := conn.Close(); err != nil {
				slog.Warn("failed close conn", slog.String("error", err.Error()))
			}
			delete(p.conns, addr)
		}
	}
}

func (p *ConnPool) checkHealth(conn *grpc.ClientConn) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	state := conn.GetState()
	if state == connectivity.Ready {
		return true
	}

	return conn.WaitForStateChange(ctx, state)
}
