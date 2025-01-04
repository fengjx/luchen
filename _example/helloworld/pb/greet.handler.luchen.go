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
	SayHelloEdnpointDefine() *luchen.EdnpointDefine
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

func RegisterGreeterGRPCHandler(gs *luchen.GRPCServer, e GreeterEndpoint) {
	impl := GreeterServiceImpl{
		sayHello: luchen.NewGRPCTransportServer(
			e.SayHelloEdnpointDefine(),
		),
	}
	RegisterGreeterServer(gs, impl)
}

func RegisterGreeterHTTPHandler(hs *luchen.HTTPServer, e GreeterEndpoint) {
	def := e.SayHelloEdnpointDefine()
	e := luchen.MakeEndpoint(e.SayHelloEdnpointDefine())
	dec := luchen.DecodeHTTPPbRequest[](ctx context.Context, req *http.Request)
	h := luchen.NewHTTPTransportServer(e, dec http.DecodeRequestFunc, enc http.EncodeResponseFunc)
	hs.Handle(def.Path, h)
}
