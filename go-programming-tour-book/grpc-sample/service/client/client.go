package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/go-programming-tour-book/grpc-sample/protobuf"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
	name = flag.String("name", "defaultName", "Name to greet")
)

func main() {
	//testClientStream()
	//testServerStream()
	testStream()
}

func test1() {
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v\n", err)
	}
	defer conn.Close()

	// 创建客户端
	client := protobuf.NewSearchClient(conn)
	req := &protobuf.SearchRequest{Request: "VernHe GO GO GO"}
	// 普通发送消息
	clientStream, err := client.ClientStreamSearch(context.Background())
	if err != nil {
		log.Fatalln("普通消息", err)
	}
	clientStream.Send(req)
	searchResponse, err := clientStream.CloseAndRecv()
	if err != nil {
		log.Fatalln("普通消息", err)
	}
	log.Printf("Response: %v\n", searchResponse.Response)

}

func test() {
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v\n", err)
	}
	defer conn.Close()

	// 创建客户端
	client := protobuf.NewSearchClient(conn)

	// 客户端发送多个请求
	clientStreamSearchClient, err := client.ClientStreamSearch(context.Background())
	for i := 0; i < 5; i++ {
		clientStreamSearchClient.Send(&protobuf.SearchRequest{Request: fmt.Sprintf("【ClientStream】request%d", i)})
	}
	searchResponse, err := clientStreamSearchClient.CloseAndRecv()
	if err != nil {
		log.Fatalln("ClientStream close 错误", err)
	}
	log.Printf("【ClientStream】response: %v", searchResponse.Response)

	// 服务端返回多个消息
	serverStreamSearchClient, err := client.ServerStreamSearch(context.Background(), &protobuf.SearchRequest{Request: "【ServerStream】request"})
	if err != nil {
		log.Fatalln("ServerStream消息", err)
	}
	searchResponse, err = serverStreamSearchClient.Recv()
	for err == nil {
		log.Printf("【ServerStream】response: %v", searchResponse.Response)
	}
	serverStreamSearchClient.CloseSend()

	// 发送多个请求，返回多个消息
	streamSearchClient, err := client.StreamSearch(context.Background())
	if err != nil {
		log.Fatalln("StreamSearch消息", err)
	}
	for i := 0; i < 5; i++ {
		streamSearchClient.Send(&protobuf.SearchRequest{Request: fmt.Sprintf("【StreamSearch】request%d", i)})
	}
	streamSearchClient.CloseSend()
	searchResponse, err = streamSearchClient.Recv()
	for err == nil {
		log.Printf("【ServerStream】response: %v", searchResponse.Response)
	}
}

func ConnectToServer() *grpc.ClientConn {
	creds, err := credentials.NewClientTLSFromFile("F:/我的文件/GoProjects/go-programming-tour-book/grpc-sample/ca/server.pem", "v.com")
	if err != nil {
		log.Fatalf("Failed to create TLS credentials %v", err)
	}
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("connect faild: %v\n", err)
	}
	return conn
}

func testClientStream() {
	conn := ConnectToServer()
	defer conn.Close()
	// 创建客户端
	client := protobuf.NewSearchClient(conn)

	req := &protobuf.SearchRequest{Request: "【ClientStream】request"}
	clientStream, err := client.ClientStreamSearch(context.Background())
	if err != nil {
		log.Fatalln("clientStream err:", err)
	}
	// 发送多个请求
	clientStream.Send(req)
	clientStream.Send(req)
	clientStream.Send(req)
	searchResponse, err := clientStream.CloseAndRecv()
	if err != nil {
		log.Fatalln("普通消息", err)
	}
	// 接收响应
	log.Printf("Response: %v\n", searchResponse.Response)

}

func testServerStream() {
	conn := ConnectToServer()
	defer conn.Close()
	// 创建客户端
	client := protobuf.NewSearchClient(conn)

	req := &protobuf.SearchRequest{Request: "【ServerStream】request"}
	serverStreamSearchClient, err := client.ServerStreamSearch(context.Background(), req)
	if err != nil {
		log.Fatalln("serverStream err:", err)
	}
	response, err := serverStreamSearchClient.Recv()
	for err == nil {
		log.Printf("【ServerStream】response: %v", response.Response)
		response, err = serverStreamSearchClient.Recv()
	}
}

func testStream() {
	conn := ConnectToServer()
	defer conn.Close()
	// 创建客户端
	client := protobuf.NewSearchClient(conn)

	// 发送多个请求，返回多个消息
	streamSearchClient, err := client.StreamSearch(context.Background())
	if err != nil {
		log.Fatalln("stream", err)
	}
	for i := 0; i < 5; i++ {
		streamSearchClient.Send(&protobuf.SearchRequest{Request: fmt.Sprintf("【StreamSearch】request%d", i)})
	}
	streamSearchClient.CloseSend()
	searchResponse, err := streamSearchClient.Recv()
	for err == nil {
		log.Printf("【ServerStream】response: %v", searchResponse.Response)
		searchResponse, err = streamSearchClient.Recv()
	}
}
