# luchen

基于 [go-kit](https://github.com/go-kit/kit) 封装的微服务框架，秉承 go-kit 的简单，自己选择使用微服务还是单体服务。

> 开始这个项目的时候，我女儿刚出生，取名【路辰】，所以将项目名命名为`luchen`。

参考文档: <https://luchen.fun>

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

