package luchen

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/fengjx/go-halo/addr"
	"github.com/go-kit/kit/endpoint"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	grpctransport "github.com/go-kit/kit/transport/grpc"
)

// GRPCServerOptions grpc server 选项
type GRPCServerOptions struct {
	addr     string
	metadata map[string]any
}

// GRPCServerOption grpc server 选项赋值
type GRPCServerOption func(*GRPCServerOptions)

// WithGRPCAddr grpc server 监听地址
func WithGRPCAddr(addr string) GRPCServerOption {
	return func(o *GRPCServerOptions) {
		o.addr = addr
	}
}

// WithGRPCMetadata grpc server 注册信息 metadata
func WithGRPCMetadata(md map[string]any) GRPCServerOption {
	return func(o *GRPCServerOptions) {
		o.metadata = md
	}
}

// GRPCServer grpc server 实现
type GRPCServer struct {
	*baseServer
	server *grpc.Server
}

// NewGRPCServer 创建 grpc server
func NewGRPCServer(serviceName string, opts ...GRPCServerOption) *GRPCServer {
	options := &GRPCServerOptions{}
	for _, opt := range opts {
		opt(options)
	}
	if options.addr == "" {
		options.addr = defaultAddress
	}
	if options.metadata == nil {
		options.metadata = make(map[string]any)
	}
	server := grpc.NewServer(grpc.UnaryInterceptor(kitgrpc.Interceptor))
	return &GRPCServer{
		baseServer: &baseServer{
			id:          uuid.NewString(),
			serviceName: serviceName,
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
	RootLogger().Infof("grpc server[%s, %s, %s] start", s.serviceName, s.address, s.id)
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

// RegisterServer 注册 grpc 接口实现
func (s *GRPCServer) RegisterServer(rh RegisterHandler) *GRPCServer {
	rh(s.server)
	return s
}

// NewGRPCHandler 绑定 grpc 请求处理逻辑
func NewGRPCHandler(
	e endpoint.Endpoint,
	dec grpctransport.DecodeRequestFunc,
	enc grpctransport.EncodeResponseFunc,
	options ...grpctransport.ServerOption,
) *grpctransport.Server {
	opts := []grpctransport.ServerOption{
		grpctransport.ServerBefore(func(ctx context.Context, md metadata.MD) context.Context {
			ctx, traceID := TraceGRPC(ctx, md)
			logger := Logger(ctx)
			logger = logger.With(zap.String("traceId", traceID))
			ctx = WithLogger(ctx, logger)
			ctx = metadata.NewOutgoingContext(ctx, md)
			return NewContext(ctx)
		}),
		grpctransport.ServerErrorHandler(NewLogGRPCErrorHandler()),
	}
	opts = append(opts, options...)
	return grpctransport.NewServer(
		e,
		dec,
		enc,
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
	logger := Logger(ctx)
	logger.Error("handle grpc err", zap.Error(err))
}

// DecodePB protobuf 解码
func DecodePB[T any](_ context.Context, req interface{}) (interface{}, error) {
	pbReq := req.(*T)
	return pbReq, nil
}

// EncodePB protobuf 编码
func EncodePB[T any](_ context.Context, resp interface{}) (interface{}, error) {
	pbResp := resp.(*T)
	return pbResp, nil
}
