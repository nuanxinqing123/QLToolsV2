package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/nuanxinqing123/QLToolsV2/internal/app/config"
	"github.com/nuanxinqing123/QLToolsV2/internal/data/ent"
	"github.com/nuanxinqing123/QLToolsV2/internal/data/ent/panel"
	"github.com/nuanxinqing123/QLToolsV2/internal/pkg/qinglong"
	"github.com/nuanxinqing123/QLToolsV2/internal/schema"
)

type PanelService struct {
	pluginService *PluginService
}

// NewPanelService 创建 PanelService
func NewPanelService() *PanelService {
	return &PanelService{
		pluginService: NewPluginService(),
	}
}

// AddPanel 添加面板
func (s *PanelService) AddPanel(req schema.AddPanelRequest) (*schema.AddPanelResponse, error) {
	ctx := context.Background()
	// 检查面板名称是否已存在
	exists, err := config.Ent.Panel.Query().
		Where(panel.NameEQ(req.Name)).
		Exist(ctx)
	if err != nil {
		return nil, fmt.Errorf("查询面板失败: %w", err)
	}
	if exists {
		return nil, errors.New("面板名称已存在")
	}

	// 获取面板Token
	qlConfig := qinglong.NewConfig(req.URL, req.ClientID, req.ClientSecret)
	tokenResp, err := qlConfig.GetConfig()
	if err != nil {
		return nil, fmt.Errorf("连接面板失败，无法获取Token: %w", err)
	}

	// 检查API响应状态
	if tokenResp.Code != 200 {
		return nil, fmt.Errorf("获取面板Token失败，错误信息: %s", tokenResp.Message)
	}

	// 创建面板记录
	p, err := config.Ent.Panel.Create().
		SetName(req.Name).
		SetURL(req.URL).
		SetClientID(req.ClientID).
		SetClientSecret(req.ClientSecret).
		SetIsEnable(req.IsEnable).
		SetToken(tokenResp.Data.Token).
		SetParams(int32(tokenResp.Data.Expiration)).
		SetCreatedAt(time.Now()).
		SetUpdatedAt(time.Now()).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("创建面板失败: %w", err)
	}

	return &schema.AddPanelResponse{
		ID:      p.ID,
		Message: "面板添加成功",
	}, nil
}

// UpdatePanel 更新面板
func (s *PanelService) UpdatePanel(req schema.UpdatePanelRequest) (*schema.UpdatePanelResponse, error) {
	ctx := context.Background()
	// 查询面板是否存在
	p, err := config.Ent.Panel.Get(ctx, req.ID)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.New("面板不存在")
		}
		return nil, fmt.Errorf("查询面板失败: %w", err)
	}

	// 检查名称是否与其他面板冲突
	if req.Name != p.Name {
		exists, err := config.Ent.Panel.Query().
			Where(panel.And(
				panel.NameEQ(req.Name),
				panel.IDNEQ(req.ID),
			)).
			Exist(ctx)
		if err != nil {
			return nil, fmt.Errorf("查询面板失败: %w", err)
		}
		if exists {
			return nil, errors.New("面板名称已存在")
		}
	}

	// 执行更新
	updater := config.Ent.Panel.UpdateOneID(req.ID).
		SetName(req.Name).
		SetURL(req.URL).
		SetClientID(req.ClientID).
		SetClientSecret(req.ClientSecret).
		SetIsEnable(req.IsEnable).
		SetUpdatedAt(time.Now())

	// 如果连接信息发生变化，需要重新获取Token
	needRefreshToken := req.URL != p.URL || req.ClientID != p.ClientID || req.ClientSecret != p.ClientSecret
	if needRefreshToken {
		// 使用新的连接信息获取Token
		qlConfig := qinglong.NewConfig(req.URL, req.ClientID, req.ClientSecret)
		tokenResp, err := qlConfig.GetConfig()
		if err != nil {
			return nil, fmt.Errorf("连接面板失败，无法获取新Token: %w", err)
		}

		// 检查API响应状态
		if tokenResp.Code != 200 {
			return nil, fmt.Errorf("获取面板Token失败，错误信息: %s", tokenResp.Message)
		}

		// 更新Token相关字段
		updater.SetToken(tokenResp.Data.Token).
			SetParams(int32(tokenResp.Data.Expiration))
	}

	if err := updater.Exec(ctx); err != nil {
		return nil, fmt.Errorf("更新面板失败: %w", err)
	}

	return &schema.UpdatePanelResponse{
		Message: "面板更新成功",
	}, nil
}

