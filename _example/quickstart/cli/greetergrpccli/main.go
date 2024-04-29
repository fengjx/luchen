package main

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/fengjx/luchen/example/quickstart/pb"
)

func main() {
	clientConn, err := grpc.Dial(
		"localhost:8088",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		panic(err)
	}
	greeterClient := pb.NewGreeterClient(clientConn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	helloResp, err := greeterClient.SayHello(ctx, &pb.HelloReq{
		Name: "fengjx",
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Println(helloResp.Message)
}
