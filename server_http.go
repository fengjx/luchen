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
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"

	httptransport "github.com/go-kit/kit/transport/http"
)

var (
	// HTTPRequestHeaderContextKey context http header
	HTTPRequestHeaderContextKey ctxType = "ctx.http.request.header"
	// HTTPRequestURLContextKey context http url
	HTTPRequestURLContextKey ctxType = "ctx.http.request.url"
)

// HTTPRouter http 请求路由注册
type HTTPRouter = *chi.Mux

// HTTPServer http server 实现
type HTTPServer struct {
	*baseServer
	httpServer *http.Server
	router     HTTPRouter
}

// NewHTTPServer 创建 http server
// opts 查看 ServerOptions
func NewHTTPServer(opts ...ServerOption) *HTTPServer {
	options := &ServerOptions{}
	for _, opt := range opts {
		opt(options)
	}
	if options.addr == "" {
		options.addr = defaultAddress
	}
	if options.serviceName == "" {
		options.serviceName = fmt.Sprintf("%s-%s", GetAppName(), "http-server")
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
			serviceName: options.serviceName,
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

// Start 启动服务
func (s *HTTPServer) Start() error {
	s.Lock()
	ln, err := net.Listen("tcp", s.address)
	if err != nil {
		s.Unlock()
		return err
	}
	address := ln.Addr().String()
	host, port, err := addr.ExtractHostPort(address)
	if err != nil {
		s.Unlock()
		return err
	}
	s.address = fmt.Sprintf("%s:%s", host, port)
	s.metadata["ts"] = time.Now().UnixMilli()
	s.started = true
	RootLogger().Infof("http server[%s, %s, %s] start", s.serviceName, s.address, s.id)
	s.Unlock()
	return s.httpServer.Serve(ln)
}

// Stop 停止服务
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

// HTTPMiddleware http 请求中间件
type HTTPMiddleware func(http.Handler) http.Handler

// TraceHTTPMiddleware 链路跟踪
func TraceHTTPMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r, traceID := TraceHTTPRequest(r)
		logger := Logger(r.Context())
		logger = logger.With(zap.String("traceId", traceID))
		ctx := WithLogger(r.Context(), logger)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// HTTPHandler http 请求处理器接口
type HTTPHandler interface {
	// Bind 绑定路由
	Bind(router HTTPRouter)
}

// NewHTTPHandler 绑定 http 请求处理逻辑
func NewHTTPHandler(
	e endpoint.Endpoint,
	dec httptransport.DecodeRequestFunc,
	enc httptransport.EncodeResponseFunc,
	options ...httptransport.ServerOption,
) *httptransport.Server {
	options = append(options, httptransport.ServerBefore(contextServerBefore))
	return httptransport.NewServer(
		e,
		dec,
		enc,
		options...,
	)
}

func contextServerBefore(ctx context.Context, req *http.Request) context.Context {
	ctx = context.WithValue(ctx, HTTPRequestHeaderContextKey, req.Header)
	ctx = context.WithValue(ctx, HTTPRequestURLContextKey, req.URL)
	return ctx
}

// DecodeParamHTTPRequest 解析 http request query 和 form 参数
func DecodeParamHTTPRequest[T any](ctx context.Context, r *http.Request) (interface{}, error) {
	logger := Logger(ctx)
	req := new(T)
	err := ShouldBind(r, req)
	if err != nil {
		logger.Error("decode request err", zap.Error(err))
		errn := &Errno{
			Code:     4,
			HTTPCode: http.StatusBadRequest,
			Msg:      err.Error(),
		}
		return nil, errn
	}
	return req, nil
}

// DecodeJSONRequest 解析 http request body json 参数
func DecodeJSONRequest[T any](ctx context.Context, r *http.Request) (interface{}, error) {
	logger := Logger(ctx)
	req := new(T)
	err := ShouldBindJSON(r, req)
	if err != nil {
		logger.Error("decode request err", zap.Error(err))
		errn := &Errno{
			Code:     4,
			HTTPCode: http.StatusBadRequest,
			Msg:      err.Error(),
		}
		return nil, errn
	}
	return req, nil
}

// CreateHTTPJSONEncoder http 返回json数据
// wrapper 对数据重新包装
func CreateHTTPJSONEncoder(wrapper DataWrapper) httptransport.EncodeResponseFunc {
	return func(ctx context.Context, w http.ResponseWriter, response interface{}) error {
		traceID := TraceID(ctx)
		if traceID != "" {
			w.Header().Set(TraceIDHeader, traceID)
		}
		w.Header().Set("Content-Type", "application/json")
		return jsoniter.NewEncoder(w).Encode(wrapper(response))
	}
}
