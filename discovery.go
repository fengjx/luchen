package luchen

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"sync"
	"sync/atomic"
	"time"

	"github.com/fengjx/go-halo/json"
	"github.com/fengjx/go-halo/logger"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/etcdv3"
	"go.uber.org/zap"
)

const (
	servicePrefix = "/luchen/services"
)

var (
	// ErrNoServer 没有可用的服务节点
	ErrNoServer = errors.New("no server available")

	etcdV3SelectorCache     = make(map[string]*EtcdV3Selector)
	etcdV3SelectorCacheLock = newSegmentLock(10)
)

// EtcdV3Registrar 基于 etcd 实现的服务注册器
type EtcdV3Registrar struct {
	delegate map[string]*etcdv3.Registrar
	servers  []Server
}

// NewEtcdV3Registrar 创建一个服务注册器
// servers 需要注册的服务
func NewEtcdV3Registrar(servers ...Server) *EtcdV3Registrar {
	return &EtcdV3Registrar{
		servers:  servers,
		delegate: make(map[string]*etcdv3.Registrar),
	}
}

// Register 启动并注册服务到 etcd
func (r *EtcdV3Registrar) Register() {
	for _, server := range r.servers {
		svr := server
		go func() {
			if err := svr.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				RootLogger().Panic(
					"server start err",
					zap.String("name", svr.GetServiceInfo().Name),
					zap.Error(err),
				)
			}
		}()
		r.register(svr.GetServiceInfo())
	}
}

func (r *EtcdV3Registrar) register(serviceInfo *ServiceInfo) {
	key := path.Join(servicePrefix, serviceInfo.Name, serviceInfo.ID)
	params := url.Values{}
	info, _ := json.ToJson(serviceInfo)
	params.Set("info", info)
	value := fmt.Sprintf("%s://%s?%s", serviceInfo.Protocol, serviceInfo.Addr, params.Encode())
	registar := etcdv3.NewRegistrar(NewDefaultEtcdV3Client(), etcdv3.Service{
		Key:   key,
		Value: value,
	}, NewKitLogger(fmt.Sprintf("%s-%s", "register", serviceInfo.Name), logger.InfoLevel))
	r.delegate[serviceInfo.ID] = registar
	registar.Register()
	RootLogger().Infof("server[%s, %s] register", serviceInfo.Name, serviceInfo.ID)
}

// Deregister 摘除服务注册信息并停止服务
func (r *EtcdV3Registrar) Deregister() {
	DoStopHook()
	for _, server := range r.servers {
		r.deregister(server.GetServiceInfo())
		// 停止服务
		if err := server.Stop(); err != nil {
			RootLogger().Error(
				"server stop err",
				zap.String("name", server.GetServiceInfo().Name),
				zap.Error(err))
		}
		RootLogger().Info(
			"server stop gracefully",
			zap.String("name", server.GetServiceInfo().Name),
		)
	}
}

func (r *EtcdV3Registrar) deregister(serviceInfo *ServiceInfo) {
	// 摘除服务节点
	r.delegate[serviceInfo.ID].Deregister()
	RootLogger().Infof("server[%s, %s] deregister", serviceInfo.Name, serviceInfo.ID)
}

// NewDefaultEtcdV3Client 创建默认 etcdv3.Client
// etcd 地址通过 GetDefaultEtcdAddress 方法获得
func NewDefaultEtcdV3Client() etcdv3.Client {
	return MustNewEtcdV3Client(GetDefaultEtcdAddress())
}

// MustNewEtcdV3Client 创建 etcdv3.Client
func MustNewEtcdV3Client(address []string) etcdv3.Client {
	options := etcdv3.ClientOptions{
		DialTimeout:   time.Second,
		DialKeepAlive: time.Second * 30,
	}
	client, err := etcdv3.NewClient(context.Background(), address, options)
	if err != nil {
		RootLogger().Panic("new etcdv3 client err", zap.Error(err))
	}
	return client
}

// Selector 节点查询
type Selector interface {
	Next() (*ServiceInfo, error)
}

// EtcdV3Selector 基于etcd3的节点查询器
type EtcdV3Selector struct {
	mtx                sync.RWMutex
	serviceName        string
	client             etcdv3.Client
	instancer          sd.Instancer
	instances          []string
	balancer           Balancer[ServiceInfo]
	ch                 chan sd.Event // 接收服务节点变更事件
	invalidateDeadline time.Time     // 定时刷新
	invalidateTimeout  time.Duration // 缓存失效时间间隔
}

