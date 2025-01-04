package luchen

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/fengjx/go-halo/addr"
	"github.com/fengjx/xin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"

	httptransport "github.com/go-kit/kit/transport/http"

	"github.com/fengjx/luchen/env"
	"github.com/fengjx/luchen/log"
	"github.com/fengjx/luchen/marshal"
)

type (
	httpRequestHeaderKey struct{}
	httpRequestURLKey    struct{}

	// HTTPTransportServer go-kit http transport server
	HTTPTransportServer = httptransport.Server
)

var (
	// HTTPRequestHeaderCtxKey context http header
	HTTPRequestHeaderCtxKey = httpRequestHeaderKey{}
	// HTTPRequestURLCtxKey context http url
	HTTPRequestURLCtxKey = httpRequestURLKey{}

	// NopHTTPRequestDecoder 不需要解析请求参数
	NopHTTPRequestDecoder = httptransport.NopRequestDecoder
)

// HTTPServer http server 实现
type HTTPServer struct {
	*baseServer
	*xin.Xin
	started bool
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
	svr := &HTTPServer{
		baseServer: &baseServer{
			id:          uuid.NewString(),
			serviceName: options.serviceName,
			protocol:    ProtocolHTTP,
			address:     options.addr,
			metadata:    make(map[string]any),
		},
		Xin: x,
	}
	svr.Use(
		ClientIPHTTPMiddleware,
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
	log.Infof("http server[%s, %s, %s] start", s.serviceName, s.address, s.id)
	s.Unlock()
	return s.Serve(ln, true)
}

// Stop 停止服务
func (s *HTTPServer) Stop() error {
	s.RLock()
	if !s.started {
		s.RUnlock()
		return nil
	}
	s.RUnlock()
	return s.Shutdown(30 * time.Second)
}

// TraceHTTPMiddleware 链路跟踪
func TraceHTTPMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r, traceID := TraceHTTPRequest(r)
		ctx := log.WithLogger(r.Context(), zap.String("traceId", traceID))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// ClientIPHTTPMiddleware 获取客户端IP
func ClientIPHTTPMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rip := xin.GetRealIP(r)
		if rip == "" {
			rip = r.RemoteAddr
		}
		ctx := withRequestClientIP(r.Context(), rip)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// NewHTTPTransportServer http handler 绑定 endpoint
func NewHTTPTransportServer(
	e Endpoint,
	dec httptransport.DecodeRequestFunc,
	enc httptransport.EncodeResponseFunc,
	options ...httptransport.ServerOption,
) *HTTPTransportServer {
	options = append(options, httptransport.ServerBefore(contextServerBefore))
	return httptransport.NewServer(
		e,
		dec,
		enc,
		options...,
	)
}

func contextServerBefore(ctx context.Context, req *http.Request) context.Context {
	startTime := time.Now()
	contentType := req.Header.Get("Content-Type")
	marshaller := marshal.GetMarshallerByContentType(contentType)

	ctx = context.WithValue(ctx, HTTPRequestHeaderCtxKey, req.Header)
	ctx = context.WithValue(ctx, HTTPRequestURLCtxKey, req.URL)
	ctx = withRequestStartTime(ctx, startTime)
	ctx = withRequestEndpoint(ctx, req.RequestURI)
	ctx = withRequestProtocol(ctx, req.Proto)
	ctx = withMethod(ctx, req.Method)
	ctx = withMarshaller(ctx, marshaller)
	return ctx
}

// DecodeHTTPPbRequest 解码 http pb 请求
func DecodeHTTPPbRequest[p proto.Message](ctx context.Context, req *http.Request) (request any, err error) {
	marshaller := Marshaller(ctx)
	payload, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	err = marshaller.Unmarshal(payload, &request)
	if err != nil {
		return nil, err
	}
	return request, nil
}

// EncodeHTTPPbResponse 编码 http pb 响应
func EncodeHTTPPbResponse[p proto.Message](ctx context.Context, w http.ResponseWriter, data any) error {
	marshaller := Marshaller(ctx)
	bytes, err := marshaller.Marshal(data)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", marshaller.ContentType())
	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
	return nil
}
