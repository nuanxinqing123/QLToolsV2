package config

import (
	"github.com/nuanxinqing123/QLToolsV2/internal/app/config/autoload"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Configuration struct {
	App   autoload.App   `mapstructure:"app" json:"app" yaml:"app"`
	DB    autoload.DB    `mapstructure:"db" json:"db" yaml:"db"`
	Cache autoload.Cache `mapstructure:"cache" json:"cache" yaml:"cache"`
}

var (
	Config Configuration
	Log    *zap.Logger
	DB     *gorm.DB
	Cache  *redis.Client
	VP     *viper.Viper
)
