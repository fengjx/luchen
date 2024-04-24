package hello

import (
	"context"

	"github.com/fengjx/luchen/log"
	"github.com/go-kit/kit/endpoint"

	"github.com/fengjx/luchen/example/greetsvr/pb"
)

type endpoints struct {
}

func newEndpoints() *endpoints {
	return &endpoints{}
}

func (e *endpoints) MakeSayHelloEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		log.InfoCtx(ctx, "greeter say hello")
		helloReq := request.(*pb.HelloReq)
		msg, err := GetInst().helloLogic.SayHello(ctx, helloReq.Name)
		if err != nil {
			return nil, err
		}
		return &pb.HelloResp{Message: msg}, nil
	}
}
