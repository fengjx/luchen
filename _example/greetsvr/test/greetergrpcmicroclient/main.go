package main

import (
	"context"
	"log"

	"github.com/fengjx/luchen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/fengjx/luchen/example/greetsvr/pb"
)

func main() {
	grpcClient := luchen.GetGRPCClient(
		"rpc.greeter",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	greeterClient := pb.NewGreeterClient(grpcClient)
	helloResp, err := greeterClient.SayHello(context.Background(), &pb.HelloReq{
		Name: "fengjx",
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Println(helloResp.Message)
}
