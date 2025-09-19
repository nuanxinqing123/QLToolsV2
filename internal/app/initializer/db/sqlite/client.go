package sqlite

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nuanxinqing123/QLToolsV2/internal/app/config"
	"gorm.io/driver/sqlite"
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

	// 根据配置设置日志级别
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
	return cfg
}

// GormSQLite 初始化SQLite数据库连接
// 创建并配置SQLite数据库连接，返回gorm.DB实例
func GormSQLite() *gorm.DB {
	s := config.Config.DB

	// SQLite数据库文件路径处理
	var dbPath string
	if s.Name == "" {
		// 如果没有指定数据库名，使用默认路径
		dbPath = "./data/QLToolsV2.db"
	} else {
		// 如果指定了数据库名，检查是否为绝对路径
		if filepath.IsAbs(s.Name) {
			dbPath = s.Name
		} else {
			// 相对路径，放在data目录下
			dbPath = filepath.Join("./data", s.Name)
		}
	}

	// 确保数据库文件目录存在
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		fmt.Printf("创建SQLite数据库目录失败: %v\n", err)
		return nil
	}

	// 调试模式下打印数据库路径
	if config.Config.App.Mode == gin.DebugMode {
		fmt.Println("SQLite数据库路径:", dbPath)
	}

	// 尝试建立数据库连接
	if db, err := gorm.Open(sqlite.Open(dbPath), orm.Config(s.Prefix, s.Singular)); err != nil {
		fmt.Printf("SQLite连接失败: %v\n", err)
		return nil
	} else {
		// SQLite连接池配置
		sqlDB, _ := db.DB()
		// SQLite通常不需要太多连接，因为它是文件数据库
		maxIdle := s.MaxIdleConns
		maxOpen := s.MaxOpenConns

		// 为SQLite设置合理的默认值
		if maxIdle <= 0 {
			maxIdle = 1
		}
		if maxOpen <= 0 {
			maxOpen = 1
		}

		sqlDB.SetMaxIdleConns(maxIdle) // 设置空闲连接池中的最大连接数
		sqlDB.SetMaxOpenConns(maxOpen) // 设置打开数据库连接的最大数量
		return db
	}
}
