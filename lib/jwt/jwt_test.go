package jwt

import (
	"context"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	apb "github.com/rombintu/GophKeeper/internal/proto/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func TestVerifyTokenInterceptor(t *testing.T) {
	// Генерация тестового токена
	token, _ := NewToken(&apb.User{Email: "test@example.com"}, "secret", time.Hour)

	// Создание контекста с метаданными
	md := metadata.Pairs("authorization", "Bearer "+token)
	ctx := metadata.NewIncomingContext(context.Background(), md)

	// Тестовый обработчик
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		claims, ok := ctx.Value("userClaims").(jwt.MapClaims)
		if !ok {
			return nil, status.Error(codes.Internal, "claims missing")
		}
		return claims["email"], nil
	}

	// Вызов интерцептора
	interceptor := VerifyTokenInterceptor("secret", nil)
	resp, err := interceptor(ctx, nil, &grpc.UnaryServerInfo{}, handler)

	// Проверки
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if resp != "test@example.com" {
		t.Errorf("Invalid email in response: %v", resp)
	}
}
