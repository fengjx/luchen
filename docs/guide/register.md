# 服务注册&发现

## 注册中心

注册中心目前只支持 etcd v3。

安装方法见官方文档：<https://etcd.io/docs/v3.5/install/>

docker 安装 etcd：<https://etcd.io/docs/v3.5/op-guide/container/>

## 服务注册

目前仅支持使用 etcd v3 作为注册中心。要把服务注册到 etcd 非常简单，只需要再单体服务的基础上增加 register 即可。

```go
hs := luchen.NewHTTPServer(
    luchen.WithServiceName("hello"),
)
gs := luchen.NewGRPCServer(
    luchen.WithServiceName("rpc.hello"),
)
registrar := luchen.NewEtcdV3Registrar(
    hs,
    gs,
)
// 注册并启动服务
registrar.Register()

// 摘除并停止服务
registrar.Deregister()
```

详细可查看：[_example/greetsvr/transport/server.go](https://github.com/fengjx/luchen/blob/dev/_example/greetsvr/transport/server.go)

etcd 地址设置

1. 通过环境变量`LUCHEN_ETCD_ADDRESS`设置
2. 调用`env.SetDefaultEtcdAddress(etcdAddrs)`

## 服务发现

接口定义，Selector 接口负责服务节点查询
```go
// Selector 服务节点查询
type Selector interface {
	Next() (*ServiceInfo, error)
}
```

根据服务名获取 Selector
```go
selector := luchen.GetEtcdV3Selector(serviceName)
// 获取服务节点
serviceInfo, err := selector.Next()
```

在大多数情况下，都不需要直接使用 `Selector`， 除非你有定制化需求。

`GRPCClient` 和 `HTTPClient` 内部会通过 Selector 获取服务节点，并向目标节点发起rpc请求。具体查看[客户端章节](/guide/http-client)


## 示例源码

完整示例代码：[quickstart](https://github.com/fengjx/luchen/tree/master/_example/quickstart)