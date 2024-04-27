# 工程规范

以下是根据工程实践不断调整后在我项目中使用的规范，可供参考。

目录结构
```
├── conf                         // 配置文件
├── connom                       // 公共模块
├── logic                        // 业务模块
│        ├── calc                // calc 模块
│        │        ├── calcapi    // 对其他模块暴露的 api
│        │        └── internal   // 内部依赖
│        │        ├── init.go    // 模块初始化
│        ├── hello               // hello 模块
│        │        ├── helloapi   // 对其他模块暴露的 api
│        │        └── internal   // 内部依赖
│        │        ├── init.go    // 模块初始化
│        └── init.go             // 业务模块初始化
├── pb                           // 协议
├── transport                    // 传输层处理逻辑
├── main.go                      // 程序入口
```

说明：从过往经验来看，模块之间的相互依赖是未来重构难度的根源，
我们应该把相关联的模块放在一起，保持单个模块代码简洁。

需要对其他模块暴露的接口，在 xxxapi 中定义，同时不能让 xxxapi 依赖其他代码，避免造成循环依赖。一个可行的做法是：在 xxxapi 中定义接口，并在 Init 时注入接口实现。

参考：[quickstart](https://github.com/fengjx/luchen/tree/master/_example/quickstart)