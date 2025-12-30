package config

import (
	"github.com/bluele/gcache"
	jsoniter "github.com/json-iterator/go"
	"github.com/nuanxinqing123/QLToolsV2/internal/app/config/autoload"
	"github.com/nuanxinqing123/QLToolsV2/internal/data/ent"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Configuration struct {
	App   autoload.App   `mapstructure:"app" json:"app" yaml:"app"`
	DB    autoload.DB    `mapstructure:"db" json:"db" yaml:"db"`
	Cache autoload.Cache `mapstructure:"cache" json:"cache" yaml:"cache"`
}

var (
	Config Configuration
	Log    *zap.Logger
	Ent    *ent.Client
	Cache  gcache.Cache
	JSON   jsoniter.API
	VP     *viper.Viper
)
