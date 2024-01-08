package luchen

// GatewayConfig 网关配置
type GatewayConfig struct {
	ServerName      string  `json:"server-name"`             // 服务名
	Listen          string  `json:"listen"`                  // 监听地址
	WebsocketPrefix string  `json:"websocket-prefix-listen"` // websocket 连接路径
	DiscoveryPrefix string  `json:"discovery-prefix"`        // 自动服务发现路由前缀
	Routes          []Route `json:"routes"`                  // 静态路由
}

// Route 服务路由
type Route struct {
	Protocol     string            `json:"protocol"`      // 协议，暂时只支持 http
	Pattern      string            `json:"pattern"`       // 匹配模式：path, host
	Prefix       string            `json:"prefix"`        // 匹配前缀
	Host         string            `json:"host"`          // 匹配 host
	ServiceName  string            `json:"service-name"`  // 注册的服务名
	Upstream     string            `json:"upstream"`      // 上游服务
	RewriteRegex string            `json:"rewrite-regex"` // url 重写正则
	Weight       int               `json:"weight"`        // 权重
	Ext          map[string]string `json:"ext"`           // 扩展信息
}
