# HTTP Client

## 请求单体服务接口

调用单体服务直接使用 go 标准库 `http.Client` 即可。

## 请求微服务接口

```go
client := luchen.GetHTTPClient("greeter")
body, _ := json.ToBytes(&pb.HelloRequest{
    Name: "fengjx",
})
req := &luchen.HTTPRequest{
    Path:   "/say-hello",
    Method: http.MethodPost,
    Body:   body,
}
response, err := client.Call(context.Background(), req)
if err != nil {
    log.Fatal(err)
}
if !response.IsSuccess() {
    log.Fatal("http call not success")
}
log.Println(response.String())
```

参考源码：[greeterhttpcli](https://github.com/fengjx/luchen/tree/dev/_example/greetsvr/test/greeterhttpcli/main.go)




