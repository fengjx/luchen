package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/fengjx/luchen"
)

func init() {
	if luchen.IsLocal() {
		luchen.SetDefaultEtcdAddress([]string{"host.etcd.dev:2379"})
	}
}

func main() {
	config := luchen.MustLoadConfig[luchen.GatewayConfig]("_example/gateway/gateway.yaml")
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
