package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/go-programming-tour-book/grpc-sample/protobuf"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
	"net"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

/**

// 定义服务的接口以及请求参数与返回值(各发一条消息)
Search(context.Context, *SearchRequest) (*SearchResponse, error)
// 客户端发多个消息
ClientStreamSearch(SearchService_ClientStreamSearchServer) error
// 服务端返回多个消息
ServerStreamSearch(*SearchRequest, SearchService_ServerStreamSearchServer) error
// 双方都发送多条消息
StreamSearch(SearchService_StreamSearchServer) error

*/

type server struct {
	// 成员变量是接口，要实现方法
	protobuf.SearchServer
}

// Search implements protobuf.SearchServiceServer
func (s *server) Search(_ context.Context, in *protobuf.SearchRequest) (*protobuf.SearchResponse, error) {
	log.Printf("Received: %v", in.GetRequest())
	return &protobuf.SearchResponse{Response: "VernHe 冲 冲 冲"}, nil
}

func (s *server) ClientStreamSearch(stream protobuf.Search_ClientStreamSearchServer) error {
	// 循环接收客户端发送的消息
	for {
		// Recv()方法不断接收数据
		request, err := stream.Recv()
		// 或者err == io.EOF
		if err != nil {
			// 数据全部接收完，返回消息并关闭
			// SendAndClose()返回数据并关闭流
			return stream.SendAndClose(&protobuf.SearchResponse{Response: "【ClientStreamSearch】相应内容"})
		}
		log.Printf("【ClientStreamSearch】收到消息: %v\n", request)
	}
}

func (s *server) ServerStreamSearch(sr *protobuf.SearchRequest, stream protobuf.Search_ServerStreamSearchServer) error {
	request := sr.Request
	log.Printf("【ServerStream】收到消息: %v\n", request)

	// 服务器通过Send()方法向客户端发送多条消息
	for i := 0; i < 5; i++ {
		stream.Send(&protobuf.SearchResponse{Response: fmt.Sprintf("【ServerStream】消息%d", i)})
	}

	return nil
}

func (s *server) StreamSearch(stream protobuf.Search_StreamSearchServer) error {
	// 循环接收客户端发送的消息
	for {
		// Recv()方法不断接收数据
		request, err := stream.Recv()
		// 或者err == io.EOF
		if err != nil {
			// 数据全部接收完，返回消息并关闭
			// SendAndClose()返回数据并关闭流
			return stream.Send(&protobuf.SearchResponse{Response: "【StreamSearch】相应内容"})
		}
		log.Printf("【StreamSearch】收到消息: %v\n", request)
	}
}

func main() {
	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	// 连接凭证
	// certFile公钥 keyFile私钥
	creds, err := credentials.NewServerTLSFromFile("F:/我的文件/GoProjects/go-programming-tour-book/grpc-sample/ca/server.pem", "F:/我的文件/GoProjects/go-programming-tour-book/grpc-sample/ca/server.key")
	if err != nil {
		log.Fatalln("公钥私钥有误:", err)
	}
	s := grpc.NewServer(grpc.Creds(creds))
	protobuf.RegisterSearchServer(s, &server{})
	err = s.Serve(listen)
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
