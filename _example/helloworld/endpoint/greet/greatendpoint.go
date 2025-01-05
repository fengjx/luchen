package greet

import (
	"context"
	"fmt"
	"github.com/fengjx/luchen"
	"github.com/fengjx/luchen/example/helloworld/pb"
)

func RegisterGreeterGRPCHandler(gs *luchen.GRPCServer) {
	pb.RegisterGreeterGRPCHandler(gs, GreeterEndpointImpl)
}

func RegisterGreeterHTTPHandler(hs *luchen.HTTPServer) {
	pb.RegisterGreeterHTTPHandler(hs, GreeterEndpointImpl)
}

var GreeterEndpointImpl = &GreeterEndpoint{
	handler: &GreeterHandlerImpl{},
}

type GreeterHandlerImpl struct {
}

type GreeterEndpoint struct {
	handler pb.GreeterHandler
}

func (e *GreeterEndpoint) SayHelloEndpoint() luchen.Endpoint {
	fn := func(ctx context.Context, request any) (any, error) {
		req, ok := request.(*pb.HelloReq)
		if !ok {
			return nil, fmt.Errorf("invalid request type: %T", request)
		}
		return e.handler.SayHello(ctx, req)
	}
	return fn
}
