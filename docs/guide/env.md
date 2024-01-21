# 环境变量

## 环境变量

| 变量名                 | 说明                                                |
|---------------------|---------------------------------------------------|
| APP_ENV             | 运行环境，默认 local；可取值：local-本地，dev-开发，test-测试，prod-生产 |
| LUCHEN_ETCD_ADDRESS | etcd 地址，默认：localhost:2379                         |


## 方法说明

```go
// IsProd 返回是否是生产环境
func IsProd() bool

// IsTest 返回是否是测试环境
func IsTest() bool

// IsDev 返回是否是开发环境
func IsDev() bool

// IsLocal 返回是否是本地环境
func IsLocal() bool

// GetAppName 可执行文件名
func GetAppName() string

// GetAppPath 可执行文件路径
func GetAppPath() string

// GetDefaultEtcdAddress 返回 etcd 连接地址
func GetDefaultEtcdAddress() (address []string)

// SetDefaultEtcdAddress 覆盖全局 etcd 地址
func SetDefaultEtcdAddress(address []string) 
```

详细可查看源码 [env.go](https://github.com/fengjx/luchen/blob/dev/env.go)

