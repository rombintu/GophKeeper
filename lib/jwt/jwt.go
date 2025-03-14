package jwt

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	proto "github.com/rombintu/GophKeeper/internal/proto/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// Объявляем тип-обёртку для ключа
type contextKey string

// Определяем конкретный ключ с этим типом
const (
	UserClaimsKey contextKey = "userClaims"
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

// VerifyToken проверяет валидность JWT токена и возвращает claims
func VerifyToken(tokenString, secret string) (jwt.MapClaims, error) {
	// Парсим токен с проверкой подписи
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Проверяем метод подписи
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("token parsing failed: %w", err)
	}

	// Проверяем валидность токена
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Дополнительная проверка срока действия
		exp, ok := claims["exp"].(float64)
		if !ok {
			return nil, fmt.Errorf("invalid expiration time")
		}

		if time.Now().Unix() > int64(exp) {
			return nil, fmt.Errorf("token expired")
		}

		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// VerifyTokenInterceptor создает интерцептор для проверки JWT токена
//
// Пример исключения методов
//
//	"/auth.AuthService/Login"
//	"/auth.AuthService/Register"
func VerifyTokenInterceptor(secret string, excludedMethods []string) grpc.UnaryServerInterceptor {
	excluded := make(map[string]struct{})
	for _, m := range excludedMethods {
		excluded[m] = struct{}{}
	}

	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// Пропускаем проверку для исключенных методов
		if _, ok := excluded[info.FullMethod]; ok {
			return handler(ctx, req)
		}

		// Извлекаем метаданные из контекста
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "missing metadata")
		}

		// Получаем токен из заголовков
		authHeaders := md.Get("authorization")
		if len(authHeaders) == 0 {
			return nil, status.Error(codes.Unauthenticated, "missing authorization header")
		}

		// Извлекаем токен из формата "Bearer <token>"
		tokenParts := strings.Split(authHeaders[0], " ")
		if len(tokenParts) != 2 || strings.ToLower(tokenParts[0]) != "bearer" {
			return nil, status.Error(codes.Unauthenticated, "invalid token format")
		}

		tokenString := tokenParts[1]

		// Проверяем валидность токена
		claims, err := VerifyToken(tokenString, secret)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, fmt.Sprintf("invalid token: %v", err))
		}

		// Добавляем claims в контекст для использования в обработчике
		ctx = context.WithValue(ctx, UserClaimsKey, claims)

		return handler(ctx, req)
	}
}
