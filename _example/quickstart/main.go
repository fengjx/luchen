package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/fengjx/luchen/env"
	"github.com/fengjx/luchen/log"

	"github.com/fengjx/luchen/example/greetsvr/service"
	"github.com/fengjx/luchen/example/greetsvr/transport"
)

func init() {
	if env.IsLocal() {
		env.SetDefaultEtcdAddress([]string{"host.etcd.dev:2379"})
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	log.Info("app start")
	service.Init()
	transport.Start(ctx)

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGKILL)

	<-quit
	log.Info("app stop")
	cancel()
	transport.Stop(ctx)
}
