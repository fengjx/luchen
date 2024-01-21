# grpc 服务


## 根据 proto 文件生成相关代码

安装 `protoc` 和 `go`, `grpc` 插件

<https://grpc.io/docs/languages/go/quickstart/>

生成代码

```bash
cd _example/greetsvr/pb/

protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    *.proto
```

将会生成`greet.pb.go`, `greet_grpc.pb.go` 两个文件。


## 创建一个 gRPC server

```go
grpcServer := luchen.NewGRPCServer(
    luchen.WithServiceName(serviceName),
    luchen.WithServerAddr(addr),
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
```

## 注册接口实现和端点绑定

```go
grpcServer.RegisterService(func(grpcServer *grpc.Server) {
    pb.RegisterGreeterServer(grpcServer, newGreeterServer())
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
    svr.sayHello = luchen.NewGRPCHandler(
        MakeSayHelloEndpoint(),
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
```

## 参数解析

`luchen` 对 grpc 参数和响应编解码简单封装了辅助方法。

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

## 端点

```go
func MakeSayHelloEndpoint() endpoint.Endpoint {
    return func(ctx context.Context, request interface{}) (interface{}, error) {
        logger := luchen.Logger(ctx)
        logger.Info("greeter say hello")
        helloReq := request.(*pb.HelloReq)
        msg, err := GetInst().helloLogic.SayHello(ctx, helloReq.Name)
        if err != nil {
			return nil, err
        }
        return &pb.HelloResp{Message: msg}, nil
    }
}
```

完整示例源码：[greetsvr](https://github.com/fengjx/luchen/tree/dev/_example/greetsvr)

