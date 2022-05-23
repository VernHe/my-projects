package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/go-programming-tour-book/grpc-sample/protobuf"
	"google.golang.org/grpc"
	"log"
	"net"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

type server struct {
	protobuf.SearchServiceServer
}

// Search implements protobuf.SearchServiceServer
func (s *server) Search(_ context.Context, in *protobuf.SearchRequest) (*protobuf.SearchResponse, error) {
	log.Printf("Received: %v", in.GetRequest())
	return &protobuf.SearchResponse{Response: "VernHe 冲 冲 冲"}, nil
}

func main() {
	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	protobuf.RegisterSearchServiceServer(s, &server{})
	err = s.Serve(listen)
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
