# 服务注册

目前仅支持使用 etcd v3 作为注册中心。要把服务注册到 etcd 非常简单，只需要再单体服务的基础上增加 register 即可。

```go
registrar = luchen.NewEtcdV3Registrar(
    grpc.GetServer(),
    http.GetServer(),
)
registrar.Register()

// 摘除并停止服务
registrar.Deregister()
```

详细可查看：[_example/greetsvr/transport/server.go](https://github.com/fengjx/luchen/blob/dev/_example/greetsvr/transport/server.go)

etcd 地址读取方式

1. 通过环境变量`LUCHEN_ETCD_ADDRESS`设置
2. 调用`luchen.SetDefaultEtcdAddress(etcdAddrs)`

详细可以查看源码[env.go](https://github.com/fengjx/luchen/blob/dev/env.go)

