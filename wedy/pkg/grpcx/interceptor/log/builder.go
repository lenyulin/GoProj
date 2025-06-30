package log

import (
	"GoProj/wedy/pkg/logger"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"runtime"
	"time"
)

type InterceptorBuilder struct {
	l logger.LoggerV1
}

func (b *InterceptorBuilder) BuildUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		start := time.Now()
		event := "Normal"
		defer func() {
			cost := time.Since(start)
			if rec := recover(); rec != nil {
				switch re := rec.(type) {
				case error:
					err = re
				default:
					err = fmt.Errorf("%v", re)
				}
				event = "Recover"
				stack := make([]byte, 4096)
				stack = stack[:runtime.Stack(stack, true)]
				err = status.Errorf(codes.Internal, "panic err"+err.Error(), event, stack)
			}
			files := []logger.Field{
				logger.Int64("cost", cost.Milliseconds()),
				logger.String("event", event),
				logger.String("method", info.FullMethod),
			}
			st, _ := status.FromError(err)
			if st != nil {
				files = append(files, logger.String("code", st.Code().String()))
				files = append(files, logger.String("code_msg", st.Message()))
			}
			b.l.Info("gRPC Interceptor:", files...)
		}()
		resp, err = handler(ctx, req)
		return
	}
}
