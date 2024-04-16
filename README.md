# luchen 基于 go-kit 的微服务框架

luchen 是一个基于 [go-kit](https://github.com/go-kit/kit) 封装的微服务框架，支持微服务架构也能单体服务运行。秉承 go-kit 的分层设计思想，同时集成了丰富的工程实践设计。它提供了一套完整的解决方案，旨在简化服务开发、提高开发效率，让开发者更专注于业务逻辑的实现。

## 特性

- 微服务架构支持： 使用 go-kit 实现微服务化架构，支持服务注册、发现、负载均衡、限流、熔断等功能。
- 多协议支持： 支持 HTTP、gRPC 传输协议，适用于不同的场景和需求，轻松扩展更多协议支持，无需改动业务层代码。
- 快速开发： 封装了工程实践中常用的组件和工具，可以快速构建和部署服务。
- 分层设计： 秉承 go-kit 的分层设计思想，包括端点（Endpoints）、传输（Transport）、服务（Service）等层次，保证了代码的可维护性和可扩展性。

> ps: 开始这个项目的时候，我女儿刚出生，取名【路辰】，所以将项目名命名为`luchen`。

## 快速开始

启动 helloworld 服务
```bash
$ cd _example
$ go run helloworld/main.go
```

source: [_example/helloworld/main.go](_example/helloworld/main.go)

```go
func main() {
	httpSvr := luchen.NewHTTPServer(
		luchen.WithServiceName("helloworld"),
		luchen.WithServerAddr(":8080"),
	).Handler(
		&helloHandler{},
	)
	luchen.Start(httpSvr)

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGKILL)

	<-quit
	luchen.Stop()
}

type helloHandler struct {
}

func (h *helloHandler) Bind(router *luchen.HTTPServeMux) {
	router.Handle("/say-hello", h.sayHello())
}

func (h *helloHandler) sayHello() *httptransport.Server {
	return httptransport.NewServer(
		makeSayHelloEndpoint(),
		decodeSayHello,
		encodeSayHello,
	)
}

func makeSayHelloEndpoint() kitendpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		name := request.(string)
		response = "hello: " + name
		return
	}
}

func decodeSayHello(_ context.Context, r *http.Request) (interface{}, error) {
	name := r.URL.Query().Get("name")
	return name, nil
}

func encodeSayHello(_ context.Context, w http.ResponseWriter, resp interface{}) error {
	_, err := w.Write([]byte(resp.(string)))
	return err
}
```

请求服务接口
```bash
$ curl http://localhost:8080/say-hello?name=fjx
hello: fjx
```

### 文档和示例

参考文档: <http://luchen.fun>

- [helloworld](_example/helloworld) 简单示例
- [greetsvc](_example/greetsvc) http + grpc 微服务示例
- [gateway](_example/gateway) 网关服务示例


## 作者

![个人微信](docs/public/assets/img/wx.jpg)

- blog: <http://blog.fengjx.com>
- email: fengjianxin2012@gmail.com

