package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/fengjx/go-halo/fs"
	"github.com/fengjx/luchen"
	"github.com/fengjx/luchen/env"
	"github.com/fengjx/luchen/log"
	"go.uber.org/zap"
)

func init() {
	if env.IsLocal() {
		env.SetDefaultEtcdAddress([]string{"host.etcd.dev:2379"})
	}
}

func main() {
	configFile, err := fs.Lookup("gateway.yaml", 3)
	if err != nil {
		log.Panic("config file not found", zap.Error(err))
	}
	config := luchen.MustLoadConfig[luchen.GatewayConfig](configFile)
	gateway := luchen.NewGateway(
		config,
		luchen.WithGatewayPlugin(
			&luchen.TraceGatewayPlugin{},
		),
	)
	luchen.Start(gateway)

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGKILL)

	<-quit
	luchen.Stop()
}
