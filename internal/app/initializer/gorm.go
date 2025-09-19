package initializer

import (
	"os"

	"github.com/nuanxinqing123/QLToolsV2/internal/app/config"
	"github.com/nuanxinqing123/QLToolsV2/internal/app/initializer/db/mysql"
	"github.com/nuanxinqing123/QLToolsV2/internal/app/initializer/db/postgresql"
	"github.com/nuanxinqing123/QLToolsV2/internal/app/initializer/db/sqlite"
	_const "github.com/nuanxinqing123/QLToolsV2/internal/const"
	"gorm.io/gorm"
)

func Gorm() *gorm.DB {
	switch config.Config.DB.Type {
	case _const.Mysql:
		return mysql.GormMysql()
	case _const.Postgres:
		return postgresql.GormPostgreSQL()
	case _const.SQLite:
		return sqlite.GormSQLite()
	default:
		// 默认使用SQLite
		return sqlite.GormSQLite()
	}
}

// RegisterTables 注册数据库表专用
func RegisterTables(db *gorm.DB) {
	// 数据表：自动迁移
	err := db.Set("gorm:table_options", "CHARSET=utf8mb4").AutoMigrate()
	if err != nil {
		os.Exit(0)
	}
}
