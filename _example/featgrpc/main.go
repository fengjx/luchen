package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/fengjx/luchen"
	kitendpoint "github.com/go-kit/kit/endpoint"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc"

	"github.com/fengjx/luchen/example/featgrpc/pb"
)

// grpc server 功能示例

func main() {
	grpcSvr := luchen.NewGRPCServer(
		luchen.WithServiceName("featgrpc"),
		luchen.WithServerAddr(":8088"),
	)
	grpcSvr.RegisterService(func(gs *grpc.Server) {
		// 注册 grpc 服务
		pb.RegisterGreeterServer(gs, newGreeterServer())
	})
	luchen.Start(grpcSvr)

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGKILL)

	<-quit
	luchen.Stop()
}

type GreeterServer struct {
	pb.UnimplementedGreeterServer
	sayHello grpctransport.Handler
}

func newGreeterServer() pb.GreeterServer {
	svr := &GreeterServer{}
	svr.sayHello = luchen.NewGRPCTransportServer(
		makeSayHelloEndpoint(),
		luchen.DecodePB[pb.HelloReq],
		luchen.EncodePB[pb.HelloResp],
	)
	return svr
}

func (s *GreeterServer) SayHello(ctx context.Context, req *pb.HelloReq) (*pb.HelloResp, error) {
	_, resp, err := s.sayHello.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.HelloResp), nil
}

func (s *GreeterServer) decodeSayHello(_ context.Context, req interface{}) (interface{}, error) {
	helloReq := req.(*pb.HelloReq)
	return &pb.HelloReq{
		Name: helloReq.Name,
	}, nil
}

func (s *GreeterServer) encodeSayHello(_ context.Context, resp interface{}) (interface{}, error) {
	helloResp := resp.(*pb.HelloResp)
	return &pb.HelloResp{
		Message: helloResp.Message,
	}, nil
}

func makeSayHelloEndpoint() kitendpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		name := request.(string)
		response = "hello: " + name
		return
	}
}
