package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/nuanxinqing123/QLToolsV2/internal/model"
	"github.com/nuanxinqing123/QLToolsV2/internal/repository"
	"github.com/nuanxinqing123/QLToolsV2/internal/schema"
	"gorm.io/gorm"
)

type EnvService struct{}

// NewEnvService 创建 EnvService
func NewEnvService() *EnvService {
	return &EnvService{}
}

// AddEnv 添加环境变量
func (s *EnvService) AddEnv(req schema.AddEnvRequest) (*schema.AddEnvResponse, error) {
	// 检查环境变量名称是否已存在
	existingEnv, err := repository.Envs.Where(
		repository.Envs.Name.Eq(req.Name),
	).Take()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("查询环境变量失败: %w", err)
	}
	if existingEnv != nil {
		return nil, errors.New("环境变量名称已存在")
	}

	now := time.Now()
	env := &model.Envs{
		CreatedAt:       now,
		UpdatedAt:       now,
		Name:            req.Name,
		Remarks:         req.Remarks,
		Quantity:        req.Quantity,
		Regex:           req.Regex,
		Mode:            req.Mode,
		RegexUpdate:     req.RegexUpdate,
		IsAutoEnvEnable: req.IsAutoEnvEnable, // 是否自动启用提交的变量
		EnableKey:       req.EnableKey,
		CdkLimit:        req.CdkLimit, // 单次消耗卡密额度
		IsPrompt:        req.IsPrompt,
		PromptLevel:     req.PromptLevel,
		PromptContent:   req.PromptContent,
		IsEnable:        true, // 默认启用
	}

	// 创建环境变量记录
	if err = repository.Envs.WithContext(context.Background()).Create(env); err != nil {
		return nil, fmt.Errorf("创建环境变量失败: %w", err)
	}

	return &schema.AddEnvResponse{
		ID:      env.ID,
		Message: "环境变量添加成功",
	}, nil
}

// UpdateEnv 更新环境变量
func (s *EnvService) UpdateEnv(req schema.UpdateEnvRequest) (*schema.UpdateEnvResponse, error) {
	// 查询环境变量是否存在
	env, err := repository.Envs.Where(
		repository.Envs.ID.Eq(req.ID),
	).Take()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("环境变量不存在")
		}
		return nil, fmt.Errorf("查询环境变量失败: %w", err)
	}

	// 检查名称是否与其他环境变量冲突
	if req.Name != env.Name {
		existingEnv, err := repository.Envs.Where(
			repository.Envs.Name.Eq(req.Name),
			repository.Envs.ID.Neq(req.ID),
		).Take()
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("查询环境变量失败: %w", err)
		}
		if existingEnv != nil {
			return nil, errors.New("环境变量名称已存在")
		}
	}

	// 构建更新数据
	updates := map[string]interface{}{
		"name":               req.Name,
		"remarks":            req.Remarks,
		"quantity":           req.Quantity,
		"regex":              req.Regex,
		"mode":               req.Mode,
		"regex_update":       req.RegexUpdate,
		"is_auto_env_enable": req.IsAutoEnvEnable, // 是否自动启用提交的变量
		"enable_key":         req.EnableKey,
		"cdk_limit":          req.CdkLimit, // 单次消耗卡密额度
		"is_prompt":          req.IsPrompt,
		"prompt_level":       req.PromptLevel,
		"prompt_content":     req.PromptContent,
		"updated_at":         time.Now(),
	}

	// 如果提供了启用状态，则更新
	if req.IsEnable != nil {
		updates["is_enable"] = *req.IsEnable
	}

	// 执行更新
	_, err = repository.Envs.Where(
		repository.Envs.ID.Eq(req.ID),
	).Updates(updates)
	if err != nil {
		return nil, fmt.Errorf("更新环境变量失败: %w", err)
	}

	return &schema.UpdateEnvResponse{
		Message: "环境变量更新成功",
	}, nil
}

