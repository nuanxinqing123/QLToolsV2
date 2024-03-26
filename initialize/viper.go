package initialize

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"

	"QLToolsV2/config"
)

// Viper 初始化配置
func Viper() *viper.Viper {
	v := viper.New()
	v.SetConfigFile("config/config.yaml")
	v.SetConfigType("yaml")
	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	v.WatchConfig()

	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("config file changed: ", e.Name)
		if err = v.Unmarshal(&config.GinConfig); err != nil {
			fmt.Println(err)
		}
	})
	if err = v.Unmarshal(&config.GinConfig); err != nil {
		fmt.Println(err)
	}

	fmt.Println("Viper初始化成功")
	return v
}
