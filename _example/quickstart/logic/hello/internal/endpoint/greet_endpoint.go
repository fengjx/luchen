package endpoint

import (
	"context"

	"github.com/fengjx/luchen/log"
	"github.com/go-kit/kit/endpoint"

	"github.com/fengjx/luchen/example/quickstart/logic/hello/service"
	"github.com/fengjx/luchen/example/quickstart/pb"
)

var greetEdp = &greetEndpoints{}

type greetEndpoints struct {
}

func (e *greetEndpoints) makeSayHelloEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		log.InfoCtx(ctx, "greeter say hello")
		helloReq := request.(*pb.HelloReq)
		msg, err := service.GreetSvc.SayHi(ctx, helloReq.Name)
		if err != nil {
			return nil, err
		}
		return &pb.HelloResp{Message: msg}, nil
	}
}
