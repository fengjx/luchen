package endpoint

import (
	"context"

	"github.com/fengjx/luchen"
	grpctransport "github.com/go-kit/kit/transport/grpc"

	"github.com/fengjx/luchen/example/quickstart/pb"
)

type calcServer struct {
	pb.UnimplementedCalcServer
	add grpctransport.Handler
}

func newCalcServer() pb.CalcServer {
	svr := &calcServer{}
	svr.add = luchen.NewGRPCTransportServer(
		calcEdp.makeAddEndpoint(),
		luchen.DecodePB[*pb.AddReq],
		luchen.EncodePB[*pb.AddResp],
	)
	return svr
}

func (s *calcServer) Add(ctx context.Context, req *pb.AddReq) (*pb.AddResp, error) {
	_, resp, err := s.add.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.AddResp), nil
}
