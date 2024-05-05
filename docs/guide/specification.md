# 工程规范

以下是根据工程实践不断调整后在我项目中使用的规范，可供参考。

目录结构
```

├── conf                        // 配置文件
├── connom                      // 公共模块
├── logic                       // 业务逻辑
│     ├── calc                  // 业务逻辑-模块a
│     │     ├── calcpub         // 对其他模块暴露的api
│     │     ├── init.go         // 模块a初始化
│     │     └── internal        // 内部依赖
│     ├── hello                 // 业务逻辑-模块b
│     │     ├── init.go         // 模块b初始化
│     │     └── internal        // 内部依赖
│     └── init.go               // 业务逻辑初始化
├── pb                          // 协议
├── transport                   // 传输层处理逻辑
│   ├── grpc                    // 传输层grpc协议处理逻辑
│   └── http                    // 传输层http协议处理逻辑
├── main.go                     // 程序入口
```

说明：从过往经验来看，模块之间的相互依赖是未来重构难度的根源，
我们应该把相关联的模块放在一起，保持单个模块代码简洁。

参考：[quickstart](https://github.com/fengjx/luchen/tree/master/_example/quickstart)