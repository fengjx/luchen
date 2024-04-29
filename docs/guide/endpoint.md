# 端点定义

一个端点就是对应一个接口实现，端点是对接口的一层抽象，不同协议对外暴露的api都将通过编解码方式转换为统一的请求和响应参数，来实现多协议支持。

相关概念可以查看[简介](/guide/introduction)章节


接口定义
```go
type Endpoint func(ctx context.Context, request interface{}) (response interface{}, err error)
```


## 创建一个端点

```go
func makeSayHelloEndpoint() luchen.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		name := request.(string)
		response = "hello: " + name
		return
	}
}
```

端点需要与协议进行绑定，可以查看对应协议章节

- [http 协议绑定](/guide/http-server#端点绑定)
- [grpc 协议绑定](/guide/grpc-server#端点绑定)


