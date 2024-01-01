package luchen

import (
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/fengjx/go-halo/hook"
	"go.uber.org/zap"
)

type Protocol string

const (
	defaultAddress = ":0"

	ProtocolHTTP Protocol = "http"
	ProtocolGRPC Protocol = "grpc"

	beforeStopHookEvent = "before-stop-hook"
)

var (
	servers []Server
)

type ServiceInfo struct {
	Id       string         `json:"id"`
	Name     string         `json:"name"`
	Protocol Protocol       `json:"protocol"`
	Addr     string         `json:"addr"`
	Metadata map[string]any `json:"metadata,omitempty"`
}

type ServerHookHandler func(info ServiceInfo)

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
			Id:       s.id,
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
	servers = svrs
	for _, server := range servers {
		svr := server
		go func() {
			if err := svr.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				RootLogger().Panic("server start err", zap.Error(err))
			}
		}()
	}
}

// Stop 停止服务
func Stop() {
	DoStopHook()
	for _, server := range servers {
		if err := server.Stop(); err != nil {
			RootLogger().Error("server stop err", zap.Error(err))
		}
	}
}

// AddBeforeStopHook 注册服务停止前回调函数
func AddBeforeStopHook(handler func()) {
	hook.AddCustomStartHook(beforeStopHookEvent, handler, 100)
}

// DoStopHook 执行服务停止前的回调函数
func DoStopHook() {
	hook.DoCustomHooks(beforeStopHookEvent)
}
