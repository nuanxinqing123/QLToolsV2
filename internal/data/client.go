package data

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/nuanxinqing123/QLToolsV2/internal/app/config"
	"github.com/nuanxinqing123/QLToolsV2/internal/data/ent"
)

var (
	Client *ent.Client
)

// InitData 初始化数据库连接
func InitData() (*ent.Client, error) {
	var (
		drv *entsql.Driver
	)

	cfg := config.Config.DB

	switch cfg.Type {
	case "mysql":
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s", cfg.UserName, cfg.Password, cfg.Host, cfg.Port, cfg.Name, cfg.Config)
		db, err := sql.Open("mysql", dsn)
		if err != nil {
			return nil, fmt.Errorf("failed opening connection to mysql: %v", err)
		}
		db.SetMaxIdleConns(cfg.MaxIdleConns)
		db.SetMaxOpenConns(cfg.MaxOpenConns)
		drv = entsql.OpenDB(dialect.MySQL, db)
	case "postgres", "postgresql":
		dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s %s",
			cfg.Host, cfg.Port, cfg.UserName, cfg.Password, cfg.Name, cfg.Config)
		db, err := sql.Open("postgres", dsn)
		if err != nil {
			return nil, fmt.Errorf("failed opening connection to postgres: %v", err)
		}
		db.SetMaxIdleConns(cfg.MaxIdleConns)
		db.SetMaxOpenConns(cfg.MaxOpenConns)
		drv = entsql.OpenDB(dialect.Postgres, db)
	case "sqlite3", "sqlite":
		db, err := sql.Open("sqlite3", cfg.Name+"?_fk=1")
		if err != nil {
			return nil, fmt.Errorf("failed opening connection to sqlite: %v", err)
		}
		drv = entsql.OpenDB(dialect.SQLite, db)
	default:
		return nil, fmt.Errorf("unsupported database type: %s", cfg.Type)
	}

	Client = ent.NewClient(ent.Driver(drv))

	// 自动迁移
	if err := Client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	return Client, nil
}

// CloseData 关闭数据库连接
func CloseData() {
	if Client != nil {
		_ = Client.Close()
	}
}
