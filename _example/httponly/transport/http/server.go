package http

import (
	"net/http"
	"sync"

	"github.com/fengjx/luchen"
	"github.com/fengjx/luchen/http/middleware"
	"github.com/fengjx/luchen/http/pprof"

	"github.com/fengjx/luchen/example/httponly/connom/config"
)

var (
	server     *luchen.HTTPServer
	serverOnce sync.Once
)

func GetServer() *luchen.HTTPServer {
	serverOnce.Do(func() {
		serverConfig := config.GetConfig().Server.HTTP
		server = luchen.NewHTTPServer(
			luchen.WithServiceName(serverConfig.ServerName),
			luchen.WithServerAddr(serverConfig.Listen),
		).Use(
			middleware.Recoverer,
			middleware.RequestID,
			middleware.RealIP,
			middleware.CorsHandler(middleware.CorsOptions{
				AllowedOrigins: []string{
					"http://localhost:3000",
					"http://localhost:8000",
					"http://admin.luchen.fun",
				},
				AllowedMethods: []string{
					http.MethodHead,
					http.MethodGet,
					http.MethodPost,
					http.MethodPut,
					http.MethodPatch,
					http.MethodDelete,
				},
				AllowedHeaders: []string{"*"},
				ExposedHeaders: []string{
					"Content-Disposition",
					"Content-Type",
				},
				AllowCredentials: true,
			}),
		).Handler(
			pprof.NewHandler().BasicAuth(map[string]string{
				"admin": "luchen",
			}),
		)
	})
	return server
}
