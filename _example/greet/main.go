package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/fengjx/luchen"
	"google.golang.org/grpc"

	"github.com/fengjx/luchen/example/greet/pb"
)

// grpc server 功能示例

func main() {
	grpcSvr := luchen.NewGRPCServer(
		luchen.WithServiceName("featgrpc"),
		luchen.WithServerAddr(":8088"),
	)
	grpcSvr.RegisterService(func(gs *grpc.Server) {
		RegisterGreeterHandler(gs)
	})
	luchen.Start(grpcSvr)

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGKILL)

	<-quit
	luchen.Stop()
}

func RegisterGreeterHandler(gs *grpc.Server, middlewares ...luchen.Middleware) {
	e := &GreeterEndpoint{}
	pb.RegisterGreeterHandler(gs, e)
}

type GreeterEndpoint struct {
}

func (e *GreeterEndpoint) MakeSayHelloEndpoint(middlewares ...luchen.Middleware) luchen.Endpoint {
	return luchen.MakeEndpoint[*pb.HelloReq, *pb.HelloResp](e.SayHello, middlewares...)
}

func (e *GreeterEndpoint) SayHello(ctx context.Context, req *pb.HelloReq) (*pb.HelloResp, error) {
	msg := "hello: " + req.Name
	return &pb.HelloResp{
		Message: msg,
	}, nil
}
