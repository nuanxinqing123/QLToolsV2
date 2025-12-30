package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/nuanxinqing123/QLToolsV2/internal/app/config"
	"github.com/nuanxinqing123/QLToolsV2/internal/data/ent"
	"github.com/nuanxinqing123/QLToolsV2/internal/data/ent/envplugin"
	"github.com/nuanxinqing123/QLToolsV2/internal/data/ent/plugin"
	"github.com/nuanxinqing123/QLToolsV2/internal/data/ent/pluginexecutionlog"
	pkgPlugin "github.com/nuanxinqing123/QLToolsV2/internal/pkg/plugin"
	"github.com/nuanxinqing123/QLToolsV2/internal/schema"
)

type PluginService struct {
	engine *pkgPlugin.Engine
}

// NewPluginService 创建 PluginService
func NewPluginService() *PluginService {
	return &PluginService{
		engine: pkgPlugin.NewEngine(5 * time.Second), // 默认5秒超时
	}
}

// CreatePlugin 创建插件
func (s *PluginService) CreatePlugin(req schema.CreatePluginRequest) (*schema.CreatePluginResponse, error) {
	ctx := context.Background()
	// 验证脚本语法
	if err := s.engine.ValidateScript(req.ScriptContent); err != nil {
		return nil, fmt.Errorf("脚本语法错误: %w", err)
	}

	// 检查插件名称是否已存在
	exists, err := config.Ent.Plugin.Query().
		Where(plugin.NameEQ(req.Name)).
		Exist(ctx)
	if err != nil {
		return nil, fmt.Errorf("查询插件失败: %w", err)
	}
	if exists {
		return nil, errors.New("插件名称已存在")
	}

	// 创建插件记录
	builder := config.Ent.Plugin.Create().
		SetName(req.Name).
		SetVersion(req.Version).
		SetScriptContent(req.ScriptContent).
		SetIsEnable(true). // 默认启用
		SetExecutionTimeout(int32(req.ExecutionTimeout)).
		SetPriority(int32(req.Priority)).
		SetTriggerEvent(req.TriggerEvent).
		SetCreatedAt(time.Now()).
		SetUpdatedAt(time.Now())

	if req.Description != "" {
		builder.SetDescription(req.Description)
	}
	if req.Author != "" {
		builder.SetAuthor(req.Author)
	}

	p, err := builder.Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("创建插件失败: %w", err)
	}

	return &schema.CreatePluginResponse{
		ID:      p.ID,
		Message: "插件创建成功",
	}, nil
}

// UpdatePlugin 更新插件
func (s *PluginService) UpdatePlugin(req schema.UpdatePluginRequest) (*schema.UpdatePluginResponse, error) {
	ctx := context.Background()
	// 验证脚本语法
	if err := s.engine.ValidateScript(req.ScriptContent); err != nil {
		return nil, fmt.Errorf("脚本语法错误: %w", err)
	}

	// 查询插件是否存在
	p, err := config.Ent.Plugin.Get(ctx, req.ID)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.New("插件不存在")
		}
		return nil, fmt.Errorf("查询插件失败: %w", err)
	}

	// 检查名称是否与其他插件冲突
	if req.Name != p.Name {
		exists, err := config.Ent.Plugin.Query().
			Where(plugin.And(
				plugin.NameEQ(req.Name),
				plugin.IDNEQ(req.ID),
			)).
			Exist(ctx)
		if err != nil {
			return nil, fmt.Errorf("查询插件失败: %w", err)
		}
		if exists {
			return nil, errors.New("插件名称已存在")
		}
	}

	// 验证触发事件类型
	if !isValidTriggerEvent(req.TriggerEvent) {
		return nil, errors.New("无效的触发事件类型")
	}

	// 执行更新
	builder := config.Ent.Plugin.UpdateOneID(req.ID).
		SetName(req.Name).
		SetDescription(req.Description).
		SetVersion(req.Version).
		SetAuthor(req.Author).
		SetScriptContent(req.ScriptContent).
		SetTriggerEvent(req.TriggerEvent).
		SetExecutionTimeout(int32(req.ExecutionTimeout)).
		SetPriority(int32(req.Priority)).
		SetUpdatedAt(time.Now())

	// 如果提供了启用状态，则更新
	if req.IsEnable != nil {
		builder.SetIsEnable(*req.IsEnable)
	}

	if err := builder.Exec(ctx); err != nil {
		return nil, fmt.Errorf("更新插件失败: %w", err)
	}

	return &schema.UpdatePluginResponse{
		Message: "插件更新成功",
	}, nil
}

