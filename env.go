package luchen

import (
	"os"
	"path/filepath"

	"github.com/fengjx/go-halo/utils"
	"go.uber.org/zap"
)

var (
	_appPath string

	defaultEtcdAddress = []string{"localhost:2379"}
)

func init() {
	appPath, err := os.Getwd()
	if err != nil {
		RootLogger().Panic("os.Getwd() return err", zap.Error(err))
	}
	_appPath = appPath
}

type ENV string

const (
	Local ENV = "local"
	Dev   ENV = "dev"
	Test  ENV = "test"
	Prod  ENV = "prod"
)

func GetEnv() ENV {
	env := os.Getenv("APP_ENV")
	switch ENV(env) {
	case Test:
		return Test
	case Prod:
		return Prod
	case Dev:
		return Dev
	default:
		return Local
	}
}

func IsProd() bool {
	return GetEnv() == Prod
}

func IsTest() bool {
	return GetEnv() == Test
}

func IsDev() bool {
	return GetEnv() == Dev
}

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
