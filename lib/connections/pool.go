package connections

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type ConnPoolAdapter interface {
	Get(addr string) (*grpc.ClientConn, error)
	CleanUp()
}

type ConnPool struct {
	mu    sync.Mutex
	conns map[string]*grpc.ClientConn
}

func NewConnPool() ConnPoolAdapter {
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

func RateLimitInterceptor(limiter *rate.Limiter) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		if !limiter.Allow() {
			return nil, status.Error(
				codes.ResourceExhausted,
				"too many requests",
			)
		}
		return handler(ctx, req)
	}
}
