package main

import (
	"context"
	"net/http"
	"time"

	"github.com/fengjx/go-halo/halo"
	"github.com/fengjx/luchen"
	"google.golang.org/grpc"

	"github.com/fengjx/luchen/example/helloworld/pb"
)

// grpc server 功能示例

func main() {

	httpSvr := luchen.NewHTTPServer(
		luchen.WithServerAddr(":8080"),
	)

	e := &GreeterEndpoint{}
	httpSvr.HandleFunc("/say-hello", func(w http.ResponseWriter, r *http.Request) {
		luchen.NewHTTPTransportServer(e.MakeSayHelloEndpoint()).ServeHTTP(w, r)
	})

	grpcSvr := luchen.NewGRPCServer(
		luchen.WithServiceName("helloworld"),
		luchen.WithServerAddr(":8088"),
	)
	grpcSvr.RegisterService(func(gs *grpc.Server) {
		RegisterGreeterHandler(gs)
	})
	luchen.Start(grpcSvr, httpSvr)

	halo.Wait(10 * time.Second)
}

func RegisterGreeterHandler(gs *grpc.Server, middlewares ...luchen.Middleware) {
	e := &GreeterEndpoint{}
	pb.RegisterGreeterHandler(gs, e, middlewares)
}

type GreeterEndpoint struct {
}

func (e *GreeterEndpoint) MakeSayHelloEndpoint() luchen.Endpoint {
	return luchen.MakeEndpoint[*pb.HelloReq, *pb.HelloResp](e.SayHello)
}

func (e *GreeterEndpoint) SayHello(ctx context.Context, req *pb.HelloReq) (*pb.HelloResp, error) {
	msg := "hello: " + req.Name
	return &pb.HelloResp{
		Message: msg,
	}, nil
}
