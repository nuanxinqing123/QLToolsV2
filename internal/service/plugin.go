package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/nuanxinqing123/QLToolsV2/internal/app/config"
	"github.com/nuanxinqing123/QLToolsV2/internal/model"
	"github.com/nuanxinqing123/QLToolsV2/internal/pkg/plugin"
	"github.com/nuanxinqing123/QLToolsV2/internal/repository"
	"github.com/nuanxinqing123/QLToolsV2/internal/schema"
	"gorm.io/gorm"
)

type PluginService struct {
	engine *plugin.Engine
}

// NewPluginService 创建 PluginService
func NewPluginService() *PluginService {
	return &PluginService{
		engine: plugin.NewEngine(5 * time.Second), // 默认5秒超时
	}
}

// CreatePlugin 创建插件
func (s *PluginService) CreatePlugin(req schema.CreatePluginRequest) (*schema.CreatePluginResponse, error) {
	// 验证脚本语法
	if err := s.engine.ValidateScript(req.ScriptContent); err != nil {
		return nil, fmt.Errorf("脚本语法错误: %w", err)
	}

	// 检查插件名称是否已存在
	existingPlugin, err := repository.Plugins.Where(
		repository.Plugins.Name.Eq(req.Name),
	).Take()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("查询插件失败: %w", err)
	}
	if existingPlugin != nil {
		return nil, errors.New("插件名称已存在")
	}

	now := time.Now()

	// 处理可选字段
	var description *string
	if req.Description != "" {
		description = &req.Description
	}

	var author *string
	if req.Author != "" {
		author = &req.Author
	}

	pluginModel := &model.Plugins{
		CreatedAt:        now,
		UpdatedAt:        now,
		Name:             req.Name,
		Description:      description,
		Version:          req.Version,
		Author:           author,
		ScriptContent:    req.ScriptContent,
		IsEnable:         true, // 默认启用
		ExecutionTimeout: int32(req.ExecutionTimeout),
		Priority:         int32(req.Priority),
		TriggerEvent:     req.TriggerEvent,
	}

	// 创建插件记录
	if err = repository.Plugins.WithContext(context.Background()).Create(pluginModel); err != nil {
		return nil, fmt.Errorf("创建插件失败: %w", err)
	}

	return &schema.CreatePluginResponse{
		ID:      pluginModel.ID,
		Message: "插件创建成功",
	}, nil
}

// UpdatePlugin 更新插件
func (s *PluginService) UpdatePlugin(req schema.UpdatePluginRequest) (*schema.UpdatePluginResponse, error) {
	// 验证脚本语法
	if err := s.engine.ValidateScript(req.ScriptContent); err != nil {
		return nil, fmt.Errorf("脚本语法错误: %w", err)
	}

	// 查询插件是否存在
	pluginModel, err := repository.Plugins.Where(
		repository.Plugins.ID.Eq(req.ID),
	).Take()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("插件不存在")
		}
		return nil, fmt.Errorf("查询插件失败: %w", err)
	}

	// 检查名称是否与其他插件冲突
	if req.Name != pluginModel.Name {
		existingPlugin, err := repository.Plugins.Where(
			repository.Plugins.Name.Eq(req.Name),
			repository.Plugins.ID.Neq(req.ID),
		).Take()
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("查询插件失败: %w", err)
		}
		if existingPlugin != nil {
			return nil, errors.New("插件名称已存在")
		}
	}

	// 验证触发事件类型
	if !isValidTriggerEvent(req.TriggerEvent) {
		return nil, errors.New("无效的触发事件类型")
	}

	// 构建更新数据
	updates := map[string]interface{}{
		"name":              req.Name,
		"description":       req.Description,
		"version":           req.Version,
		"author":            req.Author,
		"script_content":    req.ScriptContent,
		"trigger_event":     req.TriggerEvent,
		"execution_timeout": int32(req.ExecutionTimeout),
		"priority":          int32(req.Priority),
		"updated_at":        time.Now(),
	}

	// 如果提供了启用状态，则更新
	if req.IsEnable != nil {
		updates["is_enable"] = *req.IsEnable
	}

	// 执行更新
	_, err = repository.Plugins.Where(
		repository.Plugins.ID.Eq(req.ID),
	).Updates(updates)
	if err != nil {
		return nil, fmt.Errorf("更新插件失败: %w", err)
	}

	return &schema.UpdatePluginResponse{
		Message: "插件更新成功",
	}, nil
}

