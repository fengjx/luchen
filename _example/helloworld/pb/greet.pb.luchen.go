package pb

import (
	context "context"

	"github.com/fengjx/luchen"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// NewGreeterService 返回一个 GreeterClient
func NewGreeterService(serverName string) GreeterClient {
	cli := luchen.GetGRPCClient(
		serverName,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	return NewGreeterClient(cli)
}

type GreeterHandler interface {
	SayHello(ctx context.Context, in *HelloReq) (*HelloResp, error)
}

type GreeterEndpoint interface {
	GreeterHandler
	MakeSayHelloEndpoint() luchen.Endpoint
}

type GreeterServiceImpl struct {
	UnimplementedGreeterServer
	sayHello grpctransport.Handler
}

func (s *GreeterServiceImpl) SayHelHlo(ctx context.Context, req *HelloReq) (*HelloResp, error) {
	_, resp, err := s.sayHello.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*HelloResp), nil
}

func RegisterGreeterHandler(gs *grpc.Server, e GreeterEndpoint, middlewares []luchen.Middleware, options ...grpctransport.ServerOption) {
	impl := GreeterServiceImpl{
		sayHello: luchen.NewGRPCTransportServer(
			luchen.EndpointChain(e.MakeSayHelloEndpoint(), middlewares...),
			options...,
		),
	}
	RegisterGreeterServer(gs, impl)
}
