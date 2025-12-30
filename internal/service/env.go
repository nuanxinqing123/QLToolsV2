package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/nuanxinqing123/QLToolsV2/internal/app/config"
	"github.com/nuanxinqing123/QLToolsV2/internal/data/ent"
	"github.com/nuanxinqing123/QLToolsV2/internal/data/ent/env"
	"github.com/nuanxinqing123/QLToolsV2/internal/data/ent/envplugin"
	"github.com/nuanxinqing123/QLToolsV2/internal/data/ent/panel"
	"github.com/nuanxinqing123/QLToolsV2/internal/schema"
)

type EnvService struct{}

// NewEnvService 创建 EnvService
func NewEnvService() *EnvService {
	return &EnvService{}
}

// AddEnv 添加环境变量
func (s *EnvService) AddEnv(req schema.AddEnvRequest) (*schema.AddEnvResponse, error) {
	ctx := context.Background()
	// 检查环境变量名称是否已存在
	exists, err := config.Ent.Env.Query().
		Where(env.NameEQ(req.Name)).
		Exist(ctx)
	if err != nil {
		return nil, fmt.Errorf("查询环境变量失败: %w", err)
	}
	if exists {
		return nil, errors.New("环境变量名称已存在")
	}

	// 创建环境变量记录
	e, err := config.Ent.Env.Create().
		SetName(req.Name).
		SetNillableRemarks(req.Remarks).
		SetQuantity(req.Quantity).
		SetNillableRegex(req.Regex).
		SetMode(req.Mode).
		SetNillableRegexUpdate(req.RegexUpdate).
		SetIsAutoEnvEnable(req.IsAutoEnvEnable).
		SetEnableKey(req.EnableKey).
		SetCdkLimit(req.CdkLimit).
		SetIsPrompt(req.IsPrompt).
		SetNillablePromptLevel(req.PromptLevel).
		SetNillablePromptContent(req.PromptContent).
		SetIsEnable(true).
		SetCreatedAt(time.Now()).
		SetUpdatedAt(time.Now()).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("创建环境变量失败: %w", err)
	}

	return &schema.AddEnvResponse{
		ID:      e.ID,
		Message: "环境变量添加成功",
	}, nil
}

// UpdateEnv 更新环境变量
func (s *EnvService) UpdateEnv(req schema.UpdateEnvRequest) (*schema.UpdateEnvResponse, error) {
	ctx := context.Background()
	// 查询环境变量是否存在
	e, err := config.Ent.Env.Get(ctx, req.ID)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.New("环境变量不存在")
		}
		return nil, fmt.Errorf("查询环境变量失败: %w", err)
	}

	// 检查名称是否与其他环境变量冲突
	if req.Name != e.Name {
		exists, err := config.Ent.Env.Query().
			Where(env.And(
				env.NameEQ(req.Name),
				env.IDNEQ(req.ID),
			)).
			Exist(ctx)
		if err != nil {
			return nil, fmt.Errorf("查询环境变量失败: %w", err)
		}
		if exists {
			return nil, errors.New("环境变量名称已存在")
		}
	}

	// 执行更新
	updater := config.Ent.Env.UpdateOneID(req.ID).
		SetName(req.Name).
		SetNillableRemarks(req.Remarks).
		SetQuantity(req.Quantity).
		SetNillableRegex(req.Regex).
		SetMode(req.Mode).
		SetNillableRegexUpdate(req.RegexUpdate).
		SetIsAutoEnvEnable(req.IsAutoEnvEnable).
		SetEnableKey(req.EnableKey).
		SetCdkLimit(req.CdkLimit).
		SetIsPrompt(req.IsPrompt).
		SetNillablePromptLevel(req.PromptLevel).
		SetNillablePromptContent(req.PromptContent).
		SetUpdatedAt(time.Now())

	if req.IsEnable != nil {
		updater.SetIsEnable(*req.IsEnable)
	}

	if err := updater.Exec(ctx); err != nil {
		return nil, fmt.Errorf("更新环境变量失败: %w", err)
	}

	return &schema.UpdateEnvResponse{
		Message: "环境变量更新成功",
	}, nil
}

