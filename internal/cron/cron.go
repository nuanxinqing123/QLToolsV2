package cron

import (
	"time"

	"github.com/robfig/cron/v3"

	"QLToolsV2/internal/service"
)

var err error
var c *cron.Cron

// InitTask 定时任务
func InitTask() error {
	// 刷新并启用启动任务
	c = cron.New(cron.WithLocation(time.FixedZone("CST", 8*3600))) // 设置时区

	// 定时任务区

	// 定时更新面板Token（0 0 1/1 * *）
	_, err = c.AddFunc("0 0 1/1 * *", func() {
		service.RefreshPanel()
	})

	// 定时任务结束

	if err != nil {
		return err
	}
	c.Start()
	return nil
}

// StopTask 暂停任务
func StopTask() {
	c.Stop()
}