// GetPanel 获取单个面板信息
func (s *PanelService) GetPanel(id int64) (*schema.GetPanelResponse, error) {
	ctx := context.Background()
	p, err := config.Ent.Panel.Get(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.New("面板不存在")
		}
		return nil, fmt.Errorf("查询面板失败: %w", err)
	}

	return &schema.GetPanelResponse{
		ID:           p.ID,
		Name:         p.Name,
		URL:          p.URL,
		ClientID:     p.ClientID,
		ClientSecret: p.ClientSecret,
		IsEnable:     p.IsEnable,
		Token:        p.Token,
		Params:       p.Params,
		CreatedAt:    p.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:    p.UpdatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

// GetPanelList 获取面板列表
func (s *PanelService) GetPanelList(req schema.GetPanelListRequest) (*schema.GetPanelListResponse, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	ctx := context.Background()
	query := config.Ent.Panel.Query()

	if req.Name != "" {
		query.Where(panel.NameContains(req.Name))
	}
	if req.IsEnable != nil {
		query.Where(panel.IsEnableEQ(*req.IsEnable))
	}

	total, err := query.Count(ctx)
	if err != nil {
		return nil, fmt.Errorf("查询面板总数失败: %w", err)
	}

	offset := (req.Page - 1) * req.PageSize
	panels, err := query.Offset(offset).
		Limit(req.PageSize).
		Order(ent.Desc(panel.FieldCreatedAt)).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("查询面板列表失败: %w", err)
	}

	list := make([]schema.GetPanelResponse, 0, len(panels))
	for _, p := range panels {
		list = append(list, schema.GetPanelResponse{
			ID:           p.ID,
			Name:         p.Name,
			URL:          p.URL,
			ClientID:     p.ClientID,
			ClientSecret: p.ClientSecret,
			IsEnable:     p.IsEnable,
			Token:        p.Token,
			Params:       p.Params,
			CreatedAt:    p.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:    p.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return &schema.GetPanelListResponse{
		Total: int64(total),
		List:  list,
	}, nil
}

// DeletePanel 删除面板
func (s *PanelService) DeletePanel(req schema.DeletePanelRequest) (*schema.DeletePanelResponse, error) {
	ctx := context.Background()
	// 检查面板是否存在
	_, err := config.Ent.Panel.Get(ctx, req.ID)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.New("面板不存在")
		}
		return nil, fmt.Errorf("查询面板失败: %w", err)
	}

	if err := config.Ent.Panel.DeleteOneID(req.ID).Exec(ctx); err != nil {
		return nil, fmt.Errorf("删除面板失败: %w", err)
	}

	return &schema.DeletePanelResponse{
		Message: "面板删除成功",
	}, nil
}

// TogglePanelStatus 切换面板启用状态
func (s *PanelService) TogglePanelStatus(req schema.TogglePanelStatusRequest) (*schema.TogglePanelStatusResponse, error) {
	ctx := context.Background()
	_, err := config.Ent.Panel.Get(ctx, req.ID)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.New("面板不存在")
		}
		return nil, fmt.Errorf("查询面板失败: %w", err)
	}

	if err := config.Ent.Panel.UpdateOneID(req.ID).
		SetIsEnable(req.IsEnable).
		SetUpdatedAt(time.Now()).
		Exec(ctx); err != nil {
		return nil, fmt.Errorf("更新面板状态失败: %w", err)
	}

	status := "禁用"
	if req.IsEnable {
		status = "启用"
	}

	return &schema.TogglePanelStatusResponse{
		Message: fmt.Sprintf("面板已%s", status),
	}, nil
}

// RefreshPanelToken 刷新面板Token
func (s *PanelService) RefreshPanelToken(req schema.RefreshPanelTokenRequest) (*schema.RefreshPanelTokenResponse, error) {
	ctx := context.Background()
	// 查询面板信息
	p, err := config.Ent.Panel.Get(ctx, req.ID)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.New("面板不存在")
		}
		return nil, fmt.Errorf("查询面板失败: %w", err)
	}

	// 使用面板的连接信息重新获取Token
	qlConfig := qinglong.NewConfig(p.URL, p.ClientID, p.ClientSecret)
	tokenResp, err := qlConfig.GetConfig()
	if err != nil {
		return nil, fmt.Errorf("连接面板失败，无法刷新Token: %w", err)
	}

	// 检查API响应状态
	if tokenResp.Code != 200 {
		return nil, fmt.Errorf("刷新面板Token失败，错误信息: %s", tokenResp.Message)
	}

	newToken := tokenResp.Data.Token
	newParams := int32(tokenResp.Data.Expiration)

	// 更新Token信息
	err = config.Ent.Panel.UpdateOneID(req.ID).
		SetToken(newToken).
		SetParams(newParams).
		SetUpdatedAt(time.Now()).
		Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("更新面板Token失败: %w", err)
	}

	return &schema.RefreshPanelTokenResponse{
		Message: "Token刷新成功",
		Token:   newToken,
	}, nil
}

