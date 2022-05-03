package main

import (
	"context"
	"flag"
	pb "github.com/go-programming-tour-book/grpc-demo/proto"
	"google.golang.org/grpc"
	"io"
	"log"
)

var port string

func init() {
	flag.StringVar(&port, "p", "8000", "启动端口号")
	flag.Parse()
}

func SayHello(client pb.GreeterClient) error {
	resp, _ := client.SayHello(context.Background(), &pb.HelloRequest{Name: "eddycjy"})
	log.Printf("client.SayHello resp: %s", resp.Message)
	return nil
}

func SayList(client pb.GreeterClient, r *pb.HelloRequest) error {
	stream, _ := client.SayList(context.Background(), r)
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		log.Printf("resp: %v", resp)
	}

	return nil
}

func SayRecord(client pb.GreeterClient, r *pb.HelloRequest) error {
	stream, _ := client.SayRecord(context.Background())
	for i := 0; i < 6; i++ {
		_ = stream.Send(r)
	}

	recv, _ := stream.CloseAndRecv()
	log.Printf("recv: %v", recv)
	return nil
}

func SayRoute(client pb.GreeterClient, r *pb.HelloRequest) error {
	stream, _ := client.SayRoute(context.Background())
	for i := 0; i <= 6; i++ {
		_ = stream.Send(r)
		recv, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		log.Printf("resp: %v", recv)
	}
	_ = stream.CloseSend()
	return nil
}

func main() {
	// 构建连接对象
	conn, _ := grpc.Dial(":"+port, grpc.WithInsecure())
	defer conn.Close()
	// 创建客户端对象
	client := pb.NewGreeterClient(conn)

	// 发送一次
	_ = SayHello(client)

	// 发送一次，接收Stream
	_ = SayList(client, &pb.HelloRequest{Name: "test-receive-list"})

	// 发送Stream(多次发送，并多次接收响应)
	_ = SayRecord(client, &pb.HelloRequest{Name: "test-send-list"})

	// 发送响应均为Stream
	_ = SayRoute(client, &pb.HelloRequest{Name: "test-receive-send-list"})

}
