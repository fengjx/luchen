package luchen

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/fengjx/go-halo/addr"
	"github.com/fengjx/go-halo/json"
	"github.com/go-kit/kit/endpoint"
	"github.com/google/uuid"
	"go.uber.org/zap"

	httptransport "github.com/go-kit/kit/transport/http"
)

type (
	httpRequestHeaderKey struct{}
	httpRequestURLKey    struct{}
)

var (
	// HTTPRequestHeaderCtxKey context http header
	HTTPRequestHeaderCtxKey = httpRequestHeaderKey{}
	// HTTPRequestURLCtxKey context http url
	HTTPRequestURLCtxKey = httpRequestURLKey{}
)

// HTTPServer http server 实现
type HTTPServer struct {
	*baseServer
	httpServer *http.Server
	router     *ServeMux
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
	mux := NewServeMux()
	httpServer := &http.Server{
		Handler: mux,
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
		router:     mux,
	}
	svr.Use(TraceHTTPMiddleware)
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

// Static 静态文件路径
func (s *HTTPServer) Static(prefix string, dir string) *HTTPServer {
	fs := http.FileServer(http.Dir(dir))
	s.router.Handle(prefix, http.StripPrefix(prefix, fs))
	return s
}

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
	Bind(router *ServeMux)
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
	ctx = context.WithValue(ctx, HTTPRequestHeaderCtxKey, req.Header)
	ctx = context.WithValue(ctx, HTTPRequestURLCtxKey, req.URL)
	ctx = withRequestEndpoint(ctx, req.RequestURI)
	ctx = withRequestProtocol(ctx, req.Proto)
	ctx = withMethod(ctx, req.Method)
	return ctx
}

// DecodeHTTPParamRequest 解析 http request query 和 form 参数
func DecodeHTTPParamRequest[T any](ctx context.Context, r *http.Request) (interface{}, error) {
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

// DecodeHTTPJSONRequest 解析 http request body json 参数
func DecodeHTTPJSONRequest[T any](ctx context.Context, r *http.Request) (interface{}, error) {
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

// EncodeHTTPJSONResponse http 返回json数据
// wrapper 对数据重新包装
func EncodeHTTPJSONResponse(wrapper DataWrapper) httptransport.EncodeResponseFunc {
	return func(ctx context.Context, w http.ResponseWriter, response interface{}) error {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		if headerer, ok := response.(httptransport.Headerer); ok {
			for k, values := range headerer.Headers() {
				for _, v := range values {
					w.Header().Add(k, v)
				}
			}
		}
		code := http.StatusOK
		if sc, ok := response.(httptransport.StatusCoder); ok {
			code = sc.StatusCode()
		}
		w.WriteHeader(code)
		if code == http.StatusNoContent {
			return nil
		}
		traceID := TraceID(ctx)
		if traceID != "" {
			w.Header().Set(TraceIDHeader, traceID)
		}
		return json.NewEncoder(w).Encode(wrapper(response))
	}
}
