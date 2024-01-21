# gRPC Client

## 请求单体服务接口

直接使用grpc client原生调用方式即可。

```go
clientConn, err := grpc.Dial(
    "localhost:8090",
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

参考源码：[greetergrpccli](https://github.com/fengjx/luchen/tree/dev/_example/greetsvr/test/greetergrpccli/main.go)

## 请求微服务接口

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
参开源码：[greetergrpcmicroclient](https://github.com/fengjx/luchen/tree/dev/_example/greetsvr/test/greetergrpcmicroclient/main.go)