// GetEtcdV3Selector 创建 EtcdV3Selector
func GetEtcdV3Selector(serviceName string) *EtcdV3Selector {
	if selector, ok := etcdV3SelectorCache[serviceName]; ok {
		return selector
	}
	lock := etcdV3SelectorCacheLock.getLock(serviceName)
	lock.Lock()
	defer lock.Unlock()
	client := NewDefaultEtcdV3Client()
	prefix := path.Join(servicePrefix, serviceName)
	instancer, err := etcdv3.NewInstancer(client, prefix, NewKitLogger("selector-instancer", logger.InfoLevel))
	if err != nil {
		RootLogger().Panic("new etcdv3 instancer err", zap.Error(err))
	}
	s := &EtcdV3Selector{
		serviceName:       serviceName,
		client:            client,
		instancer:         instancer,
		balancer:          NewRoundRobinBalancer[ServiceInfo](),
		ch:                make(chan sd.Event),
		invalidateTimeout: time.Minute,
	}
	go s.receive()
	instancer.Register(s.ch)
	etcdV3SelectorCache[serviceName] = s
	return s
}

// receive 接收 etcd 服务节点变更事件
func (s *EtcdV3Selector) receive() {
	for event := range s.ch {
		s.update(event)
	}
}

// 更新服务节点
func (s *EtcdV3Selector) update(event sd.Event) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	if event.Err == nil {
		s.updateCache(event.Instances)
		return
	}
	RootLogger().Error("receive err event", zap.Error(event.Err))
	s.invalidateCache()
}

func (s *EtcdV3Selector) invalidateCache() {
	// 缓存立即过期
	s.invalidateDeadline = time.Now()
}

// updateCache 更新节点缓存，因为没有加锁，所以不要直接调用这个方法，统一通过 update 方法来更新
func (s *EtcdV3Selector) updateCache(instances []string) {
	RootLogger().Info("instance update", zap.Any("instances", instances))
	s.instances = instances
	// 增加缓存失效时间
	s.invalidateDeadline = time.Now().Add(s.invalidateTimeout)
	var services []*ServiceInfo
	for _, instance := range instances {
		u, err := url.Parse(instance)
		if err != nil {
			RootLogger().Error("parse instance err", zap.String("instance", instance), zap.Error(err))
			continue
		}
		info := &ServiceInfo{}
		err = json.FromJson(u.Query().Get("info"), info)
		if err != nil {
			RootLogger().Error("decode instance err", zap.String("instance", instance), zap.Error(err))
			continue
		}
		services = append(services, info)
	}
	s.balancer.Refresh(services)
}

// Next 返回一个服务节点信息
func (s *EtcdV3Selector) Next() (*ServiceInfo, error) {
	s.mtx.RLock()
	if time.Now().Before(s.invalidateDeadline) {
		defer s.mtx.RUnlock()
		return s.balancer.Pick()
	}
	s.mtx.RUnlock()
	// 超过缓存失效时间则重新查询一次etcd来刷新缓存
	prefix := path.Join(servicePrefix, s.serviceName)
	instances, err := s.client.GetEntries(prefix)
	if err != nil {
		return nil, err
	}
	s.update(sd.Event{Instances: instances})
	return s.balancer.Pick()
}

// Close 关闭查询器
func (s *EtcdV3Selector) Close() {
	s.instancer.Deregister(s.ch)
	close(s.ch)
}

// Balancer 选择器负载策略
type Balancer[T any] interface {
	// Refresh 刷新服务节点
	Refresh(services []*T)
	// Pick 选择服务节点
	Pick() (*T, error)
}

// RoundRobinBalancer 基于 RoundRobin 算法的 Balancer 实现
type RoundRobinBalancer[T any] struct {
	mtx      sync.RWMutex
	services []*T
	index    uint64
}

// NewRoundRobinBalancer 创建 NewRoundRobinBalancer 实例
func NewRoundRobinBalancer[T any]() *RoundRobinBalancer[T] {
	return &RoundRobinBalancer[T]{}
}

// Refresh 刷新节点列表
func (b *RoundRobinBalancer[T]) Refresh(services []*T) {
	b.mtx.Lock()
	defer b.mtx.Unlock()
	b.services = services
	b.index = 0
}

// Pick 选择一个服务节点
func (b *RoundRobinBalancer[T]) Pick() (*T, error) {
	b.mtx.RLock()
	if len(b.services) == 0 {
		b.mtx.RUnlock()
		return nil, ErrNoServer
	}
	b.mtx.RUnlock()
	old := atomic.AddUint64(&b.index, 1) - 1
	idx := old % uint64(len(b.services))
	return b.services[idx], nil
}

// SelectorBuilder Selector 构造函数
type SelectorBuilder func(serviceName string) Selector
