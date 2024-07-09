package env

import (
	"os"
	"path/filepath"

	"github.com/fengjx/go-halo/utils"
)

var (
	_appPath string

	defaultEtcdAddress = []string{"localhost:2379"}
)

func init() {
	appPath, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	_appPath = appPath
}

const (
	// Local 本地环境
	Local = "local"
	// Dev 开发环境
	Dev = "dev"
	// Test 测试环境
	Test = "test"
	// Prod 生产环境
	Prod = "prod"
)

// GetEnv 返回当前环境
func GetEnv() string {
	env := os.Getenv("APP_ENV")
	if env == "" {
		return Local
	}
	return env
}

// Is 返回是否是指定环境
// 框架定义了四种环境：local、dev、test、prod，这个方法可以判断这四种环境意外的情况
func Is(env string) bool {
	return GetEnv() == env
}

// IsProd 返回是否是生产环境
func IsProd() bool {
	return GetEnv() == Prod
}

// IsTest 返回是否是测试环境
func IsTest() bool {
	return GetEnv() == Test
}

// IsDev 返回是否是开发环境
func IsDev() bool {
	return GetEnv() == Dev
}

// IsLocal 返回是否是本地环境
func IsLocal() bool {
	return GetEnv() == Local
}

// GetAppName 可执行文件名
func GetAppName() string {
	app := filepath.Base(os.Args[0])
	return app
}

// GetAppPath 可执行文件路径
func GetAppPath() string {
	return _appPath
}

// GetDefaultEtcdAddress 返回 etcd 连接地址
func GetDefaultEtcdAddress() (address []string) {
	etcdAddr := os.Getenv("LUCHEN_ETCD_ADDRESS")
	if len(etcdAddr) > 0 {
		return utils.SplitTrim(etcdAddr, ",")
	}
	return defaultEtcdAddress
}

// SetDefaultEtcdAddress 覆盖全局 etcd 地址
func SetDefaultEtcdAddress(address []string) {
	defaultEtcdAddress = address
}