// GetPlugin 获取单个插件信息
func (s *PluginService) GetPlugin(id int64) (*schema.GetPluginResponse, error) {
	pluginModel, err := repository.Plugins.Where(
		repository.Plugins.ID.Eq(id),
	).Take()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("插件不存在")
		}
		return nil, fmt.Errorf("查询插件失败: %w", err)
	}

	// 处理可选字段
	var description, author string
	if pluginModel.Description != nil {
		description = *pluginModel.Description
	}
	if pluginModel.Author != nil {
		author = *pluginModel.Author
	}

	return &schema.GetPluginResponse{
		ID:               pluginModel.ID,
		Name:             pluginModel.Name,
		Description:      description,
		Version:          pluginModel.Version,
		Author:           author,
		ScriptContent:    pluginModel.ScriptContent,
		IsEnable:         pluginModel.IsEnable,
		ExecutionTimeout: int(pluginModel.ExecutionTimeout),
		CreatedAt:        pluginModel.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:        pluginModel.UpdatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

// GetPluginList 获取插件列表
func (s *PluginService) GetPluginList(req schema.GetPluginListRequest) (*schema.GetPluginListResponse, error) {
	// 设置默认分页参数
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	// 构建查询条件
	query := repository.Plugins.WithContext(context.Background())

	// 按名称模糊搜索
	if req.Name != "" {
		query = query.Where(repository.Plugins.Name.Like("%" + req.Name + "%"))
	}

	// 按启用状态筛选
	if req.IsEnable != nil {
		query = query.Where(repository.Plugins.IsEnable.Is(*req.IsEnable))
	}

	// 获取总数
	total, err := query.Count()
	if err != nil {
		return nil, fmt.Errorf("查询插件总数失败: %w", err)
	}

	// 分页查询
	offset := (req.Page - 1) * req.PageSize
	plugins, err := query.Offset(offset).Limit(req.PageSize).Order(repository.Plugins.CreatedAt.Desc()).Find()
	if err != nil {
		return nil, fmt.Errorf("查询插件列表失败: %w", err)
	}

	// 转换为响应格式
	list := make([]schema.GetPluginResponse, 0, len(plugins))
	for _, pluginModel := range plugins {
		// 处理可选字段
		var description, author string
		if pluginModel.Description != nil {
			description = *pluginModel.Description
		}
		if pluginModel.Author != nil {
			author = *pluginModel.Author
		}

		list = append(list, schema.GetPluginResponse{
			ID:               pluginModel.ID,
			Name:             pluginModel.Name,
			Description:      description,
			Version:          pluginModel.Version,
			Author:           author,
			ScriptContent:    pluginModel.ScriptContent,
			IsEnable:         pluginModel.IsEnable,
			ExecutionTimeout: int(pluginModel.ExecutionTimeout),
			TriggerEvent:     pluginModel.TriggerEvent,
			CreatedAt:        pluginModel.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:        pluginModel.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return &schema.GetPluginListResponse{
		Total: total,
		List:  list,
	}, nil
}

// DeletePlugin 删除插件
func (s *PluginService) DeletePlugin(req schema.DeletePluginRequest) (*schema.DeletePluginResponse, error) {
	// 检查插件是否存在
	_, err := repository.Plugins.Where(
		repository.Plugins.ID.Eq(req.ID),
	).Take()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("插件不存在")
		}
		return nil, fmt.Errorf("查询插件失败: %w", err)
	}

	// 开启事务
	err = repository.Q.Transaction(func(tx *repository.Query) error {
		// 删除插件与环境变量的关联关系
		_, err := tx.EnvPlugins.Where(
			tx.EnvPlugins.PluginID.Eq(req.ID),
		).Delete()
		if err != nil {
			return fmt.Errorf("删除插件环境变量关联失败: %w", err)
		}

		// 执行软删除插件
		_, err = tx.Plugins.Where(
			tx.Plugins.ID.Eq(req.ID),
		).Delete()
		if err != nil {
			return fmt.Errorf("删除插件失败: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &schema.DeletePluginResponse{
		Message: "插件删除成功",
	}, nil
}

// TogglePluginStatus 切换插件启用状态
func (s *PluginService) TogglePluginStatus(req schema.TogglePluginStatusRequest) (*schema.TogglePluginStatusResponse, error) {
	// 检查插件是否存在
	_, err := repository.Plugins.Where(
		repository.Plugins.ID.Eq(req.ID),
	).Take()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("插件不存在")
		}
		return nil, fmt.Errorf("查询插件失败: %w", err)
	}

	// 更新启用状态
	_, err = repository.Plugins.Where(
		repository.Plugins.ID.Eq(req.ID),
	).Updates(map[string]interface{}{
		"is_enable":  req.IsEnable,
		"updated_at": time.Now(),
	})
	if err != nil {
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
func (s *PluginService) ExecutePluginsForEnv(envID int64, envValue string) (*plugin.ExecutionResult, error) {
	// 查询该环境变量启用的插件
	var results []struct {
		model.EnvPlugins
		Plugins model.Plugins `gorm:"embedded;embeddedPrefix:plugins_"`
	}

	err := repository.EnvPlugins.WithContext(context.Background()).
		Select(repository.EnvPlugins.ALL, repository.Plugins.ALL).
		LeftJoin(repository.Plugins, repository.EnvPlugins.PluginID.EqCol(repository.Plugins.ID)).
		Where(
			repository.EnvPlugins.EnvID.Eq(envID),
			repository.EnvPlugins.IsEnable.Is(true),
			repository.Plugins.IsEnable.Is(true),
		).
		Order(repository.EnvPlugins.ExecutionOrder.Asc()).
		Scan(&results)

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
		return &plugin.ExecutionResult{
			Success:    true,
			OutputData: outputBytes,
		}, nil
	}

	// 依次执行插件
	var lastResult *plugin.ExecutionResult
	for _, item := range results {
		// 构建执行上下文
		var configData []byte
		if item.Config != nil {
			configData = []byte(*item.Config)
		} else {
			configData = []byte("{}")
		}

		execCtx := &plugin.ExecutionContext{
			PluginID:  item.PluginID,
			EnvID:     envID,
			EnvValue:  envValue,
			Config:    configData,
			Timestamp: time.Now().Unix(),
		}

		// 执行插件
		timeout := time.Duration(item.Plugins.ExecutionTimeout) * time.Millisecond
		result := s.engine.Execute(context.Background(), item.Plugins.ScriptContent, execCtx, timeout)

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
func (s *PluginService) logPluginExecution(pluginID, envID int64, result *plugin.ExecutionResult) {
	status := "success"
	if !result.Success {
		status = "error"
	}

	// 处理可选字段
	var outputData, errorMessage, stackTrace *string
	if len(result.OutputData) > 0 {
		outputDataStr := string(result.OutputData)
		outputData = &outputDataStr
	}
	if result.ErrorMessage != "" {
		errorMessage = &result.ErrorMessage
	}
	if result.StackTrace != "" {
		stackTrace = &result.StackTrace
	}

	logModel := &model.PluginExecutionLogs{
		CreatedAt:       time.Now(),
		PluginID:        pluginID,
		EnvID:           envID,
		ExecutionStatus: status,
		ExecutionTime:   int32(result.ExecutionTime),
		InputData:       nil, // 环境变量值不记录在日志中（敏感信息）
		OutputData:      outputData,
		ErrorMessage:    errorMessage,
		StackTrace:      stackTrace,
	}

	// 异步记录日志，不影响主流程
	go func() {
		if err := repository.PluginExecutionLogs.WithContext(context.Background()).Create(logModel); err != nil {
			config.Log.Warn(err.Error()) // 仅做记录
		}
	}()
}

// BindPluginToEnv 绑定插件到环境变量
func (s *PluginService) BindPluginToEnv(req schema.BindPluginToEnvRequest) (*schema.BindPluginToEnvResponse, error) {
	// 检查插件是否存在
	_, err := repository.Plugins.Where(
		repository.Plugins.ID.Eq(req.PluginID),
	).Take()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("插件不存在")
		}
		return nil, fmt.Errorf("查询插件失败: %w", err)
	}

	// 检查环境变量是否存在
	_, err = repository.Envs.Where(
		repository.Envs.ID.Eq(req.EnvID),
	).Take()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("环境变量不存在")
		}
		return nil, fmt.Errorf("查询环境变量失败: %w", err)
	}

	// 检查是否已经绑定
	existingBinding, err := repository.EnvPlugins.Where(
		repository.EnvPlugins.PluginID.Eq(req.PluginID),
		repository.EnvPlugins.EnvID.Eq(req.EnvID),
	).Take()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("查询绑定关系失败: %w", err)
	}

	// 设置默认执行顺序
	executionOrder := req.ExecutionOrder
	if executionOrder == 0 {
		executionOrder = 100
	}

	// 准备配置数据
	var configStr *string
	if req.Config != "" {
		configStr = &req.Config
	}

	if existingBinding != nil {
		// 如果已经绑定，更新配置
		_, err = repository.EnvPlugins.Where(
			repository.EnvPlugins.PluginID.Eq(req.PluginID),
			repository.EnvPlugins.EnvID.Eq(req.EnvID),
		).Updates(map[string]interface{}{
			"config":          configStr,
			"execution_order": executionOrder,
			"is_enable":       true,
			"updated_at":      time.Now(),
		})
		if err != nil {
			return nil, fmt.Errorf("更新绑定关系失败: %w", err)
		}
	} else {
		// 创建新的绑定关系
		now := time.Now()
		binding := &model.EnvPlugins{
			CreatedAt:      now,
			UpdatedAt:      now,
			PluginID:       req.PluginID,
			EnvID:          req.EnvID,
			IsEnable:       true,
			ExecutionOrder: executionOrder,
			Config:         configStr,
		}

		if err = repository.EnvPlugins.WithContext(context.Background()).Create(binding); err != nil {
			return nil, fmt.Errorf("创建绑定关系失败: %w", err)
		}
	}

	return &schema.BindPluginToEnvResponse{
		Message: "插件绑定成功",
	}, nil
}

// UnbindPluginFromEnv 解绑插件与环境变量
func (s *PluginService) UnbindPluginFromEnv(req schema.UnbindPluginFromEnvRequest) (*schema.UnbindPluginFromEnvResponse, error) {
	// 检查绑定关系是否存在
	_, err := repository.EnvPlugins.Where(
		repository.EnvPlugins.PluginID.Eq(req.PluginID),
		repository.EnvPlugins.EnvID.Eq(req.EnvID),
	).Take()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("绑定关系不存在")
		}
		return nil, fmt.Errorf("查询绑定关系失败: %w", err)
	}

	// 删除绑定关系
	_, err = repository.EnvPlugins.Where(
		repository.EnvPlugins.PluginID.Eq(req.PluginID),
		repository.EnvPlugins.EnvID.Eq(req.EnvID),
	).Delete()
	if err != nil {
		return nil, fmt.Errorf("删除绑定关系失败: %w", err)
	}

	return &schema.UnbindPluginFromEnvResponse{
		Message: "插件解绑成功",
	}, nil
}

