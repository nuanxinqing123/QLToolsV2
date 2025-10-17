package initializer

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/nuanxinqing123/QLToolsV2/internal/app/config"
	"github.com/nuanxinqing123/QLToolsV2/internal/app/initializer/db/mysql"
	"github.com/nuanxinqing123/QLToolsV2/internal/model"
	"gorm.io/gorm"
)

func Gorm() *gorm.DB {
	return mysql.GormMysql()
}

// RegisterTables 注册数据库表专用
func RegisterTables(db *gorm.DB) {
	if config.Config.App.Mode != gin.ReleaseMode {
		return
	}

	// 数据表：自动迁移
	err := db.Set("gorm:table_options", "CHARSET=utf8mb4").AutoMigrate(
		&model.CdKeys{},
		&model.Envs{},
		&model.LoginHistories{},
		&model.Panels{},
		&model.PluginExecutionLogs{},
		&model.Plugins{},
		&model.Users{},
		&model.EnvPanels{},
		&model.EnvPlugins{},
	)
	if err != nil {
		os.Exit(0)
	}
}
