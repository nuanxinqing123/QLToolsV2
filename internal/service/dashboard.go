package service

import (
	"context"
	"fmt"
	"time"

	"github.com/nuanxinqing123/QLToolsV2/internal/app/config"
	"github.com/nuanxinqing123/QLToolsV2/internal/data/ent/cdkey"
	"github.com/nuanxinqing123/QLToolsV2/internal/data/ent/env"
	"github.com/nuanxinqing123/QLToolsV2/internal/schema"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
)

type DashboardService struct{}

// NewDashboardService 创建DashboardService实例
func NewDashboardService() *DashboardService {
	return &DashboardService{}
}

// GetOverview 获取数据总览
func (s *DashboardService) GetOverview() (*schema.OverviewResponse, error) {
	var resp schema.OverviewResponse
	ctx := context.Background()

	// 1. 获取在线服务数量（启用的变量数量）
	onlineCount, err := config.Ent.Env.Query().
		Where(env.IsEnableEQ(true)).
		Count(ctx)
	if err != nil {
		return nil, fmt.Errorf("获取在线服务数量失败: %w", err)
	}
	resp.OnlineServices = int64(onlineCount)

	// 2. 获取总面板数
	panelCount, err := config.Ent.Panel.Query().Count(ctx)
	if err != nil {
		return nil, fmt.Errorf("获取总面板数失败: %w", err)
	}
	resp.TotalPanels = int64(panelCount)

	// 3. 获取活跃CDK数量（is_enable=true 且 count>0）
	cdkCount, err := config.Ent.CdKey.Query().
		Where(
			cdkey.IsEnableEQ(true),
			cdkey.CountGT(0),
		).
		Count(ctx)
	if err != nil {
		return nil, fmt.Errorf("获取活跃CDK数量失败: %w", err)
	}
	resp.ActiveCDK = int64(cdkCount)

	// 4. 今日提交数量（暂时返回0，等待后续实现提交记录表）
	resp.TodaySubmit = 0

	return &resp, nil
}

// GetSubmitTrend 获取提交趋势（模拟数据）
func (s *DashboardService) GetSubmitTrend() (*schema.SubmitTrendResponse, error) {
	// 生成最近7天的模拟数据
	var trend []schema.SubmitTrendItem
	now := time.Now()

	for i := 6; i >= 0; i-- {
		date := now.AddDate(0, 0, -i)
		trend = append(trend, schema.SubmitTrendItem{
			Date:  date.Format("01-02"),
			Count: int64(50 + i*10), // 模拟递增趋势
		})
	}

	return &schema.SubmitTrendResponse{
		Trend: trend,
	}, nil
}

// GetRecentActivity 获取最近活动（模拟数据）
func (s *DashboardService) GetRecentActivity() (*schema.RecentActivityResponse, error) {
	now := time.Now()

	// 生成模拟活动数据
	activities := []schema.ActivityItem{
		{
			Time:        now.Add(-5 * time.Minute).Format("15:04:05"),
			Type:        "submit",
			Description: "用户提交了环境变量 JD_COOKIE",
			Status:      "success",
		},
		{
			Time:        now.Add(-15 * time.Minute).Format("15:04:05"),
			Type:        "login",
			Description: "管理员登录系统",
			Status:      "success",
		},
		{
			Time:        now.Add(-30 * time.Minute).Format("15:04:05"),
			Type:        "submit",
			Description: "用户提交了环境变量 PT_KEY",
			Status:      "success",
		},
		{
			Time:        now.Add(-45 * time.Minute).Format("15:04:05"),
			Type:        "error",
			Description: "面板连接失败: 192.168.1.100",
			Status:      "error",
		},
		{
			Time:        now.Add(-60 * time.Minute).Format("15:04:05"),
			Type:        "submit",
			Description: "用户提交了环境变量 ELE_COOKIE",
			Status:      "success",
		},
	}

	return &schema.RecentActivityResponse{
		Activities: activities,
	}, nil
}

// GetResourceUsage 获取资源使用情况
func (s *DashboardService) GetResourceUsage() (*schema.ResourceUsageResponse, error) {
	var resp schema.ResourceUsageResponse

	// 1. 获取CPU使用率
	cpuPercent, err := cpu.Percent(time.Second, false)
	if err != nil {
		return nil, fmt.Errorf("获取CPU使用率失败: %w", err)
	}
	if len(cpuPercent) > 0 {
		resp.CPU.Percentage = cpuPercent[0]
		resp.CPU.Used = fmt.Sprintf("%.2f%%", cpuPercent[0])
		resp.CPU.Total = "100%"
	}

	// 2. 获取内存使用情况
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		return nil, fmt.Errorf("获取内存使用情况失败: %w", err)
	}
	resp.Memory.Percentage = memInfo.UsedPercent
	resp.Memory.Used = formatBytes(memInfo.Used)
	resp.Memory.Total = formatBytes(memInfo.Total)

	// 3. 获取磁盘使用情况
	diskInfo, err := disk.Usage("/")
	if err != nil {
		return nil, fmt.Errorf("获取磁盘使用情况失败: %w", err)
	}
	resp.Disk.Percentage = diskInfo.UsedPercent
	resp.Disk.Used = formatBytes(diskInfo.Used)
	resp.Disk.Total = formatBytes(diskInfo.Total)

	return &resp, nil
}

// formatBytes 格式化字节数为可读字符串
func formatBytes(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := uint64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.2f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
