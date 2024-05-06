package endpoint

import (
	"context"

	"github.com/fengjx/luchen/log"
	"github.com/go-kit/kit/endpoint"

	"github.com/fengjx/luchen/example/httponly/logic/hello/internal/protocol"
	"github.com/fengjx/luchen/example/httponly/logic/hello/internal/service"
)

var greetEdp = &greetEndpoint{}

type greetEndpoint struct {
}

func (e *greetEndpoint) makeSayHelloEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		log.InfoCtx(ctx, "greeter say hello")
		helloReq := request.(*protocol.HelloReq)
		msg, err := service.GreetSvc.SayHi(ctx, helloReq.Name)
		if err != nil {
			return nil, err
		}
		return &protocol.HelloResp{Message: msg}, nil
	}
}
