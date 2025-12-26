package interceptors

import (
	"context"
	"fmt"
	"path"

	"github.com/kitanoyoru/kgym/internal/apps/user/pkg/metrics"
	"google.golang.org/grpc"
)

func MetricsUnaryServerInterceptor(prefix string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		method := path.Base(info.FullMethod)

		metrics.GlobalRegistry.
			GetMetric(fmt.Sprintf("%s.%s", prefix, method)).
			Counter.WithLabelValues().Inc()

		return handler(ctx, req)
	}
}

func MetricsStreamServerInterceptor(prefix string) grpc.StreamServerInterceptor {
	return func(srv any, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		method := path.Base(info.FullMethod)

		metrics.GlobalRegistry.
			GetMetric(fmt.Sprintf("%s.%s", prefix, method)).
			Counter.WithLabelValues().Inc()

		return handler(srv, stream)
	}
}
