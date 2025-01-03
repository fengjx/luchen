package main

import (
	"context"
	"log"

	"github.com/fengjx/luchen/env"

	"github.com/fengjx/luchen/example/greet/pb"
)

func init() {
	if env.IsLocal() {
		// 可以设置环境变量 LUCHEN_ETCD_ADDRESS 指定 etcd 地址
		env.SetDefaultEtcdAddress([]string{"host.etcd.dev:2379"})
	}
}

func main() {
	greeterClient := pb.NewGreeterService("featgrpc")
	helloResp, err := greeterClient.SayHello(context.Background(), &pb.HelloReq{
		Name: "fengjx",
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Println(helloResp.Message)
}
