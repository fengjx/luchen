package greet

import (
	"context"

	"github.com/fengjx/luchen/example/helloworld/pb"
)

func (h *GreeterHandlerImpl) SayHello(ctx context.Context, req *pb.HelloReq) (*pb.HelloResp, error) {
	msg := "hello: " + req.Name
	return &pb.HelloResp{
		Message: msg,
	}, nil
}