// GetEnv 获取单个环境变量信息
func (s *EnvService) GetEnv(id int64) (*schema.GetEnvResponse, error) {
	env, err := repository.Envs.Where(
		repository.Envs.ID.Eq(id),
	).Take()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("环境变量不存在")
		}
		return nil, fmt.Errorf("查询环境变量失败: %w", err)
	}

	return &schema.GetEnvResponse{
		ID:              env.ID,
		Name:            env.Name,
		Remarks:         env.Remarks,
		Quantity:        env.Quantity,
		Regex:           env.Regex,
		Mode:            env.Mode,
		RegexUpdate:     env.RegexUpdate,
		IsAutoEnvEnable: env.IsAutoEnvEnable, // 是否自动启用提交的变量
		EnableKey:       env.EnableKey,
		CdkLimit:        env.CdkLimit, // 单次消耗卡密额度
		IsPrompt:        env.IsPrompt,
		PromptLevel:     env.PromptLevel,
		PromptContent:   env.PromptContent,
		IsEnable:        env.IsEnable,
		CreatedAt:       env.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:       env.UpdatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

// GetEnvList 获取环境变量列表
func (s *EnvService) GetEnvList(req schema.GetEnvListRequest) (*schema.GetEnvListResponse, error) {
	// 设置默认分页参数
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	// 构建查询条件
	query := repository.Envs.WithContext(context.Background())

	// 按名称模糊搜索
	if req.Name != "" {
		query = query.Where(repository.Envs.Name.Like("%" + req.Name + "%"))
	}

	// 按启用状态筛选
	if req.IsEnable != nil {
		query = query.Where(repository.Envs.IsEnable.Is(*req.IsEnable))
	}

	// 按模式筛选
	if req.Mode != nil {
		query = query.Where(repository.Envs.Mode.Eq(*req.Mode))
	}

	// 获取总数
	total, err := query.Count()
	if err != nil {
		return nil, fmt.Errorf("查询环境变量总数失败: %w", err)
	}

	// 分页查询
	offset := (req.Page - 1) * req.PageSize
	envs, err := query.Offset(offset).Limit(req.PageSize).Order(repository.Envs.CreatedAt.Desc()).Find()
	if err != nil {
		return nil, fmt.Errorf("查询环境变量列表失败: %w", err)
	}

	// 转换为响应格式
	list := make([]schema.GetEnvResponse, 0, len(envs))
	for _, env := range envs {
		list = append(list, schema.GetEnvResponse{
			ID:              env.ID,
			Name:            env.Name,
			Remarks:         env.Remarks,
			Quantity:        env.Quantity,
			Regex:           env.Regex,
			Mode:            env.Mode,
			RegexUpdate:     env.RegexUpdate,
			IsAutoEnvEnable: env.IsAutoEnvEnable, // 是否自动启用提交的变量
			EnableKey:       env.EnableKey,
			CdkLimit:        env.CdkLimit, // 单次消耗卡密额度
			IsPrompt:        env.IsPrompt,
			PromptLevel:     env.PromptLevel,
			PromptContent:   env.PromptContent,
			IsEnable:        env.IsEnable,
			CreatedAt:       env.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:       env.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return &schema.GetEnvListResponse{
		Total: total,
		List:  list,
	}, nil
}

// DeleteEnv 删除环境变量
func (s *EnvService) DeleteEnv(req schema.DeleteEnvConfigRequest) (*schema.DeleteEnvConfigResponse, error) {
	// 检查环境变量是否存在
	_, err := repository.Envs.Where(
		repository.Envs.ID.Eq(req.ID),
	).Take()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("环境变量不存在")
		}
		return nil, fmt.Errorf("查询环境变量失败: %w", err)
	}

	// 删除相关的面板绑定关系
	_, err = repository.EnvPanels.Where(
		repository.EnvPanels.EnvID.Eq(req.ID),
	).Delete()
	if err != nil {
		return nil, fmt.Errorf("删除环境变量面板绑定关系失败: %w", err)
	}

	// 执行软删除
	_, err = repository.Envs.Where(
		repository.Envs.ID.Eq(req.ID),
	).Delete()
	if err != nil {
		return nil, fmt.Errorf("删除环境变量失败: %w", err)
	}

	return &schema.DeleteEnvConfigResponse{
		Message: "环境变量删除成功",
	}, nil
}

// ToggleEnvStatus 切换环境变量启用状态
func (s *EnvService) ToggleEnvStatus(req schema.ToggleEnvStatusRequest) (*schema.ToggleEnvStatusResponse, error) {
	// 检查环境变量是否存在
	_, err := repository.Envs.Where(
		repository.Envs.ID.Eq(req.ID),
	).Take()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("环境变量不存在")
		}
		return nil, fmt.Errorf("查询环境变量失败: %w", err)
	}

	// 更新启用状态
	_, err = repository.Envs.Where(
		repository.Envs.ID.Eq(req.ID),
	).Updates(map[string]interface{}{
		"is_enable":  req.IsEnable,
		"updated_at": time.Now(),
	})
	if err != nil {
		return nil, fmt.Errorf("更新环境变量状态失败: %w", err)
	}

	status := "禁用"
	if req.IsEnable {
		status = "启用"
	}

	return &schema.ToggleEnvStatusResponse{
		Message: fmt.Sprintf("环境变量已%s", status),
	}, nil
}

