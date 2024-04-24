package initialize

import (
	"os"

	"gorm.io/gorm"

	"QLToolsV2/internal/model"
)

func Gorm() *gorm.DB {
	return GormSQLite()
}

// RegisterTables 注册数据库表专用
func RegisterTables(db *gorm.DB) {
	// 数据表：自动迁移
	err := db.AutoMigrate(
		&model.User{},
		&model.Panel{},
		&model.Env{},
		&model.CdKey{})
	if err != nil {
		os.Exit(0)
	}
}