// GetPlugin 获取单个插件信息
func (s *PluginService) GetPlugin(id int64) (*schema.GetPluginResponse, error) {
	ctx := context.Background()
	p, err := config.Ent.Plugin.Get(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.New("插件不存在")
		}
		return nil, fmt.Errorf("查询插件失败: %w", err)
	}

	return &schema.GetPluginResponse{
		ID:               p.ID,
		Name:             p.Name,
		Description:      p.Description,
		Version:          p.Version,
		Author:           p.Author,
		ScriptContent:    p.ScriptContent,
		IsEnable:         p.IsEnable,
		ExecutionTimeout: int(p.ExecutionTimeout),
		CreatedAt:        p.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:        p.UpdatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

// GetPluginList 获取插件列表
func (s *PluginService) GetPluginList(req schema.GetPluginListRequest) (*schema.GetPluginListResponse, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	ctx := context.Background()
	query := config.Ent.Plugin.Query()

	if req.Name != "" {
		query.Where(plugin.NameContains(req.Name))
	}
	if req.IsEnable != nil {
		query.Where(plugin.IsEnableEQ(*req.IsEnable))
	}

	total, err := query.Count(ctx)
	if err != nil {
		return nil, fmt.Errorf("查询插件总数失败: %w", err)
	}

	offset := (req.Page - 1) * req.PageSize
	plugins, err := query.Offset(offset).
		Limit(req.PageSize).
		Order(ent.Desc(plugin.FieldCreatedAt)).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("查询插件列表失败: %w", err)
	}

	list := make([]schema.GetPluginResponse, 0, len(plugins))
	for _, p := range plugins {
		list = append(list, schema.GetPluginResponse{
			ID:               p.ID,
			Name:             p.Name,
			Description:      p.Description,
			Version:          p.Version,
			Author:           p.Author,
			ScriptContent:    p.ScriptContent,
			IsEnable:         p.IsEnable,
			ExecutionTimeout: int(p.ExecutionTimeout),
			TriggerEvent:     p.TriggerEvent,
			CreatedAt:        p.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:        p.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return &schema.GetPluginListResponse{
		Total: int64(total),
		List:  list,
	}, nil
}

// DeletePlugin 删除插件
func (s *PluginService) DeletePlugin(req schema.DeletePluginRequest) (*schema.DeletePluginResponse, error) {
	ctx := context.Background()
	// 检查插件是否存在
	_, err := config.Ent.Plugin.Get(ctx, req.ID)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.New("插件不存在")
		}
		return nil, fmt.Errorf("查询插件失败: %w", err)
	}

	// 开启事务
	tx, err := config.Ent.Tx(ctx)
	if err != nil {
		return nil, fmt.Errorf("开启事务失败: %w", err)
	}

	// 确保在函数结束时正确处理事务
	defer func() {
		if v := recover(); v != nil {
			_ = tx.Rollback()
			panic(v)
		}
	}()

	// 删除插件与环境变量的关联关系
	_, err = tx.EnvPlugin.Delete().
		Where(envplugin.PluginIDEQ(req.ID)).
		Exec(ctx)
	if err != nil {
		_ = tx.Rollback()
		return nil, fmt.Errorf("删除插件环境变量关联失败: %w", err)
	}

	// 执行删除插件
	err = tx.Plugin.DeleteOneID(req.ID).Exec(ctx)
	if err != nil {
		_ = tx.Rollback()
		return nil, fmt.Errorf("删除插件失败: %w", err)
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("提交事务失败: %w", err)
	}

	return &schema.DeletePluginResponse{
		Message: "插件删除成功",
	}, nil
}

// TogglePluginStatus 切换插件启用状态
func (s *PluginService) TogglePluginStatus(req schema.TogglePluginStatusRequest) (*schema.TogglePluginStatusResponse, error) {
	ctx := context.Background()
	_, err := config.Ent.Plugin.Get(ctx, req.ID)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.New("插件不存在")
		}
		return nil, fmt.Errorf("查询插件失败: %w", err)
	}

	if err := config.Ent.Plugin.UpdateOneID(req.ID).
		SetIsEnable(req.IsEnable).
		SetUpdatedAt(time.Now()).
		Exec(ctx); err != nil {
		return nil, fmt.Errorf("更新插件状态失败: %w", err)
	}

	status := "禁用"
	if req.IsEnable {
		status = "启用"
	}

	return &schema.TogglePluginStatusResponse{
		Message: fmt.Sprintf("插件已%s", status),
	}, nil
}

