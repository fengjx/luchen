package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/fengjx/luchen"
	"github.com/fengjx/luchen/env"
	"github.com/fengjx/luchen/log"

	"github.com/fengjx/luchen/example/quickstart/connom/config"
	"github.com/fengjx/luchen/example/quickstart/logic"
)

func init() {
	if env.IsLocal() {
		env.SetDefaultEtcdAddress([]string{"host.etcd.dev:2379"})
	}
}

func main() {
	log.Info("server start")
	serverConfig := config.GetConfig().Server
	hs := luchen.NewHTTPServer(
		luchen.WithServiceName(serverConfig.HTTP.ServerName),
		luchen.WithServerAddr(serverConfig.HTTP.Listen), // 不指定则使用随机端口
	)
	gs := luchen.NewGRPCServer(
		luchen.WithServiceName(serverConfig.GRPC.ServerName),
		luchen.WithServerAddr(serverConfig.GRPC.Listen), // 不指定则使用随机端口
	)
	registrar := luchen.NewEtcdV3Registrar(
		hs,
		gs,
	)
	logic.Init(hs, gs)
	// 注册并启动服务
	registrar.Register()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGKILL)

	<-quit
	log.Info("server stop")
	// 摘除并停止服务
	registrar.Deregister()
}
