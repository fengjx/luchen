package luchen

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/fengjx/go-halo/addr"
	"github.com/go-chi/chi/v5"
	"github.com/golang/groupcache/lru"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

var (
	//nolint:gomnd
	rewriteRegexpCache = lru.New(50)
)

// GatewayOptions 网关选项定义
type GatewayOptions struct {
	plugins []GatewayPlugin
}

// GatewayOption 网关选项赋值
type GatewayOption func(*GatewayOptions)

// WithGatewayPlugin 注册网关扩展插件
func WithGatewayPlugin(plugins ...GatewayPlugin) GatewayOption {
	return func(o *GatewayOptions) {
		o.plugins = append(o.plugins, plugins...)
	}
}

// Gateway 网关服务
type Gateway struct {
	*baseServer
	*httputil.ReverseProxy
	config     GatewayConfig
	server     *http.Server
	routes     []*httpRoute
	patternMap map[string]Pattern
}

type httpRoute struct {
	protocol     string // 协议
	pattern      string // 匹配模式
	prefix       string // 匹配前缀
	host         string // 匹配 host
	serviceName  string // 注册的服务名
	rewriteRegex string // url 重写正则表达式
	upstream     string // 上游服务
	weight       int    // 权重
}

// NewGateway 创建 gateway 服务
func NewGateway(cfg GatewayConfig, opts ...GatewayOption) *Gateway {
	options := &GatewayOptions{}
	_ = options

	var routes []*httpRoute
	// 静态路由初始化
	for _, route := range cfg.Routes {
		if route.Protocol != "http" {
			continue
		}
		routes = append(routes, &httpRoute{
			pattern:      route.Pattern,
			protocol:     route.Protocol,
			prefix:       route.Prefix,
			host:         route.Host,
			serviceName:  route.ServiceName,
			upstream:     route.Upstream,
			rewriteRegex: route.RewriteRegex,
			weight:       route.Weight,
		})
	}
	// 权重排序
	sort.SliceStable(routes, func(i, j int) bool {
		return routes[i].weight > routes[j].weight
	})

	proxy := &httputil.ReverseProxy{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout: defaultConnectionTimeout,
			}).DialContext,
			MaxIdleConnsPerHost:   defaultMaxPoolSize,
			MaxIdleConns:          defaultPoolSize,
			IdleConnTimeout:       time.Second * 30,
			ExpectContinueTimeout: defaultConnectionTimeout,
		},
	}
	g := &Gateway{
		baseServer: &baseServer{
			id:          uuid.NewString(),
			serviceName: cfg.ServerName,
			protocol:    ProtocolHTTP,
			address:     cfg.Listen,
			metadata:    make(map[string]any),
		},
		ReverseProxy: proxy,
		config:       cfg,
		routes:       routes,
		patternMap:   make(map[string]Pattern),
	}
	g.RegPattern(
		&HostPattern{},
		&PrefixPattern{},
	)
	g.Director = g.director
	g.ModifyResponse = g.modifyResponse
	g.ErrorHandler = g.errorHandler
	return g
}

// Start 启动服务
func (g *Gateway) Start() error {
	g.Lock()
	address := ":8080"
	if len(g.address) > 0 {
		address = g.address
	}
	ln, err := net.Listen("tcp", address)
	if err != nil {
		g.Unlock()
		return err
	}
	router := chi.NewRouter()
	router.Handle("/*", g)
	server := &http.Server{
		Handler: router,
	}
	g.server = server
	address = ln.Addr().String()
	host, port, err := addr.ExtractHostPort(address)
	if err != nil {
		g.Unlock()
		return err
	}
	g.address = fmt.Sprintf("%s:%s", host, port)
	g.metadata["ts"] = time.Now().UnixMilli()
	g.started = true
	RootLogger().Infof("gateway server[%s, %s] start", g.serviceName, g.id)
	g.Unlock()
	return g.server.Serve(ln)
}

// Stop 停止服务
func (g *Gateway) Stop() error {
	g.RLock()
	if !g.started {
		g.RUnlock()
		return nil
	}
	g.RUnlock()
	RootLogger().Info("gateway server stop")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return g.server.Shutdown(ctx)
}

func (g *Gateway) rewrite(req *httputil.ProxyRequest) {

}

