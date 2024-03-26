package initialize

import (
	"os"

	"gorm.io/gorm"
)

func Gorm() *gorm.DB {
	return GormSQLite()
}

// RegisterTables 注册数据库表专用
func RegisterTables(db *gorm.DB) {
	// 数据表：自动迁移
	err := db.Set("gorm:table_options", "CHARSET=utf8mb4").AutoMigrate()
	if err != nil {
		os.Exit(0)
	}
}
