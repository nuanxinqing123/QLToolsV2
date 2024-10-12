package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"QLToolsV2/config"
	"QLToolsV2/initialize"
	"QLToolsV2/internal/cron"
	"QLToolsV2/utils"
	"QLToolsV2/utils/validator"
)

func main() {
	// 初始化配置
	config.GinVP = initialize.Viper()

	// 初始化日志
	config.GinLOG = initialize.Zap()

	// 初始化数据库
	config.GinDB = initialize.Gorm() // gorm连接数据库
	if config.GinDB != nil {
		// 初始化表
		initialize.RegisterTables(config.GinDB)
		fmt.Println("数据库初始化成功")
	} else {
		fmt.Println("数据库启动失败...")
		return
	}

	// 初始化 GCache
	config.GinCache = initialize.InitGCache()
	if config.GinCache == nil {
		fmt.Println("GCache初始化成功")
	}

	// 启动定时任务
	if err := cron.InitTask(); err != nil {
		fmt.Printf("定时任务初始化失败, err:%v\n", err)
		return
	}
	zap.L().Debug("定时任务初始化成功...")

	// 初始化雪花 ID 算法
	if err := utils.InitSnowflake(); err != nil {
		fmt.Println("初始化雪花 ID 算法失败...")
		return
	}

	// 初始化翻译器
	if err := validator.InitTrans("zh"); err != nil {
		fmt.Printf("翻译器初始化失败, err:%v\n", err)
		return
	}

	// 初始化路由
	router := initialize.Routers()
	if router == nil {
		fmt.Println("初始化路由失败...")
		return
	}

	// 启动服务
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.GinConfig.App.Port),
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
	config.GinLOG.Info("Service ready to shut down")

	// 创建10秒超时的Context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// 10秒内优雅关闭服务（将未处理完成的请求处理完再关闭服务），超过10秒就超时退出
	if err := srv.Shutdown(ctx); err != nil {
		config.GinLOG.Fatal("Service timed out has been shut down: ", zap.Error(err))
	}

	config.GinLOG.Info("Service has been shut down")
}
