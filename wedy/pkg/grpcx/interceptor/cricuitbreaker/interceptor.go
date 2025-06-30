package cricuitbreaker

import (
	"context"
	"github.com/go-kratos/aegis/circuitbreaker/sre"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type InterceptorBuilder struct {
	breaker sre.Breaker
}

func (i *InterceptorBuilder) BuildUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		err = i.breaker.Allow()
		if err == nil {
			resp, err = handler(ctx, req)
			if err == nil {
				i.breaker.MarkSuccess()
			} else {
				i.breaker.MarkFailed()
			}
			return
		} else {
			i.breaker.MarkFailed()
			return nil, status.Errorf(codes.Unavailable, "Service Unavailable")
		}
	}
}
