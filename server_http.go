package luchen

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"reflect"
	"time"

	"github.com/fengjx/go-halo/addr"
	"github.com/fengjx/xin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	httptransport "github.com/go-kit/kit/transport/http"

	"github.com/fengjx/go-halo/errs"
	"github.com/fengjx/go-halo/json"
	"github.com/fengjx/luchen/env"
	"github.com/fengjx/luchen/log"
	"github.com/fengjx/luchen/marshal"
	"github.com/fengjx/luchen/types"
)

const (
	HeaderRspMeta = "X-Rsp-Meta"
)

type (
	httpRequestKey struct{}

	// HTTPTransportServer go-kit http transport server
	HTTPTransportServer = httptransport.Server
)

// HTTPServer http server 实现
type HTTPServer struct {
	*baseServer
	xin *xin.Xin
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
		options.serviceName = fmt.Sprintf("%s-%s", env.GetAppName(), "http-server")
	}
	if options.metadata == nil {
		options.metadata = make(map[string]any)
	}
	x := xin.New()
	x.Use(
		RecoverHTTPMiddleware,
		TraceHTTPMiddleware,
	)
	svr := &HTTPServer{
		baseServer: &baseServer{
			id:          uuid.NewString(),
			serviceName: options.serviceName,
			protocol:    ProtocolHTTP,
			address:     options.addr,
			metadata:    make(map[string]any),
		},
		xin: x,
	}
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
	log.Infof("http server[%s, %s, %s] start", s.serviceName, s.address, s.id)
	s.Unlock()
	return s.xin.Serve(ln, true)
}

// Stop 停止服务
func (s *HTTPServer) Stop() error {
	s.RLock()
	if !s.started {
		s.RUnlock()
		return nil
	}
	s.RUnlock()
	return s.xin.Shutdown(30 * time.Second)
}

// Mux 获取路由复用器
func (s *HTTPServer) Mux() *xin.Mux {
	return s.xin.Mux()
}

// Handle 注册 http 路由
func (s *HTTPServer) Handle(def *EndpointDefine) {
	hts := NewHTTPTransportServer(def)
	s.Mux().Handle(def.Path, hts)
}

// TraceHTTPMiddleware 链路跟踪
func TraceHTTPMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r, traceID := TraceHTTPRequest(r)
		ctx := log.WithLogger(r.Context(), zap.String("traceId", traceID))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RecoverHTTPMiddleware 恢复 panic
func RecoverHTTPMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		errs.RecoverFunc(func(err any, stack *errs.Stack) {
			if err == http.ErrAbortHandler {
				// we don't recover http.ErrAbortHandler so the response
				// to the client is aborted, this should not be logged
				panic(err)
			}
			if r.Header.Get("Connection") != "Upgrade" {
				WriteError(r.Context(), w, ErrSystem.WithDetail(fmt.Sprintf("%v", stack)))
			}
		})
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

// NewHTTPTransportServer http handler 绑定 endpoint
func NewHTTPTransportServer(
	def *EndpointDefine,
	options ...httptransport.ServerOption,
) *HTTPTransportServer {
	e := def.Endpoint
	middlewares := GlobalHTTPMiddlewares
	if len(def.Middlewares) > 0 {
		middlewares = append(middlewares, def.Middlewares...)
	}
	if len(middlewares) > 0 {
		e = EndpointChain(e, middlewares...)
	}
	dec := getHTTPRequestDecoder(def.ReqType)
	options = append(options,
		httptransport.ServerBefore(contextServerBefore),
		httptransport.ServerErrorEncoder(errorEncoder),
	)
	return httptransport.NewServer(
		e,
		dec,
		encodeHTTPResponse,
		options...,
	)
}

func errorEncoder(ctx context.Context, err error, w http.ResponseWriter) {
	WriteError(ctx, w, err)
}

// WriteError write error
func WriteError(ctx context.Context, w http.ResponseWriter, err error) {
	errn, ok := FromError(err)
	if !ok {
		errn = ErrSystem
	}
	w.WriteHeader(errn.HttpCode)
	rspMeta := &types.RspMeta{
		Code:       int32(errn.Code),
		Msg:        errn.Msg,
		TraceId:    TraceID(ctx),
		ServerTime: time.Now().UnixMilli(),
	}
	if !env.IsProd() {
		rspMeta.Detail = errn.GetDetail()
	}
	rspMetaJson, _ := json.ToJson(rspMeta)
	w.Header().Set(HeaderRspMeta, rspMetaJson)
	_, _ = w.Write([]byte(""))
}

func contextServerBefore(ctx context.Context, req *http.Request) context.Context {
	startTime := time.Now()
	contentType := req.Header.Get("Content-Type")
	marshaller := marshal.GetMarshallerByContentType(contentType)

	rip := xin.GetRealIP(req)
	if rip == "" {
		rip = req.RemoteAddr
	}
	h := &Header{
		Method:      req.Method,
		Endpoint:    req.RequestURI,
		Protocol:    string(ProtocolHTTP),
		Host:        req.Host,
		CLientIP:    rip,
		TraceID:     TraceID(ctx),
		StartTime:   startTime,
		ContentType: contentType,
	}
	ctx = withHeader(ctx, h)
	ctx = context.WithValue(ctx, httpRequestKey{}, req)
	ctx = withMarshaller(ctx, marshaller)
	return ctx
}

// getHTTPRequestDecoder 解码 http pb 请求
func getHTTPRequestDecoder(typ reflect.Type) httptransport.DecodeRequestFunc {
	return func(ctx context.Context, req *http.Request) (any, error) {
		marshaller := Marshaller(ctx)
		payload, err := io.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
		if len(payload) == 0 {
			return nil, nil
		}
		data := reflect.New(typ).Interface()
		err = marshaller.Unmarshal(payload, data)
		if err != nil {
			log.ErrorCtx(ctx, "unmarshal request error", zap.Error(err))
			return nil, err
		}
		return reflect.ValueOf(data).Elem().Interface(), nil
	}
}

// encodeHTTPResponse 编码 http 响应
func encodeHTTPResponse(ctx context.Context, w http.ResponseWriter, data any) error {
	marshaller := Marshaller(ctx)
	bytes, err := marshaller.Marshal(data)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", marshaller.ContentType())
	rspMeta := &types.RspMeta{
		Code:       0,
		TraceId:    TraceID(ctx),
		ServerTime: time.Now().UnixMilli(),
	}
	rspMetaJson, _ := json.ToJson(rspMeta)
	w.Header().Set(HeaderRspMeta, rspMetaJson)
	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
	return nil
}

// GlobalHTTPMiddlewares 全局 HTTP 中间件
var GlobalHTTPMiddlewares []Middleware

// UseGlobalHTTPMiddleware 注册全局 HTTP 中间件
// 中间件的执行顺序与注册顺序相同，先注册的中间件先执行
func UseGlobalHTTPMiddleware(m ...Middleware) {
	GlobalHTTPMiddlewares = append(GlobalHTTPMiddlewares, m...)
}