// GetPluginEnvs 获取插件关联环境变量
func (s *PluginService) GetPluginEnvs(req schema.GetPluginEnvsRequest) (*schema.GetPluginEnvsResponse, error) {
	// 检查插件是否存在
	_, err := repository.Plugins.Where(
		repository.Plugins.ID.Eq(req.PluginID),
	).Take()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("插件不存在")
		}
		return nil, fmt.Errorf("查询插件失败: %w", err)
	}

	// 查询插件关联的环境变量
	var results []struct {
		model.EnvPlugins
		Envs model.Envs `gorm:"embedded;embeddedPrefix:envs_"`
	}

	err = repository.EnvPlugins.WithContext(context.Background()).
		Select(repository.EnvPlugins.ALL, repository.Envs.ALL).
		LeftJoin(repository.Envs, repository.EnvPlugins.EnvID.EqCol(repository.Envs.ID)).
		Where(repository.EnvPlugins.PluginID.Eq(req.PluginID)).
		Scan(&results)

	if err != nil {
		return nil, fmt.Errorf("查询插件环境变量关联失败: %w", err)
	}

	// 转换为响应格式
	envs := make([]schema.PluginEnvRelationInfo, 0, len(results))
	for _, result := range results {
		var configData string
		if result.Config != nil {
			configData = *result.Config
		}

		envs = append(envs, schema.PluginEnvRelationInfo{
			EnvID:          result.EnvID,
			EnvName:        result.Envs.Name,
			IsEnable:       result.IsEnable,
			ExecutionOrder: result.ExecutionOrder,
			Config:         configData,
			CreatedAt:      result.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return &schema.GetPluginEnvsResponse{
		PluginID: req.PluginID,
		Envs:     envs,
	}, nil
}

// GetPluginExecutionLogs 获取插件执行日志
func (s *PluginService) GetPluginExecutionLogs(req schema.GetPluginExecutionLogsRequest) (*schema.GetPluginExecutionLogsResponse, error) {
	// 设置默认分页参数
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	// 构建查询条件
	query := repository.PluginExecutionLogs.WithContext(context.Background()).
		Select(repository.PluginExecutionLogs.ALL, repository.Plugins.Name.As("plugin_name"), repository.Envs.Name.As("env_name")).
		LeftJoin(repository.Plugins, repository.PluginExecutionLogs.PluginID.EqCol(repository.Plugins.ID)).
		LeftJoin(repository.Envs, repository.PluginExecutionLogs.EnvID.EqCol(repository.Envs.ID))

	// 按插件ID筛选
	if req.PluginID != nil {
		query = query.Where(repository.PluginExecutionLogs.PluginID.Eq(*req.PluginID))
	}

	// 按环境变量ID筛选
	if req.EnvID != nil {
		query = query.Where(repository.PluginExecutionLogs.EnvID.Eq(*req.EnvID))
	}

	// 按执行状态筛选
	if req.ExecutionStatus != "" {
		query = query.Where(repository.PluginExecutionLogs.ExecutionStatus.Eq(req.ExecutionStatus))
	}

	// 按时间范围筛选
	if req.StartTime != "" {
		if startTime, err := time.Parse("2006-01-02 15:04:05", req.StartTime); err == nil {
			query = query.Where(repository.PluginExecutionLogs.CreatedAt.Gte(startTime))
		}
	}
	if req.EndTime != "" {
		if endTime, err := time.Parse("2006-01-02 15:04:05", req.EndTime); err == nil {
			query = query.Where(repository.PluginExecutionLogs.CreatedAt.Lte(endTime))
		}
	}

	// 获取总数
	total, err := query.Count()
	if err != nil {
		return nil, fmt.Errorf("查询日志总数失败: %w", err)
	}

	// 分页查询
	offset := (req.Page - 1) * req.PageSize
	logs, err := query.Offset(offset).Limit(req.PageSize).Order(repository.PluginExecutionLogs.CreatedAt.Desc()).Find()
	if err != nil {
		return nil, fmt.Errorf("查询执行日志失败: %w", err)
	}

	// 转换为响应格式
	list := make([]schema.PluginExecutionLogInfo, 0, len(logs))
	for _, log := range logs {
		// 处理可选字段
		var inputData, outputData string
		if log.InputData != nil {
			inputData = *log.InputData
		}
		if log.OutputData != nil {
			outputData = *log.OutputData
		}

		var errorMessage string
		if log.ErrorMessage != nil {
			errorMessage = *log.ErrorMessage
		}

		list = append(list, schema.PluginExecutionLogInfo{
			ID:              log.ID,
			PluginID:        log.PluginID,
			PluginName:      "", // 需要通过JOIN获取
			EnvID:           log.EnvID,
			EnvName:         "", // 需要通过JOIN获取
			ExecutionStatus: log.ExecutionStatus,
			ExecutionTime:   int(log.ExecutionTime),
			InputData:       inputData,
			OutputData:      outputData,
			ErrorMessage:    errorMessage,
			CreatedAt:       log.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return &schema.GetPluginExecutionLogsResponse{
		Total: total,
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
