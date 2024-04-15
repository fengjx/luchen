package grpc

import (
	"sync"

	"google.golang.org/grpc"

	"github.com/fengjx/luchen"

	"github.com/fengjx/luchen/example/greetsvr/connom/config"
	"github.com/fengjx/luchen/example/greetsvr/pb"
)

var (
	server     *luchen.GRPCServer
	serverOnce sync.Once
)

func GetServer() *luchen.GRPCServer {
	serverOnce.Do(func() {
		serverConfig := config.GetConfig().Server.GRPC
		server = luchen.NewGRPCServer(
			luchen.WithServiceName(serverConfig.ServerName),
			luchen.WithServerAddr(serverConfig.Listen),
		)
		server.RegisterService(func(grpcServer *grpc.Server) {
			pb.RegisterGreeterServer(grpcServer, newGreeterServer())
		})
	})
	return server
}
