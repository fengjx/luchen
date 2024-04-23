package luchen

import (
	"context"
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/fengjx/go-halo/hook"
	"go.uber.org/zap"

	"github.com/fengjx/luchen/log"
)

// Protocol 服务协议
type Protocol string

const (
	defaultAddress = ":0"

	// ProtocolHTTP http 协议
	ProtocolHTTP Protocol = "http"
	// ProtocolGRPC grpc 协议
	ProtocolGRPC Protocol = "grpc"

	beforeStopHookEvent = "before-stop-hook"
)

var (
	servers []Server
)

// ServerOptions server 选项
type ServerOptions struct {
	serviceName string
	addr        string
	metadata    map[string]any
}

// ServerOption grpc server 选项赋值
type ServerOption func(*ServerOptions)

// WithServerAddr server 监听地址
func WithServerAddr(addr string) ServerOption {
	return func(o *ServerOptions) {
		o.addr = addr
	}
}

// WithServerMetadata server 注册信息 metadata，单体服务无需关注
func WithServerMetadata(md map[string]any) ServerOption {
	return func(o *ServerOptions) {
		o.metadata = md
	}
}

// WithServiceName server 名称，在微服务中作为一组服务名称标识，单体服务则无需关注
func WithServiceName(serviceName string) ServerOption {
	return func(o *ServerOptions) {
		o.serviceName = serviceName
	}
}

// ServiceInfo 服务节点信息
type ServiceInfo struct {
	ID       string         `json:"id"`
	Name     string         `json:"name"`
	Protocol Protocol       `json:"protocol"`
	Addr     string         `json:"addr"`
	Metadata map[string]any `json:"metadata,omitempty"`
}

// Server server 接口定义
type Server interface {
	Start() error
	Stop() error
	GetServiceInfo() *ServiceInfo
}

type baseServer struct {
	sync.RWMutex
	id          string
	serviceName string
	protocol    Protocol
	address     string
	metadata    map[string]any

	started bool
}

func (s *baseServer) GetServiceInfo() *ServiceInfo {
	s.RLock()
	if s.started {
		s.RUnlock()
		return &ServiceInfo{
			Protocol: s.protocol,
			ID:       s.id,
			Name:     s.serviceName,
			Addr:     s.address,
			Metadata: s.metadata,
		}
	}
	s.RUnlock()
	for {
		select {
		case <-time.After(time.Millisecond):
			return s.GetServiceInfo()
		}
	}
}

// Start 启动服务
func Start(svrs ...Server) {
	start(nil, svrs...)
}

// StartWithContext 启动服务
func StartWithContext(ctx context.Context, svrs ...Server) {
	start(ctx, svrs...)
}

func start(ctx context.Context, svrs ...Server) {
	servers = svrs
	for _, server := range servers {
		svr := server
		go func() {
			if err := svr.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				log.Panic(
					"server start err",
					zap.String("server_name", svr.GetServiceInfo().Name),
					zap.Error(err),
				)
			}
		}()
	}
	if ctx == nil {
		return
	}
	select {
	case <-ctx.Done():
		Stop()
	}
}

// Stop 停止服务
func Stop() {
	DoStopHook()
	for _, server := range servers {
		// 停止服务
		if err := server.Stop(); err != nil {
			log.Error("server stop err", zap.Error(err))
		}
		log.Info(
			"server stop gracefully",
			zap.String("name", server.GetServiceInfo().Name),
		)
	}
}

// AddBeforeStopHook 注册服务停止前回调函数
func AddBeforeStopHook(handler func()) {
	hook.AddHook(beforeStopHookEvent, 100, handler)
}

// DoStopHook 执行服务停止前的回调函数
func DoStopHook() {
	hook.DoHooks(beforeStopHookEvent)
}

// DataWrapper 对数据重新组装
type DataWrapper func(data any) any
