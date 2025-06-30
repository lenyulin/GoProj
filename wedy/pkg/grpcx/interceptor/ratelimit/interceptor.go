package ratelimit

import (
	"GoProj/wedy/pkg/limiter"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strings"
)

type InterceptorBuilder struct {
	limiter limiter.Limiter
	key     string
}

// key: limiter:interceptoer-service
func NewInterceptorBuilder(limiter limiter.Limiter, key string) *InterceptorBuilder {
	return &InterceptorBuilder{limiter: limiter, key: key}
}

func (b *InterceptorBuilder) BuildUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		limited, err := b.limiter.Limit(ctx, b.key)
		if err != nil {
			return nil, status.Errorf(codes.ResourceExhausted, "rate limit exceeded")
		}
		if limited {
			return nil, status.Errorf(codes.ResourceExhausted, "rate limit exceeded")
		}
		return handler(ctx, req)
	}
}
func (b *InterceptorBuilder) BuildUnaryServerInterceptorService() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		limited, err := b.limiter.Limit(ctx, b.key)
		if strings.HasPrefix(info.FullMethod, "/UserService") {

			if err != nil {
				return nil, status.Errorf(codes.ResourceExhausted, "rate limit exceeded")
			}
			if limited {
				return nil, status.Errorf(codes.ResourceExhausted, "rate limit exceeded")
			}
		}
		return handler(ctx, req)
	}
}
