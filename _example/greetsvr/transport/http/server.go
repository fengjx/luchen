package http

import (
	"context"
	"errors"
	"net/http"
	"sync"

	"github.com/fengjx/go-halo/halo"
	"go.uber.org/zap"

	"github.com/fengjx/luchen"
	"github.com/fengjx/luchen/_example/greetsvr/connom/config"
)

var (
	server     *luchen.HTTPServer
	serverOnce sync.Once
)

func GetServer() *luchen.HTTPServer {
	serverOnce.Do(func() {
		serverConfig := config.GetConfig().Server.HTTP
		server = luchen.NewHTTPServer(
			serverConfig.ServerName,
			luchen.WithHTTPAddr(serverConfig.Listen),
		).Handler(
			newGreeterHandler(),
		)
	})
	return server
}

func Start(_ context.Context) {
	go func() {
		logger := luchen.RootLogger().With(zap.Int64("goid", halo.GetGoID()))
		if err := GetServer().Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Panic("start http server err", zap.Error(err))
		}
	}()
}

func Stop(ctx context.Context) {
	logger := luchen.Logger(ctx)
	if err := GetServer().Stop(); err != nil {
		logger.Error("http server stop err", zap.Error(err))
	} else {
		logger.Info("http server was shutdown gracefully")
	}
}