// UpdateEnvPanels 更新环境变量面板绑定关系
func (s *EnvService) UpdateEnvPanels(req schema.UpdateEnvPanelsRequest) (*schema.UpdateEnvPanelsResponse, error) {
	// 检查环境变量是否存在
	_, err := repository.Envs.Where(
		repository.Envs.ID.Eq(req.EnvID),
	).Take()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("环境变量不存在")
		}
		return nil, fmt.Errorf("查询环境变量失败: %w", err)
	}

	// 检查所有面板是否存在（如果面板ID列表不为空）
	if len(req.PanelIDs) > 0 {
		for _, panelID := range req.PanelIDs {
			_, err := repository.Panels.Where(
				repository.Panels.ID.Eq(panelID),
			).Take()
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return nil, fmt.Errorf("面板ID %d 不存在", panelID)
				}
				return nil, fmt.Errorf("查询面板失败: %w", err)
			}
		}
	}

	// 先删除该环境变量的所有现有绑定关系
	_, err = repository.EnvPanels.Where(
		repository.EnvPanels.EnvID.Eq(req.EnvID),
	).Delete()
	if err != nil {
		return nil, fmt.Errorf("删除现有绑定关系失败: %w", err)
	}

	// 创建新的绑定关系（如果面板ID列表不为空）
	if len(req.PanelIDs) > 0 {
		var envPanels []*model.EnvPanels
		for _, panelID := range req.PanelIDs {
			envPanels = append(envPanels, &model.EnvPanels{
				EnvID:   req.EnvID,
				PanelID: panelID,
			})
		}

		err = repository.EnvPanels.WithContext(context.Background()).CreateInBatches(envPanels, 100)
		if err != nil {
			return nil, fmt.Errorf("创建绑定关系失败: %w", err)
		}
	}

	message := "环境变量面板绑定关系更新成功"
	if len(req.PanelIDs) == 0 {
		message = "环境变量已解绑所有面板"
	}

	return &schema.UpdateEnvPanelsResponse{
		Message: message,
	}, nil
}

// GetEnvPanels 获取环境变量关联的面板
func (s *EnvService) GetEnvPanels(req schema.GetEnvPanelsRequest) (*schema.GetEnvPanelsResponse, error) {
	// 检查环境变量是否存在
	_, err := repository.Envs.Where(
		repository.Envs.ID.Eq(req.EnvID),
	).Take()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("环境变量不存在")
		}
		return nil, fmt.Errorf("查询环境变量失败: %w", err)
	}

	// 查询关联的面板ID
	envPanels, err := repository.EnvPanels.Where(
		repository.EnvPanels.EnvID.Eq(req.EnvID),
	).Find()
	if err != nil {
		return nil, fmt.Errorf("查询环境变量关联面板失败: %w", err)
	}

	// 提取面板ID列表
	panelIDs := make([]int64, 0, len(envPanels))
	for _, envPanel := range envPanels {
		panelIDs = append(panelIDs, envPanel.PanelID)
	}

	return &schema.GetEnvPanelsResponse{
		EnvID:    req.EnvID,
		PanelIDs: panelIDs,
	}, nil
}
