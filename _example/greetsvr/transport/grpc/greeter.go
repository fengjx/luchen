package grpc

import (
	"context"

	"github.com/fengjx/luchen"
	grpctransport "github.com/go-kit/kit/transport/grpc"

	"github.com/fengjx/luchen/example/greetsvr/pb"
	"github.com/fengjx/luchen/example/greetsvr/service/hello"
)

type GreeterServer struct {
	pb.UnimplementedGreeterServer
	sayHello grpctransport.Handler
}

func newGreeterServer() pb.GreeterServer {
	svr := &GreeterServer{}
	svr.sayHello = luchen.NewGRPCHandler(
		hello.GetInst().Endpoints.MakeSayHelloEndpoint(),
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
