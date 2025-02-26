package jwt

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	proto "github.com/rombintu/GophKeeper/internal/proto/auth"
)

// NewToken creates new JWT token for given user
func NewToken(user *proto.User, secret string, duration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	// Добавляем в токен всю необходимую информацию
	claims := token.Claims.(jwt.MapClaims)
	claims["email"] = user.Email
	claims["iat"] = time.Now().Unix()
	claims["exp"] = time.Now().Add(duration).Unix()

	// Подписываем токен, используя секретный ключ приложения
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func GenerateHMACSecret(length int) (string, error) {
	// Создаем байтовый срез нужной длины
	key := make([]byte, length)
	_, err := rand.Read(key)
	if err != nil {
		return "", fmt.Errorf("failed to generate HMAC secret: %w", err)
	}

	// Кодируем в base64 для удобства хранения
	return base64.URLEncoding.EncodeToString(key), nil
}

// func TokenValidateInterceptor(secret string) grpc.UnaryServerInterceptor {
// 	return func(
// 		ctx context.Context,
// 		req interface{},
// 		info *grpc.UnaryServerInfo,
// 		handler grpc.UnaryHandler,
// 	) (interface{}, error) {
// 		if !limiter.Allow() {
// 			return nil, status.Error(
// 				codes.ResourceExhausted,
// 				"too many requests",
// 			)
// 		}
// 		return handler(ctx, req)
// 	}
// }
