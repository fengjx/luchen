# http 服务

## 创建一个 http server

```go
httpServer := luchen.NewHTTPServer(
    "helloworld",
    luchen.WithHTTPAddr(":8080"),
)

```

```go
// NewHTTPServer 创建 http server
// opts 查看 ServerOptions
func NewHTTPServer(opts ...ServerOption) *HTTPServer

// ServerOptions server 选项
type ServerOptions struct {
    serviceName string
    addr        string
    metadata    map[string]any
}
```

## 路由和端点绑定

```go
handler := &helloHandler{}
httpServer.Handler(
    handler,
)
```

handler 实现
```go
type greeterHandler struct {
}

func (h *greeterHandler) Bind(router luchen.HTTPRouter) {
    router.Route("/hello", func(r chi.Router) {
        r.Handle("/say-hello", h.sayHello())
    })
}

func (h *greeterHandler) sayHello() *httptransport.Server {
    options := []httptransport.ServerOption{
        httptransport.ServerErrorEncoder(errorEncoder),
    }
    return luchen.NewHTTPHandler(
        hello.GetInst().Endpoints.MakeSayHelloEndpoint(),
        luchen.DecodeParamHTTPRequest[pb.HelloReq],
        luchen.CreateHTTPJSONEncoder(httpResponseWrapper),
        options...,
    )
}
```

## 参数解析

`luchen` 对 http 参数和响应编解码简单封装了辅助方法。

解码
```go
// DecodeParamHTTPRequest 解析 http request query 和 form 参数
func DecodeParamHTTPRequest[T any](ctx context.Context, r *http.Request) (interface{}, error)

// DecodeJSONRequest 解析 http request body json 参数
func DecodeJSONRequest[T any](ctx context.Context, r *http.Request) (interface{}, error)
```

编码
```go
// CreateHTTPJSONEncoder http 返回json数据
// wrapper 对数据重新包装
func CreateHTTPJSONEncoder(wrapper DataWrapper) httptransport.EncodeResponseFunc
```

## 端点

```go
func makeSayHelloEndpoint() kitendpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
        // 业务逻辑处理
		name := request.(string)
		response = "hello: " + name
		return
	}
}
```

完整示例源码：[helloworld](https://github.com/fengjx/luchen/tree/dev/_example/helloworld)
