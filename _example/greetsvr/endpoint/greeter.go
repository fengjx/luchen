package endpoint

import (
	"context"

	"github.com/go-kit/kit/endpoint"

	"github.com/fengjx/luchen"
	"github.com/fengjx/luchen/_example/greetsvr/logic"
	"github.com/fengjx/luchen/_example/greetsvr/pb"
)

type greeterEndpoints struct {
}

func newGreeterEndpoints() *greeterEndpoints {
	return &greeterEndpoints{}
}

func (e *greeterEndpoints) MakeSayHelloEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		logger := luchen.Logger(ctx)
		logger.Info("greeter say hello")
		helloReq := request.(*pb.HelloReq)
		msg, err := logic.GetInst().HelloLogic.SayHello(ctx, helloReq.Name)
		if err != nil {
			return nil, err
		}
		return &pb.HelloResp{Message: msg}, nil
	}
}