// TestPlugin 测试插件
func (s *PluginService) TestPlugin(req schema.TestPluginRequest) (*schema.TestPluginResponse, error) {
	// 执行测试
	result := s.engine.TestScript(req.ScriptContent, req.TestEnvValue)

	// 将输出数据转换为string类型
	var outputDataStr string
	if len(result.OutputData) > 0 {
		outputDataStr = string(result.OutputData)
	}

	return &schema.TestPluginResponse{
		Success:       result.Success,
		ExecutionTime: result.ExecutionTime,
		OutputData:    outputDataStr,
		ErrorMessage:  result.ErrorMessage,
	}, nil
}

// ExecutePluginsForEnv 为指定环境变量执行插件
func (s *PluginService) ExecutePluginsForEnv(envID int64, envValue string) (*pkgPlugin.ExecutionResult, error) {
	ctx := context.Background()
	// 查询该环境变量启用的插件
	results, err := config.Ent.EnvPlugin.Query().
		Where(
			envplugin.EnvIDEQ(envID),
			envplugin.IsEnableEQ(true),
		).
		WithPlugin(func(q *ent.PluginQuery) {
			q.Where(plugin.IsEnableEQ(true))
		}).
		Order(ent.Asc(envplugin.FieldExecutionOrder)).
		All(ctx)

	if err != nil {
		return nil, fmt.Errorf("查询环境变量插件失败: %w", err)
	}

	// 如果没有插件，返回原始值
	if len(results) == 0 {
		// 构建返回数据
		resultData := map[string]interface{}{
			"bool": true,
			"env":  envValue,
		}
		outputBytes, _ := config.JSON.Marshal(resultData)
		return &pkgPlugin.ExecutionResult{
			Success:    true,
			OutputData: outputBytes,
		}, nil
	}

	// 依次执行插件
	var lastResult *pkgPlugin.ExecutionResult
	for _, item := range results {
		p := item.Edges.Plugin
		if p == nil {
			continue
		}

		// 构建执行上下文
		var configData []byte
		if item.Config != nil {
			configData = []byte(*item.Config)
		} else {
			configData = []byte("{}")
		}

		execCtx := &pkgPlugin.ExecutionContext{
			PluginID:  item.PluginID,
			EnvID:     envID,
			EnvValue:  envValue,
			Config:    configData,
			Timestamp: time.Now().Unix(),
		}

		// 执行插件
		timeout := time.Duration(p.ExecutionTimeout) * time.Millisecond
		result := s.engine.Execute(context.Background(), p.ScriptContent, execCtx, timeout)

		// 记录执行日志
		s.logPluginExecution(item.PluginID, envID, result)

		lastResult = result

		// 如果执行失败，返回失败结果
		if !result.Success {
			return result, nil
		}

		// 如果插件返回了新的环境变量值，更新envValue用于下一个插件
		if len(result.OutputData) > 0 {
			var output map[string]interface{}
			if err := config.JSON.Unmarshal(result.OutputData, &output); err == nil {
				if newEnv, ok := output["env"].(string); ok {
					envValue = newEnv
				}
			}
		}
	}

	return lastResult, nil
}

