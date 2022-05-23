package main

import (
	"fmt"
	"github.com/go-programming-tour-book/grpc-sample/protobuf"
	"github.com/golang/protobuf/proto"
	"log"
)

func main() {
	testWrite()
}

func testStruct() {
	p := protobuf.Person{
		Id:    1234,
		Name:  "John Doe",
		Email: "jdoe@example.com",
		Phones: []*protobuf.Person_PhoneNumber{
			{Number: "555-4321", Type: protobuf.Person_HOME},
		},
	}
	_ = p
}

func testWrite() {
	book := &protobuf.AddressBook{
		People: make([]*protobuf.Person, 10),
	}

	book.People[0] = &protobuf.Person{
		Id:    1234,
		Name:  "John Doe",
		Email: "jdoe@example.com",
		Phones: []*protobuf.Person_PhoneNumber{
			{Number: "555-4321", Type: protobuf.Person_HOME},
		},
	}

	// 将地址book写入磁盘
	out, err := proto.Marshal(book)
	if err != nil {
		log.Fatalln(err)
	}
	// 中间忽略写磁盘操作

	readBook := &protobuf.AddressBook{}
	// 从磁盘读取
	err = proto.Unmarshal(out, readBook)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("readBook: %v\n", readBook)
}
