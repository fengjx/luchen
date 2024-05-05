package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/fengjx/luchen"
	"github.com/fengjx/luchen/env"
	"github.com/fengjx/luchen/log"

	"github.com/fengjx/luchen/example/quickstart/logic"
	"github.com/fengjx/luchen/example/quickstart/transport/grpc"
	"github.com/fengjx/luchen/example/quickstart/transport/http"
)

func init() {
	// 这里是为了演示，实际开发建议通过环境变量设置
	// 查看文档：https://luchen.fun/guide/env
	if env.IsLocal() {
		env.SetDefaultEtcdAddress([]string{"host.etcd.dev:2379"})
	}
}

func main() {
	hs := http.GetServer()
	gs := grpc.GetServer()
	registrar := luchen.NewEtcdV3Registrar(
		hs,
		gs,
	)
	logic.Init(hs, gs)
	// 注册并启动服务
	registrar.Register()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	<-quit
	log.Info("server shutdown")
	// 摘除并停止服务
	registrar.Deregister()
}
