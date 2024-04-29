# HTTP Client

## 通过ip:port访问

未使用服务注册的服务，可以通过固定ip:port访问，直接使用 go 标准库 `http.Client` 即可。

## 通过服务发现请求

```go
client := luchen.GetHTTPClient("quickstart")
params := url.Values{}
params.Set("name", "fengjx")
req := &luchen.HTTPRequest{
    Path:   "/hello/say-hello",
    Method: http.MethodPost,
    Params: params,
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

## 示例源码

完整示例代码：[greeterhttpcli](https://github.com/fengjx/luchen/tree/dev/_example/quickstart/cli/greeterhttpcli/main.go)

