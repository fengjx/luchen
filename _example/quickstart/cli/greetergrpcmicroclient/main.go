package main

import (
	"context"
	"log"

	"github.com/fengjx/luchen"
	"github.com/fengjx/luchen/env"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/fengjx/luchen/example/quickstart/connom/config"
	"github.com/fengjx/luchen/example/quickstart/pb"
)

func init() {
	if env.IsLocal() {
		// 可以设置环境变量 LUCHEN_ETCD_ADDRESS 指定 etcd 地址
		env.SetDefaultEtcdAddress([]string{"host.etcd.dev:2379"})
	}
}

func main() {
	sname := config.GetConfig().Server.GRPC.ServerName
	grpcClient := luchen.GetGRPCClient(
		sname,
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
