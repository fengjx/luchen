# luchen

luchen 是一个基于 [go-kit](https://github.com/go-kit/kit) 封装的微服务框架。秉承 go-kit 的分层设计思想，同时集成了丰富的工程实践设计。它提供了一套完整的解决方案，旨在简化服务开发、提高开发效率，让开发者更专注于业务逻辑的实现。

无论是构建复杂的微服务系统还是简单的单体应用，[luchen](https://github.com/fengjx/luchen) 都能满足你的需求。

> ps: 开始这个项目的时候，我女儿刚出生，取名【路辰】，所以将项目名命名为`luchen`。

## 特性

- 快速开发： 封装了工程实践中常用的组件和工具，可以快速构建和部署服务。
- 多协议支持： 支持 HTTP、gRPC 传输协议，适用于不同的场景和需求，轻松扩展更多协议支持，无需改动业务层代码。
- 分层设计： 保留 go-kit 的分层设计思想，包括端点（Endpoints）、传输（Transport）、服务（Service）等层次，保证了代码的可维护性和可扩展性。
- 微服务支持： 使用 go-kit 实现微服务化架构，支持服务注册、发现、负载均衡、限流、熔断等功能。

## 快速体验

启动 helloworld 服务
```bash
$ cd _example/helloworld
$ go run main.go
```

请求服务接口
```bash
$ curl http://localhost:8080/say-hello?name=fjx
hello: fjx
```

### 示例

- [helloworld](_example/helloworld) 简单示例
- [feathttp](_example/feathttp) http 功能特性示例
- [featgrpc](_example/featgrpc) grpc 功能特性示例
- [quickstart](_example/quickstart) 多协议支持示例
- [gateway](_example/gateway) 网关服务示例

## 文档

详细的文档请查阅: <http://luchen.fun>

## 生态

- [lca](https://github.com/fengjx/lca) 基于 amis 开发的低代码后台系统
- [glca](https://github.com/fengjx/lca) 基于`luchen`开发 `lca` 后台服务实现
- [lc](https://github.com/fengjx/lc) glca 命令行工具，可以快速生成通用代码

## 作者

![个人微信](docs/public/assets/img/wx.jpg)

- blog: <http://blog.fengjx.com>
- email: fengjianxin2012@gmail.com

