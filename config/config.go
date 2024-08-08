package config

import (
	"github.com/bluele/gcache"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"QLToolsV2/config/autoload"
)

type Configuration struct {
	App autoload.App `mapstructure:"app" json:"app" yaml:"app"`
	JWT autoload.JWT `mapstructure:"jwt" json:"jwt" yaml:"jwt"`
}

var (
	GinConfig Configuration
	GinDB     *gorm.DB
	GinCache  gcache.Cache
	GinLOG    *zap.Logger
	GinVP     *viper.Viper
)
