package main

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/fengjx/go-halo/halo"
	"github.com/fengjx/luchen"

	"github.com/fengjx/luchen/env"
	"github.com/fengjx/luchen/example/helloworld/endpoint/greet"
	"github.com/fengjx/luchen/example/helloworld/pb"
)

func init() {
	if env.IsDev() {
		// 可以设置环境变量 LUCHEN_ETCD_ADDRESS 指定 etcd 地址
		env.SetDefaultEtcdAddress([]string{"192.168.6.121:2379"})
	}
}

// grpc server 功能示例

func main() {

	hs := luchen.NewHTTPServer(
		luchen.WithServiceName("http.helloworld"),
		luchen.WithServerAddr(":8080"),
	)

	gs := luchen.NewGRPCServer(
		luchen.WithServiceName("grpc.helloworld"),
		luchen.WithServerAddr(":8088"),
	)

	greet.RegisterGreeterGRPCHandler(gs)
	greet.RegisterGreeterHTTPHandler(hs)

	registrar := luchen.NewEtcdV3Registrar(
		hs,
		gs,
	)
	registrar.Register()
	// luchen.Start(gs, hs)
	halo.Wait(10 * time.Second)
}

type GreeterEndpoint struct {
}

func (e *GreeterEndpoint) SayHelloEndpointDefine() *luchen.EndpointDefine {
	def := &luchen.EndpointDefine{
		Name:    "Greet.SayHello",
		Path:    "/say-hello",
		ReqType: reflect.TypeOf(&pb.HelloReq{}),
		RspType: reflect.TypeOf(&pb.HelloResp{}),
		Endpoint: func(ctx context.Context, request any) (any, error) {
			req, ok := request.(*pb.HelloReq)
			if !ok {
				return nil, fmt.Errorf("invalid request type: %T", request)
			}
			return e.SayHello(ctx, req)
		},
	}
	return def
}

func (e *GreeterEndpoint) SayHello(ctx context.Context, req *pb.HelloReq) (*pb.HelloResp, error) {
	msg := "hello: " + req.Name
	return &pb.HelloResp{
		Message: msg,
	}, nil
}
