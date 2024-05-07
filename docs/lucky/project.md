# 工程说明

## 后端工程

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


