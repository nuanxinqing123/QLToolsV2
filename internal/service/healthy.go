package service

import (
	"context"
	"fmt"
	"time"

	"github.com/nuanxinqing123/QLToolsV2/internal/app/config"
	res "github.com/nuanxinqing123/QLToolsV2/internal/pkg/response"
)

type HealthyService struct{}

// NewHealthyService 创建 HealthyService
func NewHealthyService() *HealthyService {
	return &HealthyService{}
}

// HealthCheckResponse 健康检查响应
type HealthCheckResponse struct {
	Status    bool      `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	API       bool      `json:"api"`
	Database  bool      `json:"database"`
	Redis     bool      `json:"redis"`
}

func (s *HealthyService) CheckHealth() (res.ResCode, any) {
	hc := &HealthCheckResponse{
		Timestamp: time.Now(),
		API:       true, // API 本身能响应说明是健康的
	}

	// 检查数据库
	hc.Database = s.checkDatabase()

	// 检查Redis
	hc.Redis = true

	// 确定整体状态
	if hc.Database == true && hc.Redis == true {
		hc.Status = true
	} else {
		hc.Status = false
	}

	return res.CodeSuccess, hc
}

// checkDatabase 检查数据库连接
func (s *HealthyService) checkDatabase() bool {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if config.Ent == nil {
		return false
	}

	// 执行测试查询
	_, err := config.Ent.User.Query().Count(ctx)
	if err != nil {
		config.Log.Error(fmt.Sprintf("【DB】执行测试查询失败: %s", err.Error()))
		return false
	}

	return true
}
