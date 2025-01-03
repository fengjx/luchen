package main

import (
	"context"
	"log"

	"github.com/fengjx/luchen/env"

	"github.com/fengjx/luchen/example/helloworld/pb"
)

func init() {
	if env.IsDev() {
		// 可以设置环境变量 LUCHEN_ETCD_ADDRESS 指定 etcd 地址
		env.SetDefaultEtcdAddress([]string{"192.168.6.129:2379"})
	}
}

func main() {
	greeterClient := pb.NewGreeterService("greet")
	helloResp, err := greeterClient.SayHello(context.Background(), &pb.HelloReq{
		Name: "fengjx",
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Println(helloResp.Message)
}
