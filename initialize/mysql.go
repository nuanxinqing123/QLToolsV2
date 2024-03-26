package initialize

import (
	"fmt"
	"log"
	"os"
	"time"

	"QLToolsV2/config"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type Writer struct {
	logger.Writer
}

// NewWriter writer 构造函数
func NewWriter(w logger.Writer) *Writer {
	return &Writer{Writer: w}
}

// Printf 格式化打印日志
func (w *Writer) Printf(message string, data ...interface{}) {
	mode := config.GinConfig.App.Mode
	if mode == "release" {
		config.GinLOG.Info(fmt.Sprintf(message+"\n", data...))
	} else {
		w.Writer.Printf(message, data...)
	}
}

type DbBase interface {
	GetLogMode() string
}

var orm = new(_gorm)

type _gorm struct{}

// Config gorm 自定义配置
func (g *_gorm) Config() *gorm.Config {
	cfg := &gorm.Config{
		// 命名策略
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 是否使用单数形式的表名，如果设置为 true，那么 User 模型会对应 users 表
		},

		DisableForeignKeyConstraintWhenMigrating: true,
	}
	_default := logger.New(NewWriter(log.New(os.Stdout, "\r\n", log.LstdFlags)), logger.Config{
		SlowThreshold: 200 * time.Millisecond,
		LogLevel:      logger.Warn,
		Colorful:      true,
	})

	cfg.Logger = _default.LogMode(logger.Info)
	return cfg

}

func GormSQLite() *gorm.DB {
	if db, err := gorm.Open(sqlite.Open("config/app.db"), orm.Config()); err != nil {
		return nil
	} else {
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetMaxOpenConns(50)
		return db
	}
}