// TestPanelConnection 测试面板连接
func (s *PanelService) TestPanelConnection(req schema.TestPanelConnectionRequest) (*schema.TestPanelConnectionResponse, error) {
	// 创建青龙配置实例
	qlConfig := qinglong.NewConfig(req.URL, req.ClientID, req.ClientSecret)

	// 尝试获取Token来测试连接
	tokenResp, err := qlConfig.GetConfig()
	if err != nil {
		// 连接失败，返回失败响应
		return &schema.TestPanelConnectionResponse{
			Success:     false,
			Message:     "连接失败",
			Token:       "",
			Expiration:  0,
			ResponseMsg: fmt.Sprintf("网络连接错误: %s", err.Error()),
		}, nil
	}

	// 检查API响应状态
	if tokenResp.Code != 200 {
		// API返回错误状态，连接失败
		return &schema.TestPanelConnectionResponse{
			Success:     false,
			Message:     "认证失败",
			Token:       "",
			Expiration:  0,
			ResponseMsg: fmt.Sprintf("API响应错误 (Code: %d): %s", tokenResp.Code, tokenResp.Message),
		}, nil
	}

	// 连接成功，返回成功响应
	return &schema.TestPanelConnectionResponse{
		Success:     true,
		Message:     "连接成功",
		Token:       tokenResp.Data.Token,
		Expiration:  tokenResp.Data.Expiration,
		ResponseMsg: "面板连接正常，认证成功",
	}, nil
}

// SubmitEnvToPanel 提交环境变量到面板（集成插件执行流程）
func (s *PanelService) SubmitEnvToPanel(panelID int64, envID int64, envValue string) (interface{}, error) {
	// 执行环境变量的插件验证 verification logic preserved
	result, err := s.pluginService.ExecutePluginsForEnv(envID, envValue)
	if err != nil {
		return nil, fmt.Errorf("执行环境变量插件失败: %w", err)
	}

	// 检查插件执行结果
	if !result.Success {
		return nil, fmt.Errorf("插件验证失败: %s", result.ErrorMessage)
	}

	// 解析插件返回的结果
	var pluginResult map[string]interface{}
	if len(result.OutputData) > 0 {
		if err := json.Unmarshal(result.OutputData, &pluginResult); err != nil {
			return nil, fmt.Errorf("解析插件结果失败: %w", err)
		}
	}

	// 检查插件返回的bool值，决定是否继续提交
	if pluginResult != nil {
		if boolVal, ok := pluginResult["bool"].(bool); ok && !boolVal {
			// 插件返回false，表示验证失败
			errorMsg := "插件验证失败"
			if envVal, ok := pluginResult["env"].(string); ok {
				errorMsg = envVal
			}
			return nil, fmt.Errorf(errorMsg)
		}

		// 如果插件返回了新的环境变量值，使用新值
		if newEnvVal, ok := pluginResult["env"].(string); ok {
			envValue = newEnvVal
		}
	}

	// 这里应该是实际向面板提交数据的逻辑
	submitResult := map[string]interface{}{
		"env_id":    envID,
		"env_value": envValue,
		"panel_id":  panelID,
		"status":    "success",
	}

	return submitResult, nil
}

// CreateTokenRefreshCallback 创建token刷新回调函数
func (s *PanelService) CreateTokenRefreshCallback() qinglong.TokenRefreshCallback {
	return func(panelID int64) (newToken string, err error) {
		ctx := context.Background()
		// 查询面板信息
		p, err := config.Ent.Panel.Get(ctx, panelID)
		if err != nil {
			return "", fmt.Errorf("查询面板失败: %w", err)
		}

		// 使用ClientID和ClientSecret重新获取token
		qlConfig := qinglong.NewConfig(p.URL, p.ClientID, p.ClientSecret)
		tokenResp, err := qlConfig.GetConfig()
		if err != nil {
			return "", fmt.Errorf("获取新token失败: %w", err)
		}

		if tokenResp.Code != 200 {
			return "", fmt.Errorf("获取token失败，响应码: %d, 消息: %s", tokenResp.Code, tokenResp.Message)
		}

		newToken = tokenResp.Data.Token

		// 更新数据库中的token
		err = config.Ent.Panel.UpdateOneID(panelID).SetToken(newToken).Exec(ctx)
		if err != nil {
			return "", fmt.Errorf("更新数据库token失败: %w", err)
		}

		return newToken, nil
	}
}

// CreateQlAPIWithAutoRefresh 创建带自动刷新功能的API实例
func (s *PanelService) CreateQlAPIWithAutoRefresh(panelID int64) (*qinglong.QlAPI, error) {
	ctx := context.Background()
	// 查询面板信息
	p, err := config.Ent.Panel.Query().
		Where(
			panel.IDEQ(panelID),
			panel.IsEnableEQ(true),
		).Only(ctx)
	if err != nil {
		return nil, fmt.Errorf("查询面板失败: %w", err)
	}

	// 创建带面板信息和回调函数的API实例
	callback := s.CreateTokenRefreshCallback()
	api := qinglong.NewAPIWithPanel(p.URL, p.Token, int(p.Params), panelID, callback)

	return api, nil
}
