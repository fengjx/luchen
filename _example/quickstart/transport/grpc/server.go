package grpc

import (
	"sync"

	"github.com/fengjx/luchen"

	"github.com/fengjx/luchen/example/quickstart/connom/config"
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
			luchen.WithServerAddr(serverConfig.Listen), // 不指定则使用随机端口
		)
	})
	return server
}