// logPluginExecution 记录插件执行日志
func (s *PluginService) logPluginExecution(pluginID, envID int64, result *pkgPlugin.ExecutionResult) {
	status := "success"
	if !result.Success {
		status = "error"
	}

	// 处理可选字段
	var outputDataStr string
	if len(result.OutputData) > 0 {
		outputDataStr = string(result.OutputData)
	}

	// 异步记录日志，不影响主流程
	go func() {
		ctx := context.Background()
		_, err := config.Ent.PluginExecutionLog.Create().
			SetCreatedAt(time.Now()).
			SetPluginID(pluginID).
			SetEnvID(envID).
			SetExecutionStatus(status).
			SetExecutionTime(int32(result.ExecutionTime)).
			SetOutputData(outputDataStr).
			SetErrorMessage(result.ErrorMessage).
			SetStackTrace(result.StackTrace).
			Save(ctx)
		if err != nil {
			config.Log.Warn(err.Error()) // 仅做记录
		}
	}()
}

// BindPluginToEnv 绑定插件到环境变量
func (s *PluginService) BindPluginToEnv(req schema.BindPluginToEnvRequest) (*schema.BindPluginToEnvResponse, error) {
	ctx := context.Background()
	// 检查插件是否存在
	_, err := config.Ent.Plugin.Get(ctx, req.PluginID)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.New("插件不存在")
		}
		return nil, fmt.Errorf("查询插件失败: %w", err)
	}

	// 检查环境变量是否存在
	_, err = config.Ent.Env.Get(ctx, req.EnvID)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.New("环境变量不存在")
		}
		return nil, fmt.Errorf("查询环境变量失败: %w", err)
	}

	// 检查是否已经绑定
	existingBinding, err := config.Ent.EnvPlugin.Query().
		Where(
			envplugin.PluginIDEQ(req.PluginID),
			envplugin.EnvIDEQ(req.EnvID),
		).
		Only(ctx)
	if err != nil && !ent.IsNotFound(err) {
		return nil, fmt.Errorf("查询绑定关系失败: %w", err)
	}

	// 设置默认执行顺序
	executionOrder := req.ExecutionOrder
	if executionOrder == 0 {
		executionOrder = 100
	}

	if existingBinding != nil {
		// 如果已经绑定，更新配置
		err = config.Ent.EnvPlugin.UpdateOne(existingBinding).
			SetConfig(req.Config).
			SetExecutionOrder(executionOrder).
			SetIsEnable(true).
			SetUpdatedAt(time.Now()).
			Exec(ctx)
		if err != nil {
			return nil, fmt.Errorf("更新绑定关系失败: %w", err)
		}
	} else {
		// 创建新的绑定关系
		err = config.Ent.EnvPlugin.Create().
			SetPluginID(req.PluginID).
			SetEnvID(req.EnvID).
			SetIsEnable(true).
			SetExecutionOrder(executionOrder).
			SetConfig(req.Config).
			SetCreatedAt(time.Now()).
			SetUpdatedAt(time.Now()).
			Exec(ctx)
		if err != nil {
			return nil, fmt.Errorf("创建绑定关系失败: %w", err)
		}
	}

	return &schema.BindPluginToEnvResponse{
		Message: "插件绑定成功",
	}, nil
}

// UnbindPluginFromEnv 解绑插件与环境变量
func (s *PluginService) UnbindPluginFromEnv(req schema.UnbindPluginFromEnvRequest) (*schema.UnbindPluginFromEnvResponse, error) {
	ctx := context.Background()
	// 检查绑定关系是否存在
	_, err := config.Ent.EnvPlugin.Query().
		Where(
			envplugin.PluginIDEQ(req.PluginID),
			envplugin.EnvIDEQ(req.EnvID),
		).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.New("绑定关系不存在")
		}
		return nil, fmt.Errorf("查询绑定关系失败: %w", err)
	}

	// 删除绑定关系
	_, err = config.Ent.EnvPlugin.Delete().
		Where(
			envplugin.PluginIDEQ(req.PluginID),
			envplugin.EnvIDEQ(req.EnvID),
		).
		Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("删除绑定关系失败: %w", err)
	}

	return &schema.UnbindPluginFromEnvResponse{
		Message: "插件解绑成功",
	}, nil
}

