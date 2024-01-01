package luchen

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/fengjx/go-halo/addr"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-kit/kit/endpoint"
	"github.com/google/uuid"
	"go.uber.org/zap"

	httptransport "github.com/go-kit/kit/transport/http"
)

type HTTPServerOptions struct {
	addr     string
	metadata map[string]any
}

type HTTPServerOption func(*HTTPServerOptions)

func WithHTTPAddr(addr string) HTTPServerOption {
	return func(o *HTTPServerOptions) {
		o.addr = addr
	}
}

func WithHTTPMetadata(md map[string]any) HTTPServerOption {
	return func(o *HTTPServerOptions) {
		o.metadata = md
	}
}

type HTTPRouter = *chi.Mux

type HTTPServer struct {
	*baseServer
	httpServer *http.Server
	router     HTTPRouter
}

func NewHTTPServer(serviceName string, opts ...HTTPServerOption) *HTTPServer {
	options := &HTTPServerOptions{}
	for _, opt := range opts {
		opt(options)
	}
	if options.addr == "" {
		options.addr = defaultAddress
	}
	if options.metadata == nil {
		options.metadata = make(map[string]any)
	}
	router := chi.NewRouter()
	httpServer := &http.Server{
		Handler: router,
	}
	svr := &HTTPServer{
		baseServer: &baseServer{
			id:          uuid.NewString(),
			serviceName: serviceName,
			protocol:    ProtocolHTTP,
			address:     options.addr,
			metadata:    make(map[string]any),
		},
		httpServer: httpServer,
		router:     router,
	}
	svr.Use(
		middleware.Recoverer,
		middleware.RealIP,
		middleware.RequestID,
		TraceHTTPMiddleware,
	)
	return svr
}

func (s *HTTPServer) Start() error {
	s.Lock()
	ln, err := net.Listen("tcp", s.address)
	if err != nil {
		return err
	}
	address := ln.Addr().String()
	host, port, err := addr.ExtractHostPort(address)
	if err != nil {
		return err
	}
	s.address = fmt.Sprintf("%s:%s", host, port)
	s.metadata["ts"] = time.Now().UnixMilli()
	s.started = true
	RootLogger().Infof("http server[%s, %s] start", s.serviceName, s.id)
	s.Unlock()
	return s.httpServer.Serve(ln)
}

func (s *HTTPServer) Stop() error {
	s.RLock()
	if !s.started {
		s.RUnlock()
		return nil
	}
	s.RUnlock()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return s.httpServer.Shutdown(ctx)
}

// Use 中间件
func (s *HTTPServer) Use(middlewares ...HTTPMiddleware) *HTTPServer {
	for _, m := range middlewares {
		s.router.Use(m)
	}
	return s
}

// Handler 请求处理
func (s *HTTPServer) Handler(handlers ...HTTPHandler) *HTTPServer {
	for _, handler := range handlers {
		handler.Bind(s.router)
	}
	return s
}

type HTTPMiddleware func(http.Handler) http.Handler

// TraceHTTPMiddleware 链路跟踪
func TraceHTTPMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		traceID := r.Header.Get(TraceIDHeader)
		if traceID == "" {
			traceID = TraceID(r.Context())
		}
		if traceID == "" {
			traceID = uuid.NewString()
		}
		logger := Logger(r.Context())
		logger = logger.With(zap.String("traceId", traceID))
		ctx := WithLogger(r.Context(), logger)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

type HTTPHandler interface {
	Bind(router HTTPRouter)
}

// NewHTTPHandler 绑定 http 请求处理逻辑
func NewHTTPHandler(
	e endpoint.Endpoint,
	dec httptransport.DecodeRequestFunc,
	enc httptransport.EncodeResponseFunc,
	options ...httptransport.ServerOption,
) *httptransport.Server {
	return httptransport.NewServer(
		e,
		dec,
		enc,
		options...,
	)
}

func DecodeKvRequest[T any](ctx context.Context, r *http.Request) (interface{}, error) {
	logger := Logger(ctx)
	req := new(T)
	err := ShouldBind(r, req)
	if err != nil {
		logger.Error("decode request err", zap.Error(err))
		errn := &Errno{
			Code:     4,
			HttpCode: http.StatusBadRequest,
			Msg:      err.Error(),
		}
		return nil, errn
	}
	return req, nil
}

func DecodeJsonRequest[T any](ctx context.Context, r *http.Request) (interface{}, error) {
	logger := Logger(ctx)
	req := new(T)
	err := ShouldBindJSON(r, req)
	if err != nil {
		logger.Error("decode request err", zap.Error(err))
		errn := &Errno{
			Code:     4,
			HttpCode: http.StatusBadRequest,
			Msg:      err.Error(),
		}
		return nil, errn
	}
	return req, nil
}
