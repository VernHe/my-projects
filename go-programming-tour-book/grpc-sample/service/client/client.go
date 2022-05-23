package main

import (
	"context"
	"flag"
	"github.com/go-programming-tour-book/grpc-sample/protobuf"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
	name = flag.String("name", "defaultName", "Name to greet")
)

func main() {
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v\n", err)
	}
	defer conn.Close()

	// 创建客户端
	client := protobuf.NewSearchServiceClient(conn)
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second)
	defer cancelFunc()
	req := &protobuf.SearchRequest{Request: "VernHe GO GO GO"}
	response, err := client.Search(ctx, req)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("Response: %v\n", response.GetResponse())

}
