# 配置加载

## 方法说明

```go
// MustLoadConfig 加载配置，异常时 panic
func MustLoadConfig[T any](files ...string) T
```
`luchen` 提供了配置文件加载辅助方法，支持泛型，内部使用 `github.com/spf13/viper` 来加载配置文件。读取配置异常时程序会 panic。

支持加载多个配置文件，当多个文件配置存在相同配置 key 时，后加载的将会覆盖之前的配置。

## 参考示例

```go
package config

import (
	"os"

	"github.com/fengjx/go-halo/fs"

	"github.com/fengjx/luchen"
)

var appConfig AppConfig

type AppConfig struct {
	Server Server `json:"server"`
}

type Server struct {
	HTTP HTTPServerConfig
	GRPC GRPCServerConfig
}

type HTTPServerConfig struct {
	ServerName string `json:"server-name"`
	Listen     string `json:"listen"`
}

type GRPCServerConfig struct {
	ServerName string `json:"server-name"`
	Listen     string `json:"listen"`
}

func init() {
	var configFile string
	envConfigPath := os.Getenv("APP_CONFIG")
	if envConfigPath != "" {
		configFile = envConfigPath
	}
	if configFile == "" && len(os.Args) > 1 {
		configFile = os.Args[1]
	}
	if configFile == "" {
		confFile, err := fs.Lookup("conf/app.yaml", 3)
		if err != nil {
			luchen.RootLogger().Panic("config file not found")
		}
		configFile = confFile
	}
	configFile, err := fs.Lookup(configFile, 3)
	if err != nil {
		panic(err)
	}
	appConfig = luchen.MustLoadConfig[AppConfig](configFile)
}

func GetConfig() AppConfig {
	return appConfig
}

```


