package pb

import (
	"context"

	grpctransport "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/fengjx/luchen"
	"reflect"
	"sync"
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
	SayHelloEndpoint() luchen.Endpoint
}

type GreeterServiceImpl struct {
	UnimplementedGreeterServer
	middlewares    []luchen.Middleware
	endpoint       GreeterEndpoint
	sayHelloDefine *luchen.EdnpointDefine
	sayHello       grpctransport.Handler
}

var (
	greeterServiceImplOnce = sync.Once{}
	greeterServiceImpl     *GreeterServiceImpl
)

func GetGreeterServiceImpl(e GreeterEndpoint, middlewares ...luchen.Middleware) *GreeterServiceImpl {
	greeterServiceImplOnce.Do(func() {
		greeterServiceImpl = newGreeterServiceImpl(e, middlewares...)
	})
	return greeterServiceImpl
}

func newGreeterServiceImpl(e GreeterEndpoint, middlewares ...luchen.Middleware) *GreeterServiceImpl {
	sayHelloDefine := &luchen.EdnpointDefine{
		Name:        "Greet.SayHello",
		Path:        "/say-hello",
		ReqType:     reflect.TypeOf(&HelloReq{}),
		RspType:     reflect.TypeOf(&HelloResp{}),
		Endpoint:    e.SayHelloEndpoint(),
		Middlewares: middlewares,
	}
	impl := &GreeterServiceImpl{
		endpoint:       e,
		sayHelloDefine: sayHelloDefine,
	}
	impl.sayHello = luchen.NewGRPCTransportServer(sayHelloDefine)
	return impl
}

func (s *GreeterServiceImpl) SayHello(ctx context.Context, req *HelloReq) (*HelloResp, error) {
	_, resp, err := s.sayHello.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*HelloResp), nil
}

// RegisterGreeterGRPCHandler 注册  GRPC 接口实现
func RegisterGreeterGRPCHandler(gs *luchen.GRPCServer, e GreeterEndpoint, middlewares ...luchen.Middleware) {
	impl := GetGreeterServiceImpl(e, middlewares...)
	RegisterGreeterServer(gs, impl)
}

// RegisterGreeterHTTPHandler 注册HTTP请求路由
func RegisterGreeterHTTPHandler(hs *luchen.HTTPServer, e GreeterEndpoint, middlewares ...luchen.Middleware) {
	impl := GetGreeterServiceImpl(e, middlewares...)
	if impl.sayHelloDefine.Path != "" {
		hs.Mux().Handle(impl.sayHelloDefine.Path, luchen.NewHTTPTransportServer(impl.sayHelloDefine))
	}
}
