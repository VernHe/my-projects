package main

import (
	"context"
	"flag"
	pb "github.com/go-programming-tour-book/grpc-demo/proto"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
)

var port string

func init() {
	flag.StringVar(&port, "p", "8000", "启动端口号")
	flag.Parse()
}

type GreeterServer struct{}

func (s *GreeterServer) SayRecord(stream pb.Greeter_SayRecordServer) error {
	for {
		// 不断的去接收多次消息
		recv, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.HelloReplay{Message: "say.record"})
		}
		if err != nil {
			return err
		}
		log.Printf("recv: %v", recv)
	}
}

func (s *GreeterServer) SayHello(ctx context.Context, r *pb.HelloRequest) (*pb.HelloReplay, error) {
	return &pb.HelloReplay{Message: "hello.world " + r.GetName()}, nil
}

func (s *GreeterServer) SayList(r *pb.HelloRequest, stream pb.Greeter_SayListServer) error {
	// 连续多次发消息
	for n := 0; n <= 6; n++ {
		_ = stream.Send(&pb.HelloReplay{Message: "hello.list"})
		log.Printf("receive request: %v", r.GetName())
	}

	return nil
}

func (s *GreeterServer) SayRoute(stream pb.Greeter_SayRouteServer) error {
	for {
		stream.Send(&pb.HelloReplay{Message: "say.route"})

		recv, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		log.Printf("recv: %v", recv)
	}
}

func main() {
	server := grpc.NewServer()
	pb.RegisterGreeterServer(server, &GreeterServer{})
	lis, _ := net.Listen("tcp", ":"+port)
	server.Serve(lis)
}
