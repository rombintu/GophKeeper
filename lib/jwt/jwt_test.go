package jwt

import (
	"context"
	"encoding/base64"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	proto "github.com/rombintu/GophKeeper/internal/proto/auth"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func TestNewTokenAndVerifyToken(t *testing.T) {
	secret := "test-secret"
	user := &proto.User{Email: "test@example.com"}
	duration := time.Hour

	// Генерация токена
	token, err := NewToken(user, secret, duration)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Проверка токена
	claims, err := VerifyToken(token, secret)
	assert.NoError(t, err)
	assert.Equal(t, user.Email, claims["email"])
	assert.NotEmpty(t, claims["iat"])
	assert.NotEmpty(t, claims["exp"])

	// Создаем токен с истекшим сроком действия
	expiredDuration := -time.Hour // Отрицательная длительность для истечения срока
	expiredToken, err := NewToken(user, secret, expiredDuration)
	assert.NoError(t, err)

	// Проверка истечения срока действия токена
	_, err = VerifyToken(expiredToken, secret)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "token is expired")
}

func TestGenerateHMACSecret(t *testing.T) {
	secret, err := GenerateHMACSecret(32)
	assert.NoError(t, err)
	assert.NotEmpty(t, secret)

	// Проверка длины закодированного секрета
	decoded, err := base64.URLEncoding.DecodeString(secret)
	assert.NoError(t, err)
	assert.Equal(t, 32, len(decoded))
}

func TestVerifyTokenInterceptor(t *testing.T) {
	secret := "test-secret"
	user := &proto.User{Email: "test@example.com"}
	duration := time.Hour

	// Генерация токена
	token, err := NewToken(user, secret, duration)
	assert.NoError(t, err)

	// Создание интерцептора
	interceptor := VerifyTokenInterceptor(secret, []string{"/auth.AuthService/Login"})

	// Тест для исключенного метода
	ctx := metadata.NewIncomingContext(context.Background(), metadata.New(map[string]string{}))
	_, err = interceptor(ctx, nil, &grpc.UnaryServerInfo{FullMethod: "/auth.AuthService/Login"}, func(ctx context.Context, req interface{}) (interface{}, error) {
		return "success", nil
	})
	assert.NoError(t, err)

	// Тест для метода, требующего авторизации
	ctx = metadata.NewIncomingContext(context.Background(), metadata.New(map[string]string{"authorization": "Bearer " + token}))
	_, err = interceptor(ctx, nil, &grpc.UnaryServerInfo{FullMethod: "/auth.AuthService/SomeMethod"}, func(ctx context.Context, req interface{}) (interface{}, error) {
		// Проверка наличия claims в контексте
		claims := ctx.Value(UserClaimsKey).(jwt.MapClaims)
		assert.Equal(t, user.Email, claims["email"])
		return "success", nil
	})
	assert.NoError(t, err)

	// Тест для отсутствующего заголовка авторизации
	ctx = metadata.NewIncomingContext(context.Background(), metadata.New(map[string]string{}))
	_, err = interceptor(ctx, nil, &grpc.UnaryServerInfo{FullMethod: "/auth.AuthService/SomeMethod"}, func(ctx context.Context, req interface{}) (interface{}, error) {
		return nil, nil
	})
	assert.Error(t, err)
	assert.Equal(t, codes.Unauthenticated, status.Code(err))

	// Тест для неверного формата токена
	ctx = metadata.NewIncomingContext(context.Background(), metadata.New(map[string]string{"authorization": "InvalidToken"}))
	_, err = interceptor(ctx, nil, &grpc.UnaryServerInfo{FullMethod: "/auth.AuthService/SomeMethod"}, func(ctx context.Context, req interface{}) (interface{}, error) {
		return nil, nil
	})
	assert.Error(t, err)
	assert.Equal(t, codes.Unauthenticated, status.Code(err))

	// Тест для неверного токена
	ctx = metadata.NewIncomingContext(context.Background(), metadata.New(map[string]string{"authorization": "Bearer invalid-token"}))
	_, err = interceptor(ctx, nil, &grpc.UnaryServerInfo{FullMethod: "/auth.AuthService/SomeMethod"}, func(ctx context.Context, req interface{}) (interface{}, error) {
		return nil, nil
	})
	assert.Error(t, err)
	assert.Equal(t, codes.Unauthenticated, status.Code(err))
}
