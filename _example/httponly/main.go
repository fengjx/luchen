package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/fengjx/luchen"
	"github.com/fengjx/luchen/log"

	"github.com/fengjx/luchen/example/httponly/logic"
	"github.com/fengjx/luchen/example/httponly/transport/http"
)

func main() {
	httpServer := http.GetServer()
	logic.Init(httpServer)
	luchen.Start(httpServer)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	<-quit
	log.Info("server shutdown")
	luchen.Stop()
}
