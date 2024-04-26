package endpoint

import (
	"github.com/fengjx/luchen"
	"google.golang.org/grpc"

	"github.com/fengjx/luchen/example/quickstart/pb"
)

func Init(hs *luchen.HTTPServer, gs *luchen.GRPCServer) {
	// 注册 http 路由
	hs.Handler(
		&greeterHandler{},
	)
	// 注册 grpc 服务
	gs.RegisterService(func(s *grpc.Server) {
		pb.RegisterGreeterServer(s, newGreeterServer())
	})
}
