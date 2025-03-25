package connections

import (
	"context"
	"testing"
	"time"

	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestConnPool_Get(t *testing.T) {
	pool := NewConnPool().(*ConnPool)
	addr := "localhost:50051"

	// Test creating a new connection
	conn, err := pool.Get(addr)
	if err != nil {
		t.Fatalf("Failed to get connection: %v", err)
	}
	if conn == nil {
		t.Error("Expected a valid connection, got nil")
	}

	// Test retrieving an existing connection
	conn2, err := pool.Get(addr)
	if err != nil {
		t.Fatalf("Failed to get existing connection: %v", err)
	}
	if conn != conn2 {
		t.Error("Expected the same connection, got different connections")
	}
}

// func TestConnPool_CleanUp(t *testing.T) {
// 	pool := NewConnPool().(*ConnPool)
// 	addr := "localhost:50051"

// 	// Создаем новое соединение и добавляем его в пул
// 	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
// 	if err != nil {
// 		t.Fatalf("Failed to create connection: %v", err)
// 	}
// 	pool.mu.Lock()
// 	pool.conns[addr] = conn
// 	pool.mu.Unlock()

// 	// Даем соединению время для перехода в состояние Ready
// 	time.Sleep(2 * time.Second)

// 	// Проверяем, что соединение остается в пуле, если оно здоровое
// 	pool.CleanUp()
// 	if _, ok := pool.conns[addr]; !ok {
// 		t.Error("Expected connection to remain in pool, but it was removed")
// 	}

// 	// Закрываем соединение и проверяем, что оно удаляется из пула
// 	conn.Close()
// 	pool.CleanUp()
// 	if _, ok := pool.conns[addr]; ok {
// 		t.Error("Expected connection to be removed from pool, but it remained")
// 	}
// }
// func TestConnPool_checkHealth(t *testing.T) {
// 	pool := NewConnPool().(*ConnPool)
// 	addr := "localhost:50051"

// 	// Создаем новое соединение
// 	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
// 	if err != nil {
// 		t.Fatalf("Failed to create connection: %v", err)
// 	}

// 	// Даем соединению время для перехода в состояние Ready
// 	time.Sleep(2 * time.Second)

// 	// Проверяем, что соединение здоровое
// 	if !pool.checkHealth(conn) {
// 		t.Error("Expected connection to be healthy, but it was not")
// 	}

// 	// Закрываем соединение и проверяем, что оно нездоровое
// 	conn.Close()
// 	if pool.checkHealth(conn) {
// 		t.Error("Expected connection to be unhealthy, but it was healthy")
// 	}
// }

type mockHandler struct {
	called bool
}

func (m *mockHandler) handle(ctx context.Context, req interface{}) (interface{}, error) {
	m.called = true
	return "response", nil
}

func TestRateLimitInterceptor(t *testing.T) {
	limiter := rate.NewLimiter(rate.Every(time.Second), 1)
	interceptor := RateLimitInterceptor(limiter)

	handler := &mockHandler{}
	mockInfo := &grpc.UnaryServerInfo{}

	// Test allowed request
	_, err := interceptor(context.Background(), "request", mockInfo, handler.handle)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if !handler.called {
		t.Error("Expected handler to be called, but it was not")
	}

	// Test rate-limited request
	_, err = interceptor(context.Background(), "request", mockInfo, handler.handle)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
	if status.Code(err) != codes.ResourceExhausted {
		t.Errorf("Expected error code %v, got %v", codes.ResourceExhausted, status.Code(err))
	}
}
