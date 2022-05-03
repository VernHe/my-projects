package main

import (
	"context"
	"github.com/go-programming-tour-book/tag-service/internal/middleware"
	"github.com/go-programming-tour-book/tag-service/proto"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"log"
)

type Auth struct {
	AppKey    string
	AppSecret string
}

func (a *Auth) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{"app_key": a.AppKey, "app_secret": a.AppSecret}, nil
}

func (a *Auth) RequireTransportSecurity() bool {
	return false
}

/**
grpc客户端，调用grpc服务端
*/
func main() {
	ctx := context.Background()
	// 发送gprc请求时带上metadata
	md := metadata.New(map[string]string{"go": "programming", "tour": "book"})
	outgoingContext := metadata.NewOutgoingContext(ctx, md)

	auth := &Auth{
		AppKey:    "go-programming-tour-book",
		AppSecret: "eddycjy",
	}
	option := []grpc.DialOption{grpc.WithPerRPCCredentials(auth)}
	conn := GetClientConn(outgoingContext, "127.0.0.1:9999", option)
	defer conn.Close()

	client := proto.NewTagServiceClient(conn)
	reply, err := client.GetList(ctx, &proto.GetTagListRequest{Name: "Java"})
	if err != nil {
		log.Fatalf("err: %v\n", err)
	}

	log.Printf("reply: %v \n", reply)
}

func GetClientConn(ctx context.Context, target string, opts []grpc.DialOption) *grpc.ClientConn {
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	opts = append(opts, grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(
		// 客户端超时控制
		middleware.UnaryContextTimeout(),
		// grpc自带的重试拦截器
		grpc_retry.UnaryClientInterceptor(
			// 最大重试次数
			grpc_retry.WithMax(2),
			// 需要重试的状态码
			grpc_retry.WithCodes(
				codes.Unknown,
				codes.Internal,
				codes.DeadlineExceeded,
			),
		),
	)))
	opts = append(opts, grpc.WithStreamInterceptor(grpc_middleware.ChainStreamClient()))
	conn, _ := grpc.DialContext(ctx, target, opts...)
	return conn
}
