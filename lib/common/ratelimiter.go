package common

import (
	"context"

	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

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
