# 工程说明

## 后端工程

源码：<https://github.com/fengjx/lucky>

```bash
├── conf                // 配置文件
├── connom              // 公共模块
│     ├── auth          // 登录认证
│     ├── config        // 应用配置
│     ├── errno         // 错误码定义
│     ├── kit           // 工具方法
│     ├── lifecycle     // 生命周期事件定义
│     └── types         // 公共结构体定义
├── current             // context 上下文
├── deployments         // 服务部署相关配置参考
├── integration         // 集成外部服务、中间件
│     └── db            // 数据集连接
├── logic               // 业务逻辑
│     ├── common        // 公共逻辑
│     ├── init.go       // 业务逻辑初始化
│     └── sys           // 系统功能
├── static              // 静态文件
│     └── pages         // 使用 amis 协议编写的后台页面
├── t                   // 测试相关
├── tools               // 工具
│     ├── gen           // 代码生成配置
│     └── init          // 初始化数据
└── transport           // 传输层
    └── http            // http 协议处理
```

更多工程规范，查看：<a href="/guide/specification" target="_blank">luchen推荐工程规范</a>

## 前端工程

源码：<https://github.com/fengjx/lucky-web>

前端工程是一个页面框架，包括了登录和页面渲染。绝大多数情况下，你都不需要对这个工程做修改。除非你需要定制化一些功能，例如：接入sso企业账号登录。


### 本地启动

需要安装好`nodejs`，`pnpm(推荐)`

```bash
pnpm i
pnmm run dev
```

### 打包

```bash
pnpm i
pnmm run build
```

打包好的文件在`dist`目录下，部署到静态服务器（如 nginx）即可。

如果不想另外部署一个静态服务器，还提供了一下方式直接启动一个服务

- 使用 nodejs 启动一个 http 服务器。
```bash
node server.js
```

具体可以查看[server.js](https://github.com/fengjx/lucky-web/blob/master/server.js)源码


- 直接运行一个二进制程序（后续提供）


