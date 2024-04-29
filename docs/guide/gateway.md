# 网关服务（alpha）

网关目前还处于内部测试阶段，未来可能会有比较大的调整。

## 编写一个网关服务

```go
package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/fengjx/go-halo/fs"
	"github.com/fengjx/luchen"
	"go.uber.org/zap"
)

func main() {
	configFile, err := fs.Lookup("gateway.yaml", 3)
	if err != nil {
		luchen.RootLogger().Panic("config file not found", zap.Error(err))
	}
	config := luchen.MustLoadConfig[luchen.GatewayConfig](configFile)
	gateway := luchen.NewGateway(
		config,
		luchen.WithGatewayPlugin(
			&luchen.TraceGatewayPlugin{},
		),
	)
	luchen.Start(gateway)

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGKILL)

	<-quit
	luchen.Stop()
}
```

## 网关配置

```yml
server-name: "luchen-gateway"
listen: ":9000"
routes:
  - protocol: http
    pattern: prefix
    prefix: /open/api/greeter
    service-name: greeter
    rewrite-regex: "^/open/api/greeter(.*)"
```

参数说明

```go
// GatewayConfig 网关配置
type GatewayConfig struct {
	ServerName      string  `json:"server-name"`             // 服务名
	Listen          string  `json:"listen"`                  // 监听地址
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
```





