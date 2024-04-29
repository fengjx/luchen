package endpoint

import (
	"context"

	"github.com/fengjx/luchen/log"
	"github.com/go-kit/kit/endpoint"

	"github.com/fengjx/luchen/example/quickstart/logic/calc/internal/service"
	"github.com/fengjx/luchen/example/quickstart/pb"
)

var calcEdp = &calcEndpoint{}

type calcEndpoint struct {
}

func (e *calcEndpoint) makeAddEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		log.InfoCtx(ctx, "calc add")
		req := request.(*pb.AddReq)
		result, err := service.CalcSvc.Add(ctx, req.A, req.B)
		if err != nil {
			return nil, err
		}
		return &pb.AddResp{Result: result}, nil
	}
}
