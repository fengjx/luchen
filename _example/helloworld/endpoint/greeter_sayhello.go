package endpoint

import (
	"context"
	"fmt"
	"net/http"

	"github.com/fengjx/luchen"
	"github.com/fengjx/luchen/example/helloworld/pbgreet"
)

func (e *GreeterEndpoint) SayHelloEndpoint() luchen.Endpoint {
	fn := func(ctx context.Context, request any) (any, error) {
		req, ok := request.(*pbgreet.HelloReq)
		if !ok {
			msg := fmt.Sprintf("invalid request type: %T", request)
			return nil, luchen.NewErrno(http.StatusBadRequest, msg)
		}
		return e.handler.SayHello(ctx, req)
	}
	return fn
}

// SayHello Sends a greeting
// http.path=/say-hello
func (h *GreeterHandlerImpl) SayHello(ctx context.Context, req *pbgreet.HelloReq) (*pbgreet.HelloResp, error) {
	msg := "hello: " + req.Name
	return &pbgreet.HelloResp{
		Message: msg,
	}, nil
}
