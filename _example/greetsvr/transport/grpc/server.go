package grpc

import (
	"context"
	"sync"

	"github.com/fengjx/go-halo/halo"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/fengjx/luchen"
	"github.com/fengjx/luchen/_example/greetsvr/connom/config"
	"github.com/fengjx/luchen/_example/greetsvr/pb"
)

var (
	server     *luchen.GRPCServer
	serverOnce sync.Once
)

func GetServer() *luchen.GRPCServer {
	serverOnce.Do(func() {
		serverConfig := config.GetConfig().Server.GRPC
		server = luchen.NewGRPCServer(
			serverConfig.ServerName,
			luchen.WithGRPCAddr(serverConfig.Listen),
		)
		server.RegisterServer(func(grpcServer *grpc.Server) {
			pb.RegisterGreeterServer(grpcServer, newGreeterServer())
		})
	})
	return server
}

func Start(_ context.Context) {
	go func() {
		logger := luchen.RootLogger().With(zap.Int64("goid", halo.GetGoID()))
		if err := GetServer().Start(); err != nil {
			logger.Panic("start grpc server err", zap.Error(err))
		}
	}()
}

func Stop(ctx context.Context) {
	logger := luchen.Logger(ctx)
	if err := GetServer().Stop(); err != nil {
		logger.Error("grpc server stop err", zap.Error(err))
	} else {
		logger.Info("grpc server was shutdown gracefully")
	}
}
