package luchen

import (
	"context"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type GatewayOptions struct {
}

type GatewayOption func(*GatewayOptions)

func With() GatewayOption {
	return func(o *GatewayOptions) {
	}
}

// Gateway 网关服务
type Gateway struct {
	config GatewayConfig
	server *http.Server
}

// NewGGateway 创建 gateway 服务
func NewGGateway(cfg GatewayConfig, opts ...GatewayOption) *Gateway {
	options := &GatewayOptions{}
	_ = options
	g := &Gateway{
		config: cfg,
	}
	return g
}

func (g *Gateway) Start(ctx context.Context) error {
	proxy := &httputil.ReverseProxy{}
	router := chi.NewRouter()
	router.Handle("/", proxy)
	server := &http.Server{
		Addr:    g.config.Listen,
		Handler: router,
	}
	g.server = server
	return g.server.ListenAndServe()
}

func (g *Gateway) Stop() {
	RootLogger().Warn("gateway server stop")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := g.server.Shutdown(ctx); err != nil {
		RootLogger().Error("error while shutting down gateway", zap.Error(err))
	} else {
		RootLogger().Info("gateway was shutdown gracefully")
	}
}
