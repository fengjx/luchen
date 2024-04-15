# luchen

基于 [go-kit](https://github.com/go-kit/kit) 封装的微服务框架

1. 多协议支持（目前支持 http 和 grpc）
2. 秉承 go-kit 的简单，支持单体服务和微服务，可以自己选择
3. 在单体服务的基础上，只需要增加一个 Register 即可完成服务注册
4. 实现了支持静态路由和动态服务发现网关服务，通过插件化很容易对功能进行扩展

参考文档: <https://luchen.fun>

> 开始这个项目的时候，我女儿刚出生，取名【路辰】，所以将项目名命名为`luchen`。

## 快速体验

启动 helloworld 服务
```bash
$ cd _example
$ go run helloworld/main.go
```

source: [_example/helloworld/main.go](_example/helloworld/main.go)

请求服务接口
```bash
$ curl http://localhost:8080/say-hello?name=fjx
hello: fjx
```

## 示例参考

- [helloworld](_example/helloworld) 简单示例
- [greetsvc](_example/greetsvc) http + grpc 微服务示例
- [gateway](_example/gateway) 网关服务示例


## 作者

![个人微信](docs/public/assets/img/wx.jpg)

- blog: <http://blog.fengjx.com>
- email: fengjianxin2012@gmail.com

