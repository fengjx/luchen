package luchen

import (
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

// MustLoadConfig 加载配置，异常时 panic
func MustLoadConfig[T any](files ...string) T {
	viperConfig := viper.New()
	for _, file := range files {
		mergeConfig(viperConfig, file)
	}
	cfg := new(T)
	err := viperConfig.Unmarshal(&cfg, func(decoderConfig *mapstructure.DecoderConfig) {
		decoderConfig.TagName = "json"
	})
	if err != nil {
		panic(err)
	}
	return *cfg
}

func mergeConfig(viperConfig *viper.Viper, configFile string) {
	viperConfig.SetConfigFile(configFile)
	err := viperConfig.MergeInConfig()
	if err != nil {
		panic(err)
	}
}
