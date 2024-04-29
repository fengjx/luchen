# grpc 服务

## 创建一个 gRPC server

```go
grpcServer := luchen.NewGRPCServer(
    luchen.WithServiceName("featgrpc"),
    luchen.WithServerAddr(":8088"),
)
```

```go
// NewGRPCServer 创建 grpc server
// opts 查看 ServerOptions
func NewGRPCServer(opts ...ServerOption) *GRPCServer

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

## 根据 proto 文件生成相关代码

### 安装 grpc 编译工具

- protoc 安装：<https://grpc.io/docs/protoc-installation/>
- protoc grpc 插件安装：<https://grpc.io/docs/languages/go/quickstart/>

### 生成代码

```bash
cd _example/featgrpc/pb/

bash build.sh
```

build.sh
```bash
#!/usr/bin/env bash

protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    *.proto
```

编译后会生成 `greet.pb.go` 和 `greet_grpc.pb.go` 两个文件
```bash
$ ls    
build.sh greet.pb.go greet.proto greet_grpc.pb.go
```

## 端点绑定

需要将 grpc 接口实现与端点进行绑定。

```go
grpcSvr.RegisterService(func(gs *grpc.Server) {
    // 注册 grpc 服务
    pb.RegisterGreeterServer(gs, newGreeterServer())
})
```

GreeterServer 实现
```go
type GreeterServer struct {
    pb.UnimplementedGreeterServer
    sayHello grpctransport.Handler
}

func newGreeterServer() pb.GreeterServer {
    svr := &GreeterServer{}
	// 绑定端点，将接口实现交给 endpoint 处理
    svr.sayHello = luchen.NewGRPCTransportServer(
        makeSayHelloEndpoint(),
        luchen.DecodePB[pb.HelloReq],
        luchen.EncodePB[pb.HelloResp],
    )
    return svr
}

func (s *GreeterServer) SayHello(ctx context.Context, req *pb.HelloReq) (*pb.HelloResp, error) {
    _, resp, err := s.sayHello.ServeGRPC(ctx, req)
    if err != nil {
        return nil, err
    }
    return resp.(*pb.HelloResp), nil
}

func (s *GreeterServer) decodeSayHello(_ context.Context, req interface{}) (interface{}, error) {
    helloReq := req.(*pb.HelloReq)
    return &pb.HelloReq{
        Name: helloReq.Name,
    }, nil
}

func (s *GreeterServer) encodeSayHello(_ context.Context, resp interface{}) (interface{}, error) {
    helloResp := resp.(*pb.HelloResp)
    return &pb.HelloResp{
        Message: helloResp.Message,
    }, nil
}

func makeSayHelloEndpoint() kitendpoint.Endpoint {
    return func(ctx context.Context, request interface{}) (response interface{}, err error) {
        name := request.(string)
        response = "hello: " + name
        return
    }
}
```

## 参数编解码

通过编解码处理将不同协议转换为统一的结构体，交给 endpoint 处理。

接口定义：<https://github.com/go-kit/kit/blob/master/transport/grpc/encode_decode.go>
```go
// 请求参数解码
type DecodeRequestFunc func(context.Context, interface{}) (request interface{}, err error)

// 响应参数编码
type EncodeRequestFunc func(context.Context, interface{}) (request interface{}, err error)
```

`luchen` 对 grpc 参数和响应编解码简单封装了辅助方法，如果不满足需求可以自己实现。

解码
```go
// DecodePB protobuf 解码
func DecodePB[T any](_ context.Context, req interface{}) (interface{}, error) 
```

编码
```go
// EncodePB protobuf 编码
func EncodePB[T any](_ context.Context, resp interface{}) (interface{}, error)
```

## 示例源码

完整示例源码：[featgrpc](https://github.com/fengjx/luchen/tree/master/_example/featgrpc)

