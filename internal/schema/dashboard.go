package schema

// OverviewResponse 数据总览响应
type OverviewResponse struct {
	OnlineServices int64 `json:"online_services"` // 在线服务数量（启用的变量数量）
	TotalPanels    int64 `json:"total_panels"`    // 总面板数
	ActiveCDK      int64 `json:"active_cdk"`      // 活跃CDK数量
	TodaySubmit    int64 `json:"today_submit"`    // 今日提交数量
}

// SubmitTrendItem 提交趋势项
type SubmitTrendItem struct {
	Date  string `json:"date"`  // 日期
	Count int64  `json:"count"` // 提交数量
}

// SubmitTrendResponse 提交趋势响应
type SubmitTrendResponse struct {
	Trend []SubmitTrendItem `json:"trend"` // 趋势数据
}

// ActivityItem 活动项
type ActivityItem struct {
	Time        string `json:"time"`        // 时间
	Type        string `json:"type"`        // 类型：submit/login/error
	Description string `json:"description"` // 描述
	Status      string `json:"status"`      // 状态：success/warning/error
}

// RecentActivityResponse 最近活动响应
type RecentActivityResponse struct {
	Activities []ActivityItem `json:"activities"` // 活动列表
}

// ResourceUsageResponse 资源使用响应
type ResourceUsageResponse struct {
	CPU    ResourceItem `json:"cpu"`    // CPU使用情况
	Memory ResourceItem `json:"memory"` // 内存使用情况
	Disk   ResourceItem `json:"disk"`   // 磁盘使用情况
}

// ResourceItem 资源项
type ResourceItem struct {
	Percentage float64 `json:"percentage"` // 使用百分比
	Used       string  `json:"used"`       // 已使用（格式化字符串）
	Total      string  `json:"total"`      // 总量（格式化字符串）
}
