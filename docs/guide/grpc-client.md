# gRPC Client

## 通过ip:port访问

未使用服务注册的服务，可以通过固定ip:port访问，直接使用 grpc client 原生调用方式即可。

```go
clientConn, err := grpc.Dial(
    "localhost:8088",
    grpc.WithTransportCredentials(insecure.NewCredentials()),
)
if err != nil {
    panic(err)
}
greeterClient := pb.NewGreeterClient(clientConn)
ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
defer cancel()
helloResp, err := greeterClient.SayHello(ctx, &pb.HelloReq{
    Name: "fengjx",
})
log.Println(helloResp.Message)
```

```bash
$ cd _example/quickstart/cli/greetergrpccli
$ go run main.go

2024/04/27 15:53:52 Hi: fengjx
```

源码：[greetergrpccli](https://github.com/fengjx/luchen/tree/dev/_example/quickstart/cli/greetergrpccli/main.go)

## 通过服务发现请求

```go
grpcClient := luchen.GetGRPCClient(
    "rpc.greeter",
    grpc.WithTransportCredentials(insecure.NewCredentials()),
)
greeterClient := pb.NewGreeterClient(grpcClient)
helloResp, err := greeterClient.SayHello(context.Background(), &pb.HelloReq{
    Name: "fengjx",
})
if err != nil {
    log.Fatal(err)
}
log.Println(helloResp.Message)
```

```bash
$ cd _example/quickstart/cli/greetergrpcmicroclient
$ go run main.go

2024/04/27 15:53:00 Hi: fengjx
```

源码：[greetergrpcmicroclient](https://github.com/fengjx/luchen/tree/dev/_example/quickstart/cli/greetergrpcmicroclient/main.go)

