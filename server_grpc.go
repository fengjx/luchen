package luchen

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/fengjx/go-halo/addr"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	grpctransport "github.com/go-kit/kit/transport/grpc"

	"github.com/fengjx/luchen/env"
	"github.com/fengjx/luchen/log"
)

type (
	// GRPCTransportServer grpc transport server
	GRPCTransportServer = grpctransport.Server
)

// GRPCServer grpc server 实现
type GRPCServer struct {
	*baseServer
	server *grpc.Server
}

// NewGRPCServer 创建 grpc server
// opts 查看 ServerOptions
func NewGRPCServer(opts ...ServerOption) *GRPCServer {
	options := &ServerOptions{}
	for _, opt := range opts {
		opt(options)
	}
	if options.addr == "" {
		options.addr = defaultAddress
	}
	if options.serviceName == "" {
		options.serviceName = fmt.Sprintf("%s-%s", env.GetAppName(), "grpc-server")
	}
	if options.metadata == nil {
		options.metadata = make(map[string]any)
	}
	server := grpc.NewServer(grpc.UnaryInterceptor(grpctransport.Interceptor))
	return &GRPCServer{
		baseServer: &baseServer{
			id:          uuid.NewString(),
			serviceName: options.serviceName,
			protocol:    ProtocolGRPC,
			address:     options.addr,
			metadata:    make(map[string]any),
		},
		server: server,
	}
}

// Start 停止服务
func (s *GRPCServer) Start() error {
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
	log.Infof("grpc server[%s, %s, %s] start", s.serviceName, s.address, s.id)
	s.Unlock()
	return s.server.Serve(ln)
}

// Stop 停止服务
func (s *GRPCServer) Stop() error {
	s.RLock()
	if !s.started {
		s.RUnlock()
		return nil
	}
	s.RUnlock()
	s.server.GracefulStop()
	return nil
}

// RegisterHandler 注册 grpc handler
type RegisterHandler func(grpcServer *grpc.Server)

// RegisterService 注册 grpc 接口实现
func (s *GRPCServer) RegisterService(desc *grpc.ServiceDesc, impl any) {
	s.server.RegisterService(desc, impl)
}

// NewGRPCTransportServer grpc handler 绑定 endpoint
func NewGRPCTransportServer(
	def *EdnpointDefine,
	options ...grpctransport.ServerOption,
) *GRPCTransportServer {
	e := EndpointChain(def.Endpoint, def.Middlewares...)
	opts := []grpctransport.ServerOption{
		grpctransport.ServerBefore(func(ctx context.Context, md metadata.MD) context.Context {
			ctx, traceID := TraceGRPC(ctx, md)
			ctx = log.WithLogger(ctx, zap.String("traceId", traceID))
			ctx = metadata.NewOutgoingContext(ctx, md)
			return ctx
		}),
		grpctransport.ServerErrorHandler(NewLogGRPCErrorHandler()),
	}
	opts = append(opts, options...)
	return grpctransport.NewServer(
		e,
		decodePB,
		encodePB,
		opts...,
	)
}

// LogGRPCErrorHandler grpc 接口错误处理器
type LogGRPCErrorHandler struct {
}

// NewLogGRPCErrorHandler 创建 LogGRPCErrorHandler
func NewLogGRPCErrorHandler() *LogGRPCErrorHandler {
	return &LogGRPCErrorHandler{}
}

// Handle 统一错误处理
func (h *LogGRPCErrorHandler) Handle(ctx context.Context, err error) {
	log.ErrorCtx(ctx, "handle grpc err", zap.Error(err), zap.Stack("stack"))
}

func decodePB(_ context.Context, req interface{}) (interface{}, error) {
	return req, nil
}

func encodePB(_ context.Context, resp interface{}) (interface{}, error) {
	return resp, nil
}
