package config

import (
	"flag"
	"os"

	"github.com/fengjx/go-halo/fs"
	"github.com/fengjx/luchen/log"
	"go.uber.org/zap"

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
	configArg := flag.String("c", "", "custom config file path")
	flag.Parse()

	appConfigFile, err := fs.Lookup("conf/app.yml", 5)
	if err != nil {
		log.Panic("config file not found")
	}
	configs := []string{appConfigFile}
	var configFile string
	if configArg != nil {
		configFile = *configArg
	}
	if configFile == "" {
		envConfig := os.Getenv("APP_CONFIG")
		if envConfig != "" {
			configFile = envConfig
		}
	}
	if configFile != "" {
		configFile, err = fs.Lookup(configFile, 5)
		if err != nil {
			log.Panic("config file not found", zap.String("file", configFile), zap.Error(err))
		}
		log.Infof("load custom config file: %s", configFile)
		configs = append(configs, configFile)
	}
	appConfig = luchen.MustLoadConfig[AppConfig](configs...)
}

func GetConfig() AppConfig {
	return appConfig
}