func (g *Gateway) director(req *http.Request) {
	var r *httpRoute
	// 静态路由匹配
	for _, route := range g.routes {
		r = g.match(req, route)
		if r != nil {
			break
		}
	}

	if r == nil {
		RootLogger().Warn("no route match", zap.String("path", req.URL.Path))
		return
	}
	upstream := r.upstream
	if upstream == "" {
		serviceInfo, err := g.selectServiceInfo(r)
		if err != nil {
			RootLogger().Error("select service err",
				zap.String("service_name", r.serviceName),
				zap.Error(err),
			)
		}
		if serviceInfo != nil {
			upstream = serviceInfo.Addr
		}
	}
	if upstream == "" {
		// to write none available
		req.URL = nil
		return
	}
	req.URL.Scheme = r.protocol
	req.URL.Host = upstream
	reg := getRewriteRegexp(r.rewriteRegex)
	// url 重写
	if reg != nil {
		req.URL.Path = reg.ReplaceAllString(req.URL.Path, "$1")
	}
	RootLogger().Info("upstream info",
		zap.String("service_name", r.serviceName),
		zap.String("upstream", upstream),
		zap.String("path", req.URL.Path),
	)
	innerIP := addr.InnerIP()
	req.Header.Set("X-Real-Ip", getClientIP(req))
	req.Header.Set("X-Proxy-Server", "luchen-gateway")
	req.Header.Set("X-Upstream-Service", r.serviceName)
	req.Header.Set("X-Upstream-Node", upstream)
	req.Header.Set("X-Proxy-Ip", innerIP)
}

func (g *Gateway) modifyResponse(resp *http.Response) error {
	resp.Header.Set("Server", "luchen-gateway")
	return nil
}

func (g *Gateway) errorHandler(w http.ResponseWriter, r *http.Request, err error) {
	RootLogger().Error("handler err", zap.Error(err))
}

// RegPattern 注册匹配模式
func (g *Gateway) RegPattern(patterns ...Pattern) {
	for _, p := range patterns {
		g.patternMap[p.Name()] = p
	}
}

// match 静态路由匹配
func (g *Gateway) match(req *http.Request, route *httpRoute) *httpRoute {
	if route == nil {
		return nil
	}
	// 静态路由匹配
	if p, ok := g.patternMap[route.pattern]; ok {
		if p.Match(req, route) {
			return route
		}
	}
	RootLogger().Warn("route pattern not support", zap.String("pattern", route.pattern))
	return nil
}

// selectNode 查询服务节点
func (g *Gateway) selectServiceInfo(r *httpRoute) (*ServiceInfo, error) {
	if r.serviceName == "" {
		return nil, nil
	}
	selector := GetEtcdV3Selector(r.serviceName)
	serviceInfo, err := selector.Next()
	if err != nil {
		return nil, err
	}
	return serviceInfo, nil
}

// Pattern 路由匹配模式
type Pattern interface {
	// Name 匹配模式名称
	Name() string
	// Match 匹配路由
	Match(req *http.Request, route *httpRoute) bool
}

// HostPattern host 匹配
type HostPattern struct {
}

// Name 模式名称
func (h *HostPattern) Name() string {
	return "host"
}

// Match host 匹配判断
func (h *HostPattern) Match(req *http.Request, route *httpRoute) bool {
	return req.Host == route.host
}

// PrefixPattern 前缀匹配
type PrefixPattern struct {
}

// Name 模式名称
func (h *PrefixPattern) Name() string {
	return "prefix"
}

// Match 前缀匹配判断
func (h *PrefixPattern) Match(req *http.Request, route *httpRoute) bool {
	return strings.HasPrefix(req.URL.Path, route.prefix)
}

func getClientIP(req *http.Request) string {
	forwardIPs := getHeader(req, "X-Forwarded-For")
	if forwardIPs == "" {
		return ""
	}
	ips := strings.Split(forwardIPs, ",")
	for _, ip := range ips {
		ip = strings.TrimSpace(ip)
		if ip == "unknown" || ip == "unknow" {
			continue
		}
		parseIP := net.ParseIP(ip)
		if parseIP.IsLoopback() || parseIP.IsPrivate() {
			continue
		}
		return ip
	}
	return ""
}

func getHeader(req *http.Request, key string) string {
	vals := req.Header[key]
	if len(vals) == 0 {
		return ""
	}
	return vals[0]
}

func getRewriteRegexp(rewriteRegex string) *regexp.Regexp {
	if rewriteRegex == "" {
		return nil
	}
	var reg *regexp.Regexp
	regCache, ok := rewriteRegexpCache.Get(rewriteRegex)
	if ok {
		reg = regCache.(*regexp.Regexp)
	} else {
		var err error
		reg, err = regexp.Compile(rewriteRegex)
		if err != nil {
			RootLogger().Error("regexp compile error", zap.String("regexp", rewriteRegex), zap.Error(err))
		} else {
			rewriteRegexpCache.Add(rewriteRegex, reg)
		}
	}
	return reg
}

// GatewayPlugin 网关插件接口
type GatewayPlugin interface {
}
