package postgresql

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nuanxinqing123/QLToolsV2/internal/app/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type Writer struct {
	logger.Writer
}

// NewWriter writer 构造函数
// 创建一个新的日志写入器，用于自定义日志输出格式
func NewWriter(w logger.Writer) *Writer {
	return &Writer{Writer: w}
}

// Printf 格式化打印日志
// 根据配置决定使用zap日志还是标准输出
func (w *Writer) Printf(message string, data ...any) {
	var logZap bool
	logZap = config.Config.DB.LogZap
	if logZap {
		config.Log.Info(fmt.Sprintf(message+"\n", data...))
	} else {
		w.Writer.Printf(message, data...)
	}
}

var orm = new(_gorm)

type _gorm struct{}

// Config gorm 自定义配置
// 配置GORM的基本设置，包括命名策略和日志级别
func (g *_gorm) Config(prefix string, singular bool) *gorm.Config {
	cfg := &gorm.Config{
		// 命名策略配置
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   prefix,   // 表前缀，在表名前添加前缀
			SingularTable: singular, // 是否使用单数形式的表名
		},
		// 迁移时禁用外键约束
		DisableForeignKeyConstraintWhenMigrating: true,
	}

	// 创建默认日志配置
	_default := logger.New(NewWriter(log.New(os.Stdout, "\r\n", log.LstdFlags)), logger.Config{
		SlowThreshold: 300 * time.Millisecond, // 慢查询阈值
		LogLevel:      logger.Warn,            // 默认日志级别
		Colorful:      true,                   // 启用彩色输出
	})

	if config.Config.App.Mode == gin.DebugMode {
		// 调试模式下，根据配置设置日志级别
		switch config.Config.DB.LogLevel {
		case "silent", "Silent":
			cfg.Logger = _default.LogMode(logger.Silent)
		case "error", "Error":
			cfg.Logger = _default.LogMode(logger.Error)
		case "warn", "Warn":
			cfg.Logger = _default.LogMode(logger.Warn)
		case "info", "Info":
			cfg.Logger = _default.LogMode(logger.Info)
		default:
			cfg.Logger = _default.LogMode(logger.Info)
		}
	} else {
		// 非调试模式下，关闭SQL日志输出
		cfg.Logger = _default.LogMode(logger.Silent)
	}
	return cfg
}

// GormPostgreSQL 初始化PostgreSQL数据库连接
// 创建并配置PostgreSQL数据库连接，返回gorm.DB实例
func GormPostgreSQL() *gorm.DB {
	p := config.Config.DB
	if p.Name == "" {
		return nil
	}

	// 构建PostgreSQL DSN连接字符串
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d %s",
		p.Host, p.UserName, p.Password, p.Name, p.Port, p.Config)

	// 调试模式下打印DSN
	if config.Config.App.Mode == gin.DebugMode {
		fmt.Println("PostgreSQL DSN:", dsn)
	}

	// 配置PostgreSQL连接参数
	postgresConfig := postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true, // 禁用隐式prepared statement
	}

	// 尝试建立数据库连接
	if db, err := gorm.Open(postgres.New(postgresConfig), orm.Config(p.Prefix, p.Singular)); err != nil {
		fmt.Printf("PostgreSQL连接失败: %v\n", err)
		return nil
	} else {
		// 配置连接池参数
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(p.MaxIdleConns) // 设置空闲连接池中的最大连接数
		sqlDB.SetMaxOpenConns(p.MaxOpenConns) // 设置打开数据库连接的最大数量
		return db
	}
}
