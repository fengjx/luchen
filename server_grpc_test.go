package luchen_test

import (
	"context"

	grpctransport "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc"

	pb "google.golang.org/grpc/examples/helloworld/helloworld"

	"github.com/fengjx/luchen"
)

func newHelloGRPCServer(serviceName, addr string) *luchen.GRPCServer {
	server := luchen.NewGRPCServer(
		serviceName,
		luchen.WithGRPCAddr(addr),
	)
	server.RegisterServer(func(grpcServer *grpc.Server) {
		pb.RegisterGreeterServer(grpcServer, newGreeterServer())
	})
	return server
}

func newGreeterServer() pb.GreeterServer {
	svr := &GreeterServer{}
	svr.sayHello = luchen.NewGRPCHandler(
		makeSayHelloEndpoint(),
		luchen.DecodePB[pb.HelloRequest],
		luchen.EncodePB[pb.HelloReply],
	)
	return svr
}

type GreeterServer struct {
	pb.UnimplementedGreeterServer
	sayHello grpctransport.Handler
}

func (s *GreeterServer) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
	_, resp, err := s.sayHello.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.HelloReply), nil
}
