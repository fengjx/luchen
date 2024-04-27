package config

import (
	"os"

	"github.com/fengjx/go-halo/fs"
	"github.com/fengjx/luchen/log"

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
		confFile, err := fs.Lookup("conf/app.yml", 3)
		if err != nil {
			log.Panic("config file not found")
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
