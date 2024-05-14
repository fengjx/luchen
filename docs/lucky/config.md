# 系统配置

如果你系统一些功能配置可以在不重启服务的情况下修改，可以使用`系统配置`。

噢诶之数据保存在`sys_config`表。

## 配置定义

配置会以 key-value形式存储和获取。

- 范围：即数据可见性定义，例如：sys, backend, web, android, ios, 你可以根据需求自行定义
- 配置键：数据 key 定义
- 配置值：数据 value

## 使用

你可以在其他模块中使用`syspub.ConfigAPI`来获取配置。

提供了一下两个方法：

```go
var ConfigAPI configAPI

type configAPI interface {
	// GetConfigString 返回key对应的配置
	GetConfigString(scope string, key string) string

	// GetConfig 返回key对应的配置，并序列化成对象
	GetConfig(scope string, key string, data any) error
}
```