// GetPluginEnvs 获取插件关联环境变量
func (s *PluginService) GetPluginEnvs(req schema.GetPluginEnvsRequest) (*schema.GetPluginEnvsResponse, error) {
	ctx := context.Background()
	// 检查插件是否存在
	_, err := config.Ent.Plugin.Get(ctx, req.PluginID)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.New("插件不存在")
		}
		return nil, fmt.Errorf("查询插件失败: %w", err)
	}

	// 查询插件关联的环境变量
	results, err := config.Ent.EnvPlugin.Query().
		Where(envplugin.PluginIDEQ(req.PluginID)).
		WithEnv().
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("查询插件环境变量关联失败: %w", err)
	}

	// 转换为响应格式
	envs := make([]schema.PluginEnvRelationInfo, 0, len(results))
	for _, res := range results {
		envName := ""
		if res.Edges.Env != nil {
			envName = res.Edges.Env.Name
		}
		envs = append(envs, schema.PluginEnvRelationInfo{
			EnvID:          res.EnvID,
			EnvName:        envName,
			IsEnable:       res.IsEnable,
			ExecutionOrder: res.ExecutionOrder,
			Config:         res.Config,
			CreatedAt:      res.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return &schema.GetPluginEnvsResponse{
		PluginID: req.PluginID,
		Envs:     envs,
	}, nil
}

// GetPluginExecutionLogs 获取插件执行日志
func (s *PluginService) GetPluginExecutionLogs(req schema.GetPluginExecutionLogsRequest) (*schema.GetPluginExecutionLogsResponse, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	ctx := context.Background()
	query := config.Ent.PluginExecutionLog.Query()

	// 按插件ID筛选
	if req.PluginID != nil {
		query.Where(pluginexecutionlog.PluginIDEQ(*req.PluginID))
	}

	// 按环境变量ID筛选
	if req.EnvID != nil {
		query.Where(pluginexecutionlog.EnvIDEQ(*req.EnvID))
	}

	// 按执行状态筛选
	if req.ExecutionStatus != "" {
		query.Where(pluginexecutionlog.ExecutionStatusEQ(req.ExecutionStatus))
	}

	// 按时间范围筛选
	if req.StartTime != "" {
		if startTime, err := time.Parse("2006-01-02 15:04:05", req.StartTime); err == nil {
			query.Where(pluginexecutionlog.CreatedAtGTE(startTime))
		}
	}
	if req.EndTime != "" {
		if endTime, err := time.Parse("2006-01-02 15:04:05", req.EndTime); err == nil {
			query.Where(pluginexecutionlog.CreatedAtLTE(endTime))
		}
	}

	total, err := query.Count(ctx)
	if err != nil {
		return nil, fmt.Errorf("查询日志总数失败: %w", err)
	}

	offset := (req.Page - 1) * req.PageSize
	logs, err := query.Offset(offset).
		Limit(req.PageSize).
		Order(ent.Desc(pluginexecutionlog.FieldCreatedAt)).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("查询执行日志失败: %w", err)
	}

	list := make([]schema.PluginExecutionLogInfo, 0, len(logs))
	for _, l := range logs {
		list = append(list, schema.PluginExecutionLogInfo{
			ID:              l.ID,
			PluginID:        l.PluginID,
			EnvID:           l.EnvID,
			ExecutionStatus: l.ExecutionStatus,
			ExecutionTime:   int(l.ExecutionTime),
			OutputData:      l.OutputData,
			ErrorMessage:    l.ErrorMessage,
			CreatedAt:       l.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return &schema.GetPluginExecutionLogsResponse{
		Total: int64(total),
		List:  list,
	}, nil
}

// isValidTriggerEvent 验证触发事件类型
func isValidTriggerEvent(event string) bool {
	validEvents := []string{"before_submit", "after_submit", "on_error"}
	for _, validEvent := range validEvents {
		if event == validEvent {
			return true
		}
	}
	return false
}
