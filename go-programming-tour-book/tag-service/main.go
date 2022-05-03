package main

import (
	"context"
	"encoding/json"
	"flag"
	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/go-programming-tour-book/tag-service/internal/middleware"
	"github.com/go-programming-tour-book/tag-service/pkg/swagger"
	"github.com/go-programming-tour-book/tag-service/proto"
	"github.com/go-programming-tour-book/tag-service/server"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"log"
	"net"
	"net/http"
	"path"
	"strings"
)

//http错误
type httpError struct {
	Code    int32  `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

/**

通过命令"grpcurl -plaintext -d '{"name":"Java"}' localhost:9999 proto.TagService.GetList"，发起grpc调用

*/

// 监听的端口
var (
	port string
)

func init() {
	flag.StringVar(&port, "grpc_port", "9999", "启动端口号")
	flag.Parse()
}

func main() {
	err := RunServer(port)
	if err != nil {
		log.Fatalf("Run Serve err: %v", err)
	}
}

func RunTCPServer(port string) (net.Listener, error) {
	return net.Listen("tcp", ":"+port)
}

// RunHttpServer 创建一个HTTP多路复用器
func runHttpServer() *http.ServeMux {
	serveMux := http.NewServeMux()
	serveMux.HandleFunc("/ping", func(writer http.ResponseWriter, request *http.Request) {
		_, _ = writer.Write([]byte("pong"))
	})

	// 映射静态资源
	prefix := "/swagger-ui/"
	fileServer := http.FileServer(&assetfs.AssetFS{
		Asset:    swagger.Asset,
		AssetDir: swagger.AssetDir,
		Prefix:   "third_party/swagger-ui",
	})
	serveMux.Handle(prefix, http.StripPrefix(prefix, fileServer))

	serveMux.HandleFunc("/swagger/", func(writer http.ResponseWriter, request *http.Request) {
		log.Println("/swagger/ 收到请求")
		if !strings.HasSuffix(request.URL.Path, "swagger.json") {
			// 响应404
			http.NotFound(writer, request)
			return
		}
		// 去除 '/swagger/' 前缀
		trimPrefix := strings.TrimPrefix(request.URL.Path, "/swagger/")
		// 拼接前缀 'proto/'
		trimPrefix = path.Join("proto", trimPrefix)
		writer.Header().Set("Access-Control-Allow-Origin", "*")
		writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		writer.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		// 响应资源
		// 使用命名文件或目录的内容回复请求,name 参数中提供的文件或目录
		http.ServeFile(writer, request, trimPrefix)
	})

	return serveMux
}

// RunGrpcServer 运行Grpc服务器
func runGrpcServer() *grpc.Server {
	option := []grpc.ServerOption{
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			// 日志拦截器
			middleware.AccessLog,
			middleware.ErrorLog,
			middleware.Recpvery,
		)),
	}
	// 创建一个grpcServer
	grpcServer := grpc.NewServer(option...)
	// 注册grpc服务
	proto.RegisterTagServiceServer(grpcServer, server.NewTagServer())
	reflection.Register(grpcServer)

	return grpcServer
}

// grpcHandlerFunc 根据Content-Type，将grpc的请求交由grpc的方法进行处理
func grpcHandlerFunc(grpcServer *grpc.Server, otherHandler http.Handler) http.Handler {
	return h2c.NewHandler(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.ProtoMajor == 2 && strings.Contains(request.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(writer, request)
		} else {
			otherHandler.ServeHTTP(writer, request)
		}
	}), &http2.Server{})
}

// grpcGateWayError 处理错误方法
func grpcGateWayError(_ context.Context, _ *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, _ *http.Request, err error) {
	// 检查是否为grpc
	s, ok := status.FromError(err)
	if !ok {
		// err为nil或者未知
		s = status.New(codes.Unknown, err.Error())
	}

	httpError := httpError{
		Code:    int32(s.Code()),
		Message: s.Message(),
	}
	details := s.Details()
	for _, detail := range details {
		if v, ok := detail.(*proto.Error); ok {
			//log.Printf("code: %d, msg %v\n", v.Code, v.Message)
			httpError.Code = v.Code
			httpError.Message = v.Message
		}
	}

	resp, _ := json.Marshal(httpError)
	//log.Printf("resp: %v\n", string(resp))
	w.Header().Set("Content-type", marshaler.ContentType())
	w.WriteHeader(runtime.HTTPStatusFromCode(s.Code()))
	_, _ = w.Write(resp)
}

// runGrpcGatewayServer 运行一个网关作用的多路复用器
func runGrpcGatewayServer() *runtime.ServeMux {
	endpoint := "0.0.0.0:" + port
	runtime.HTTPError = grpcGateWayError
	serveMux := runtime.NewServeMux()
	options := []grpc.DialOption{grpc.WithInsecure()}
	proto.RegisterTagServiceHandlerFromEndpoint(context.Background(), serveMux, endpoint, options)

	return serveMux
}

func RunServer(port string) error {
	// http多路复用器
	httpServerMux := runHttpServer()
	// grpcServer
	grpcServer := runGrpcServer()
	gatewayServer := runGrpcGatewayServer()

	httpServerMux.Handle("/", gatewayServer)
	return http.ListenAndServe(":"+port, grpcHandlerFunc(grpcServer, httpServerMux))
}
