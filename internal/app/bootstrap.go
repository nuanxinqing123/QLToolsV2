package app

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"github.com/nuanxinqing123/QLToolsV2/internal/app/config"
	"github.com/nuanxinqing123/QLToolsV2/internal/app/initializer"
	"github.com/nuanxinqing123/QLToolsV2/internal/repository"
	"go.uber.org/zap"
)

// Start 启动服务
func Start() {
	// 解析命令行参数，支持通过 -config 或 -c 指定配置文件路径
	// 示例：./app -config=/path/to/config.yaml 或 ./app -c ./configs/dev.yaml
	var configPath string
	flag.StringVar(&configPath, "config", "", "配置文件路径（默认：configs/config.yaml）")
	flag.StringVar(&configPath, "c", "", "配置文件路径（默认：configs/config.yaml）")
	flag.Parse()

	// 初始化配置文件
	config.VP = initializer.Viper(configPath)
	if config.VP != nil {
		fmt.Println("Viper初始化成功")
	}

	// 初始化日志
	config.Log = initializer.Zap()

	// 初始化数据库
	config.DB = initializer.Gorm()
	if config.DB == nil {
		config.Log.Fatal("数据库连接失败")
		return
	} else {
		// 初始化表
		initializer.RegisterTables(config.DB)
		// 设置 Gorm Gen 使用的默认数据库
		repository.SetDefault(config.DB)
		config.Log.Info("数据库连接成功")
	}

	// 初始化缓存
	config.Cache = initializer.Cache()
	if config.Cache == nil {
		config.Log.Fatal("Redis连接失败")
		return
	} else {
		config.Log.Info("Redis连接成功")
	}

	// 初始化JSON编解码器
	config.JSON = jsoniter.ConfigCompatibleWithStandardLibrary

	// 启动限速器清理任务
	initializer.StartRateLimitCleanup()

	router := initializer.Routers()

	fmt.Println(" ")
	switch config.Config.App.Mode {
	case gin.DebugMode:
		fmt.Println("运行模式: Debug模式")
		gin.SetMode(gin.DebugMode)
	case gin.TestMode:
		fmt.Println("运行模式: Test模式")
		gin.SetMode(gin.TestMode)
	default:
		fmt.Println("运行模式: Release模式")
		gin.SetMode(gin.ReleaseMode)
	}
	fmt.Println("监听端口: " + strconv.Itoa(config.Config.App.Port))
	fmt.Println(" ")

	// 启动服务
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.Config.App.Port),
		Handler: router,
	}

	// 启动
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Listten: %s\n", err)
		}
	}()

	// 等待终端信号来优雅关闭服务器，为关闭服务器设置10秒超时
	quit := make(chan os.Signal, 1) // 创建一个接受信号的通道

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // 此处不会阻塞
	<-quit                                               // 阻塞此处，当接受到上述两种信号时，才继续往下执行
	config.Log.Info("Service ready to shut down")

	// 创建10秒超时的Context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// 10秒内优雅关闭服务（将未处理完成的请求处理完再关闭服务），超过10秒就超时退出
	if err := srv.Shutdown(ctx); err != nil {
		config.Log.Fatal("Service timed out has been shut down: ", zap.Error(err))
	}

	config.Log.Info("Service has been shut down")
}
