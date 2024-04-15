package hello

import (
	"context"

	"github.com/go-kit/kit/endpoint"

	"github.com/fengjx/luchen"

	"github.com/fengjx/luchen/example/greetsvr/pb"
)

type endpoints struct {
}

func newEndpoints() *endpoints {
	return &endpoints{}
}

func (e *endpoints) MakeSayHelloEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		logger := luchen.Logger(ctx)
		logger.Info("greeter say hello")
		helloReq := request.(*pb.HelloReq)
		msg, err := GetInst().helloLogic.SayHello(ctx, helloReq.Name)
		if err != nil {
			return nil, err
		}
		return &pb.HelloResp{Message: msg}, nil
	}
}
