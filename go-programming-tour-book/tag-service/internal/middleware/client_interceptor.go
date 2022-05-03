package middleware

import (
	"context"
	"google.golang.org/grpc"
	"time"
)

func defaultContextTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	var cancelFunc context.CancelFunc
	if _, ok := ctx.Deadline(); !ok {
		// 未设置超市时长，就使用默认的
		defaultTimeout := 60 * time.Second
		ctx, cancelFunc = context.WithTimeout(ctx, defaultTimeout)
	}
	return ctx, cancelFunc
}

// UnaryContextTimeout 超时控制
func UnaryContextTimeout() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {

		// 设置默认超市时长
		ctx, cancelFunc := defaultContextTimeout(ctx)
		if cancelFunc != nil {
			defer cancelFunc()
		}

		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