// GetEnv 获取单个环境变量信息
func (s *EnvService) GetEnv(id int64) (*schema.GetEnvResponse, error) {
	ctx := context.Background()
	e, err := config.Ent.Env.Get(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.New("环境变量不存在")
		}
		return nil, fmt.Errorf("查询环境变量失败: %w", err)
	}

	return &schema.GetEnvResponse{
		ID:              e.ID,
		Name:            e.Name,
		Remarks:         e.Remarks,
		Quantity:        e.Quantity,
		Regex:           e.Regex,
		Mode:            e.Mode,
		RegexUpdate:     e.RegexUpdate,
		IsAutoEnvEnable: e.IsAutoEnvEnable,
		EnableKey:       e.EnableKey,
		CdkLimit:        e.CdkLimit,
		IsPrompt:        e.IsPrompt,
		PromptLevel:     e.PromptLevel,
		PromptContent:   e.PromptContent,
		IsEnable:        e.IsEnable,
		CreatedAt:       e.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:       e.UpdatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

// GetEnvList 获取环境变量列表
func (s *EnvService) GetEnvList(req schema.GetEnvListRequest) (*schema.GetEnvListResponse, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	ctx := context.Background()
	query := config.Ent.Env.Query()

	if req.Name != "" {
		query.Where(env.NameContains(req.Name))
	}
	if req.IsEnable != nil {
		query.Where(env.IsEnableEQ(*req.IsEnable))
	}
	if req.Mode != nil {
		query.Where(env.ModeEQ(*req.Mode))
	}

	total, err := query.Count(ctx)
	if err != nil {
		return nil, fmt.Errorf("查询环境变量总数失败: %w", err)
	}

	offset := (req.Page - 1) * req.PageSize
	envs, err := query.Offset(offset).
		Limit(req.PageSize).
		Order(ent.Desc(env.FieldCreatedAt)).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("查询环境变量列表失败: %w", err)
	}

	list := make([]schema.GetEnvResponse, 0, len(envs))
	for _, e := range envs {
		list = append(list, schema.GetEnvResponse{
			ID:              e.ID,
			Name:            e.Name,
			Remarks:         e.Remarks,
			Quantity:        e.Quantity,
			Regex:           e.Regex,
			Mode:            e.Mode,
			RegexUpdate:     e.RegexUpdate,
			IsAutoEnvEnable: e.IsAutoEnvEnable,
			EnableKey:       e.EnableKey,
			CdkLimit:        e.CdkLimit,
			IsPrompt:        e.IsPrompt,
			PromptLevel:     e.PromptLevel,
			PromptContent:   e.PromptContent,
			IsEnable:        e.IsEnable,
			CreatedAt:       e.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:       e.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return &schema.GetEnvListResponse{
		Total: int64(total),
		List:  list,
	}, nil
}

// DeleteEnv 删除环境变量
func (s *EnvService) DeleteEnv(req schema.DeleteEnvConfigRequest) (*schema.DeleteEnvConfigResponse, error) {
	ctx := context.Background()
	// 检查环境变量是否存在
	_, err := config.Ent.Env.Get(ctx, req.ID)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.New("环境变量不存在")
		}
		return nil, fmt.Errorf("查询环境变量失败: %w", err)
	}

	// Ent 删除会自动处理多对多关联表的清除 (如果定义了正确的外键/级联)
	// 在我们的 Schema 中，env_panels 是通过 StorageKey 配置的
	if err := config.Ent.Env.DeleteOneID(req.ID).Exec(ctx); err != nil {
		return nil, fmt.Errorf("删除环境变量失败: %w", err)
	}

	return &schema.DeleteEnvConfigResponse{
		Message: "环境变量删除成功",
	}, nil
}

// ToggleEnvStatus 切换环境变量启用状态
func (s *EnvService) ToggleEnvStatus(req schema.ToggleEnvStatusRequest) (*schema.ToggleEnvStatusResponse, error) {
	ctx := context.Background()
	_, err := config.Ent.Env.Get(ctx, req.ID)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.New("环境变量不存在")
		}
		return nil, fmt.Errorf("查询环境变量失败: %w", err)
	}

	if err := config.Ent.Env.UpdateOneID(req.ID).
		SetIsEnable(req.IsEnable).
		SetUpdatedAt(time.Now()).
		Exec(ctx); err != nil {
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
	ctx := context.Background()
	// 检查环境变量是否存在
	_, err := config.Ent.Env.Get(ctx, req.EnvID)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.New("环境变量不存在")
		}
		return nil, fmt.Errorf("查询环境变量失败: %w", err)
	}

	// 验证面板是否存在
	if len(req.PanelIDs) > 0 {
		cnt, err := config.Ent.Panel.Query().Where(panel.IDIn(req.PanelIDs...)).Count(ctx)
		if err != nil {
			return nil, fmt.Errorf("查询面板失败: %w", err)
		}
		if cnt != len(req.PanelIDs) {
			return nil, errors.New("部分面板不存在")
		}
	}

	// 更新关联关系 (Ent 会自动处理解绑和新绑定)
	err = config.Ent.Env.UpdateOneID(req.EnvID).
		ClearPanels().
		AddPanelIDs(req.PanelIDs...).
		Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("更新绑定关系失败: %w", err)
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
	ctx := context.Background()
	e, err := config.Ent.Env.Query().
		Where(env.IDEQ(req.EnvID)).
		WithPanels().
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.New("环境变量不存在")
		}
		return nil, fmt.Errorf("查询环境变量关联面板失败: %w", err)
	}

	panelIDs := make([]int64, 0, len(e.Edges.Panels))
	for _, p := range e.Edges.Panels {
		panelIDs = append(panelIDs, p.ID)
	}

	return &schema.GetEnvPanelsResponse{
		EnvID:    req.EnvID,
		PanelIDs: panelIDs,
	}, nil
}

// GetEnvPlugins 获取环境变量关联的插件
func (s *EnvService) GetEnvPlugins(req schema.GetEnvPluginsRequest) (*schema.GetEnvPluginsResponse, error) {
	ctx := context.Background()
	// 查询关联的插件信息（通过中间表 EnvPlugin）
	eps, err := config.Ent.EnvPlugin.Query().
		Where(envplugin.EnvIDEQ(req.EnvID)).
		WithPlugin().
		Order(ent.Asc(envplugin.FieldExecutionOrder)).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("查询环境变量关联插件失败: %w", err)
	}

	plugins := make([]schema.EnvPluginRelationInfo, 0, len(eps))
	for _, ep := range eps {
		configStr := ""
		if ep.Config != nil {
			configStr = *ep.Config
		}

		plugins = append(plugins, schema.EnvPluginRelationInfo{
			PluginID:       ep.PluginID,
			PluginName:     ep.Edges.Plugin.Name,
			IsEnable:       ep.IsEnable,
			ExecutionOrder: ep.ExecutionOrder,
			Config:         configStr,
			CreatedAt:      ep.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return &schema.GetEnvPluginsResponse{
		EnvID:   req.EnvID,
		Plugins: plugins,
	}, nil
}
