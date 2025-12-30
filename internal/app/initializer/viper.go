package initializer

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/nuanxinqing123/QLToolsV2/internal/app/config"
	"github.com/spf13/viper"
)

// Viper 初始化配置
func Viper(configPath string) *viper.Viper {
	v := viper.New()
	// 如果未指定配置文件路径，则使用默认路径
	if configPath == "" {
		configPath = "configs/config.yaml"
	}
	// 指定配置文件路径（支持绝对路径或相对路径）
	v.SetConfigFile(configPath)
	v.SetConfigType("yaml")
	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	v.WatchConfig()

	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("config file changed: ", e.Name)
		if err = v.Unmarshal(&config.Config); err != nil {
			fmt.Println(err)
		}
	})
	if err = v.Unmarshal(&config.Config); err != nil {
		fmt.Println(err)
	}

	return v
}
