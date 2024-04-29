# http 服务

## 创建一个 http server

```go
httpServer := luchen.NewHTTPServer(
    "helloworld",
    luchen.WithHTTPAddr(":8080"),
)
```

可选参数
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

// WithServiceName server 名称，在微服务中作为一组服务名称标识，单体服务则无需关注
func WithServiceName(serviceName string) ServerOption
// WithServerAddr server 监听地址
func WithServerAddr(addr string) ServerOption
// WithServerMetadata server 注册信息 metadata，单体服务无需关注
func WithServerMetadata(md map[string]any) ServerOption
```

## 路由和端点绑定

handler 接口定义
```go
// HTTPHandler http 请求处理器接口
type HTTPHandler interface {
	// Bind 绑定路由
	Bind(router *HTTPServeMux)
}
```

示例
```go
type helloHandler struct {
}

func (h *helloHandler) Bind(router *luchen.HTTPServeMux) {
    router.Handle("/say-hello", h.sayHello())
}

func (h *helloHandler) sayHello() *httptransport.Server {
    return luchen.NewHTTPTransportServer(
        makeSayHelloEndpoint(), // 端点绑定，端点的定义在下面说明
        decodeSayHello,
        encodeSayHello,
    )
}
```

注册路由
```go
handler := &helloHandler{}
httpServer.Handler(
    handler,
)
```

## 中间件

接口定义
```go
// HTTPMiddleware http 请求中间件
type HTTPMiddleware func(http.Handler) http.Handler
```

示例：实现提个请求耗时打印中间件
```go
func timeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestURI := r.RequestURI
		start := time.Now()
		defer func() {
			log.Infof("take time: %s, %v", requestURI, time.Since(start))
		}()
		next.ServeHTTP(w, r)
	})
}
```

```go
// 使用中间件
httpServer.Use(
    timeMiddleware,
)
```


## 端点绑定

```go
// 端点定义，端点即对应一个接口，不同协议转换成相同的参数，交给端点进行处理
func makeSayHelloEndpoint() kitendpoint.Endpoint {
    return func(ctx context.Context, request interface{}) (response interface{}, err error) {
        name := request.(string)
        response = "hello: " + name
        return
    }
}

// 处理http协议参数解码
func decodeSayHello(_ context.Context, r *http.Request) (interface{}, error) {
    name := r.URL.Query().Get("name")
    return name, nil
}

// 处理http协议响应参数编码
func encodeSayHello(_ context.Context, w http.ResponseWriter, resp interface{}) error {
    _, err := w.Write([]byte(resp.(string)))
    return err
}
```

## 参数编解码

通过编解码处理将不同协议转换为统一的结构体，交给 endpoint 处理。 

接口定义在：<https://github.com/go-kit/kit/blob/master/transport/http/encode_decode.go>
```go
// http 请求参数解码
type DecodeRequestFunc func(context.Context, *http.Request) (request interface{}, err error)

// http 响应参数编码
type EncodeRequestFunc func(context.Context, *http.Request, interface{}) error
```

`luchen` 内置了一些 http 协议的请求和响应参数编解码方法，如不满足需求，可以自己实现编解码接口。

解码
```go
// DecodeParamHTTPRequest 解析 http request query 和 form 参数
func DecodeParamHTTPRequest[T any](ctx context.Context, r *http.Request) (interface{}, error)

// DecodeJSONRequest 解析 http request body json 参数
func DecodeJSONRequest[T any](ctx context.Context, r *http.Request) (interface{}, error)
```

编码
```go
// EncodeHTTPJSONResponse http 返回json数据
// wrapper 对数据重新包装
func EncodeHTTPJSONResponse(wrapper DataWrapper) httptransport.EncodeResponseFunc
```

## 静态文件服务

```go
// 注册静态文件服务访问路径和文件路径
httpServer.Static("/assets/", "static")
```

## 示例源码

完整示例代码：[feathttp](https://github.com/fengjx/luchen/tree/master/_example/feathttp)

