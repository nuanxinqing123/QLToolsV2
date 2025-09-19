package service

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
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

func (s *HealthyService) CheckHealth(ctx *gin.Context) (res.ResCode, any) {
	hc := &HealthCheckResponse{
		Timestamp: time.Now(),
		API:       true, // API 本身能响应说明是健康的
	}

	// 检查数据库
	hc.Database = s.checkDatabase()

	// 检查Redis
	hc.Redis = s.checkRedis(ctx)

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
	if config.DB == nil {
		return false
	}

	var result int
	if err := config.DB.Raw("SELECT 1").Scan(&result).Error; err != nil {
		config.Log.Error(fmt.Sprintf("【DB】执行测试查询失败: %s", err.Error()))
		return false
	}

	return true
}

// checkRedis 检查Redis连接
func (s *HealthyService) checkRedis(ctx *gin.Context) bool {
	if config.Cache == nil {
		return false
	}

	// 执行 ping 命令
	if err := config.Cache.Ping(ctx).Err(); err != nil {
		config.Log.Error(fmt.Sprintf("【Cache】执行测试查询失败: %s", err.Error()))
		return false
	}

	return true
}
