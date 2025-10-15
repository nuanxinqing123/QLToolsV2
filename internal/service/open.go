package service

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"sort"
	"sync"

	"github.com/nuanxinqing123/QLToolsV2/internal/app/config"
	_const "github.com/nuanxinqing123/QLToolsV2/internal/const"
	"github.com/nuanxinqing123/QLToolsV2/internal/repository"
	"github.com/nuanxinqing123/QLToolsV2/internal/schema"
	"gorm.io/gorm"
)

type OpenService struct {
	cdkMutexMap   sync.Map       // 基于卡密值的锁映射，每个卡密有独立的锁
	pluginService *PluginService // 插件服务
	panelService  *PanelService  // 面板服务
}

// NewOpenService 创建 OpenService
func NewOpenService() *OpenService {
	return &OpenService{
		pluginService: NewPluginService(),
		panelService:  NewPanelService(),
	}
}

// getCDKMutex 获取指定卡密的互斥锁
func (s *OpenService) getCDKMutex(cdkKey string) *sync.Mutex {
	mutex, _ := s.cdkMutexMap.LoadOrStore(cdkKey, &sync.Mutex{})
	return mutex.(*sync.Mutex)
}

// CheckCDK 检查卡密
func (s *OpenService) CheckCDK(req schema.CheckCDKRequest) (*schema.CheckCDKResponse, error) {
	// 查询CDK是否存在
	cdk, err := repository.CdKeys.Where(
		repository.CdKeys.Key.Eq(req.Key),
	).Take()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &schema.CheckCDKResponse{
				Valid:         false,
				RemainingUses: 0,
				Message:       "卡密不存在",
			}, nil
		}
		return nil, fmt.Errorf("查询卡密失败: %w", err)
	}

	// 检查是否禁用
	if !cdk.IsEnable {
		return &schema.CheckCDKResponse{
			Valid:         false,
			RemainingUses: cdk.Count_,
			Message:       "卡密已被禁用",
		}, nil
	}

	// 检查使用次数是否足够
	if cdk.Count_ <= 0 {
		return &schema.CheckCDKResponse{
			Valid:         false,
			RemainingUses: 0,
			Message:       "卡密使用次数已用完",
		}, nil
	}

	return &schema.CheckCDKResponse{
		Valid:         true,
		RemainingUses: cdk.Count_,
		Message:       "卡密有效",
	}, nil
}

// GetOnlineServices 获取在线服务
func (s *OpenService) GetOnlineServices() (*schema.GetOnlineServicesResponse, error) {
	// 查询所有启用的环境变量
	query := repository.Envs.WithContext(context.Background()).Where(
		repository.Envs.IsEnable.Is(true),
	)

	// 获取总数
	total, err := query.Count()
	if err != nil {
		return nil, fmt.Errorf("查询环境变量总数失败: %w", err)
	}

	// 查询所有数据
	envs, err := query.Order(repository.Envs.CreatedAt.Desc()).Find()
	if err != nil {
		return nil, fmt.Errorf("查询环境变量列表失败: %w", err)
	}

	// 转换为响应格式并计算可用位置数
	var list []schema.OnlineServiceInfo
	for _, env := range envs {
		// 查询该环境变量绑定的启用面板ID
		var panelIDs []int64
		err := repository.EnvPanels.WithContext(context.Background()).
			Select(repository.Panels.ID).
			Join(repository.Panels, repository.EnvPanels.PanelID.EqCol(repository.Panels.ID)).
			Where(
				repository.EnvPanels.EnvID.Eq(env.ID),
				repository.Panels.IsEnable.Is(true),
			).Scan(&panelIDs)
		if err != nil {
			return nil, fmt.Errorf("查询绑定面板失败: %w", err)
		}

		// 计算可用位置数：配置变量总数 - 所有面板中该变量的实际数量
		// 注意：Quantity是所有面板共享的总配额，不是每个面板的配额
		totalSlots := env.Quantity
		usedSlots := int32(0)

		// 遍历所有绑定的面板，统计该变量在每个面板中的数量
		for _, panelID := range panelIDs {
			// 创建青龙API实例
			qlAPI, err := s.panelService.CreateQlAPIWithAutoRefresh(panelID)
			if err != nil {
				config.Log.Warn(fmt.Sprintf("创建面板%d的API实例失败: %v", panelID, err))
				continue
			}

			// 获取面板中的所有环境变量
			envResponse, err := qlAPI.GetEnvs()
			if err != nil {
				config.Log.Warn(fmt.Sprintf("获取面板%d环境变量失败: %v", panelID, err))
				continue
			}

			if envResponse.Code != 200 {
				config.Log.Warn(fmt.Sprintf("获取面板%d环境变量失败，响应码: %d", panelID, envResponse.Code))
				continue
			}

			// 统计该面板中匹配的变量数量
			for _, panelEnv := range envResponse.Data {
				if panelEnv.Name == env.Name {
					usedSlots++
				}
			}
		}

		// 计算可用位置数
		availableSlots := totalSlots - usedSlots
		if availableSlots < 0 {
			availableSlots = 0
		}

		// 调试日志：输出计算过程
		config.Log.Debug(fmt.Sprintf("环境变量[%s] ID=%d: 绑定面板数=%d, 总配额=%d, 已使用=%d, 可用=%d",
			env.Name, env.ID, len(panelIDs), totalSlots, usedSlots, availableSlots))

		list = append(list, schema.OnlineServiceInfo{
			ID:             env.ID,
			Name:           env.Name,
			Remarks:        env.Remarks,
			Quantity:       env.Quantity,
			EnableKey:      env.EnableKey,
			CdkLimit:       env.CdkLimit,
			IsPrompt:       env.IsPrompt,
			PromptLevel:    env.PromptLevel,
			PromptContent:  env.PromptContent,
			AvailableSlots: availableSlots,
		})
	}

	return &schema.GetOnlineServicesResponse{
		Total: total,
		List:  list,
	}, nil
}

// CalculateAvailableSlots 计算剩余位置
func (s *OpenService) CalculateAvailableSlots(req schema.CalculateAvailableSlotsRequest) (*schema.CalculateAvailableSlotsResponse, error) {
	// 查询环境变量是否存在且启用
	env, err := repository.Envs.Where(
		repository.Envs.ID.Eq(req.EnvID),
		repository.Envs.IsEnable.Is(true),
	).Take()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("环境变量不存在或已禁用")
		}
		return nil, fmt.Errorf("查询环境变量失败: %w", err)
	}

	// 查询该环境变量绑定的启用面板ID
	var panelIDs []int64
	err = repository.EnvPanels.WithContext(context.Background()).
		Select(repository.Panels.ID).
		Join(repository.Panels, repository.EnvPanels.PanelID.EqCol(repository.Panels.ID)).
		Where(
			repository.EnvPanels.EnvID.Eq(req.EnvID),
			repository.Panels.IsEnable.Is(true),
		).Scan(&panelIDs)
	if err != nil {
		return nil, fmt.Errorf("查询绑定面板失败: %w", err)
	}

	// 计算总位置数和已使用位置数
	// 注意：Quantity是所有面板共享的总配额，不是每个面板的配额
	totalSlots := env.Quantity
	usedSlots := int32(0)

	// 遍历所有绑定的面板，统计该变量在每个面板中的数量
	for _, panelID := range panelIDs {
		// 创建青龙API实例
		qlAPI, err := s.panelService.CreateQlAPIWithAutoRefresh(panelID)
		if err != nil {
			config.Log.Warn(fmt.Sprintf("创建面板%d的API实例失败: %v", panelID, err))
			continue
		}

		// 获取面板中的所有环境变量
		envResponse, err := qlAPI.GetEnvs()
		if err != nil {
			config.Log.Warn(fmt.Sprintf("获取面板%d环境变量失败: %v", panelID, err))
			continue
		}

		if envResponse.Code != 200 {
			config.Log.Warn(fmt.Sprintf("获取面板%d环境变量失败，响应码: %d", panelID, envResponse.Code))
			continue
		}

		// 统计该面板中匹配的变量数量
		for _, panelEnv := range envResponse.Data {
			if panelEnv.Name == env.Name {
				usedSlots++
			}
		}
	}

	// 计算可用位置数
	availableSlots := totalSlots - usedSlots
	if availableSlots < 0 {
		availableSlots = 0
	}

	return &schema.CalculateAvailableSlotsResponse{
		EnvID:          req.EnvID,
		TotalSlots:     totalSlots,
		UsedSlots:      usedSlots,
		AvailableSlots: availableSlots,
	}, nil
}

// SubmitVariable 提交变量
func (s *OpenService) SubmitVariable(req schema.SubmitVariableRequest) (*schema.SubmitVariableResponse, error) {
	// 1. 判断是否为空内容
	if req.Value == "" {
		return &schema.SubmitVariableResponse{
			Success: false,
			Message: "变量值不能为空",
		}, nil
	}

	// 2. 检查变量名是否存在并启用
	env, err := repository.Envs.Where(
		repository.Envs.ID.Eq(req.EnvID),
		repository.Envs.IsEnable.Is(true),
	).Take()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &schema.SubmitVariableResponse{
				Success: false,
				Message: "环境变量不存在或已禁用",
			}, nil
		}
		return nil, fmt.Errorf("查询环境变量失败: %w", err)
	}

	var remainingCDK int32 = 0

	// 3. 检查是否启用KEY，并且用户提交的KEY是否有效
	if env.EnableKey {
		if req.Key == "" {
			return &schema.SubmitVariableResponse{
				Success: false,
				Message: "该服务需要提供有效的卡密",
			}, nil
		}

		// 使用基于卡密值的互斥锁防止并发问题
		cdkMutex := s.getCDKMutex(req.Key)
		cdkMutex.Lock()
		defer cdkMutex.Unlock()

		// 检查卡密
		cdkResp, err := s.CheckCDK(schema.CheckCDKRequest{Key: req.Key})
		if err != nil {
			return nil, fmt.Errorf("检查卡密失败: %w", err)
		}

		if !cdkResp.Valid {
			return &schema.SubmitVariableResponse{
				Success: false,
				Message: cdkResp.Message,
			}, nil
		}

		// 检查卡密次数是否足够
		if cdkResp.RemainingUses < env.CdkLimit {
			return &schema.SubmitVariableResponse{
				Success: false,
				Message: fmt.Sprintf("卡密剩余次数不足，需要%d次，剩余%d次", env.CdkLimit, cdkResp.RemainingUses),
			}, nil
		}

		remainingCDK = cdkResp.RemainingUses
	}

	// 4. 校验正则，判断是否满足提交条件，并提取匹配内容
	if env.Regex != nil && *env.Regex != "" {
		re, err := regexp.Compile(*env.Regex)
		if err != nil {
			return nil, fmt.Errorf("正则表达式错误: %w", err)
		}

		// 查找匹配的内容
		matched := re.FindString(req.Value)
		if matched == "" {
			return &schema.SubmitVariableResponse{
				Success: false,
				Message: "变量值格式不符合要求",
			}, nil
		}

		// 将匹配到的内容替换原始值
		req.Value = matched
	}

	// 5. 执行实时计算，判断是否还有空余提交位置
	slotsResp, err := s.CalculateAvailableSlots(schema.CalculateAvailableSlotsRequest{EnvID: req.EnvID})
	if err != nil {
		return nil, fmt.Errorf("计算可用位置失败: %w", err)
	}

	if slotsResp.AvailableSlots <= 0 {
		return &schema.SubmitVariableResponse{
			Success: false,
			Message: "当前服务已满，暂无可用位置",
		}, nil
	}

	// 6. todo 判断是否启用插件，并且执行插件处理
	processedValue := req.Value
	pluginResult, err := s.pluginService.ExecutePluginsForEnv(req.EnvID, req.Value)
	if err != nil {
		return nil, fmt.Errorf("执行插件处理失败: %w", err)
	}
	if pluginResult != nil && pluginResult.Success && len(pluginResult.OutputData) > 0 {
		// 如果插件处理成功且有输出数据，使用处理后的数据
		processedValue = string(pluginResult.OutputData)
	}

	// 7. 提交数据到所有绑定的面板，并根据IsAutoEnvEnable判断是否需要启用提交变量
	// 查询该环境变量绑定的启用面板
	var panelIDs []int64
	err = repository.EnvPanels.WithContext(context.Background()).
		Select(repository.Panels.ID).
		Join(repository.Panels, repository.EnvPanels.PanelID.EqCol(repository.Panels.ID)).
		Where(
			repository.EnvPanels.EnvID.Eq(req.EnvID),
			repository.Panels.IsEnable.Is(true),
		).Scan(&panelIDs)
	if err != nil {
		return nil, fmt.Errorf("查询绑定面板失败: %w", err)
	}

	// 根据模式选择提交策略
	submittedTo := int32(0)

	if env.Mode == _const.CreateMode {
		// 新建模式：使用负载均衡，选择可用位置最多的面板
		err = s.submitAndAutoEnable(req.EnvID, panelIDs, env.Name, processedValue, req.Remarks, env.IsAutoEnvEnable)
		if err != nil {
			return nil, err
		}
		submittedTo = 1

	} else if env.Mode == _const.UpdateMode {
		// 更新模式：遍历所有面板，根据正则表达式匹配并更新
		if env.RegexUpdate == nil || *env.RegexUpdate == "" {
			return nil, errors.New("更新模式下必须设置更新正则表达式")
		}

		updatedCount, _, err := s.updateExistingVariables(panelIDs, env.Name, *env.RegexUpdate, processedValue, req.Remarks)
		if err != nil {
			return nil, fmt.Errorf("更新现有变量失败: %w", err)
		}

		if updatedCount == 0 {
			// 没有匹配到任何变量，使用新建逻辑
			config.Log.Info("更新模式下未匹配到任何变量，使用新建逻辑")
			err = s.submitAndAutoEnable(req.EnvID, panelIDs, env.Name, processedValue, req.Remarks, env.IsAutoEnvEnable)
			if err != nil {
				return nil, err
			}
			submittedTo = 1
		} else {
			submittedTo = int32(updatedCount)
		}

	} else {
		return nil, fmt.Errorf("不支持的模式: %d", env.Mode)
	}

	// 如果启用了KEY验证，扣减卡密次数
	if env.EnableKey {
		// 扣减卡密次数
		_, err = repository.CdKeys.Where(
			repository.CdKeys.Key.Eq(req.Key),
		).UpdateSimple(repository.CdKeys.Count_.Sub(env.CdkLimit))
		if err != nil {
			return nil, fmt.Errorf("扣减卡密次数失败: %w", err)
		}
		remainingCDK -= env.CdkLimit
	}

	return &schema.SubmitVariableResponse{
		Success:      true,
		Message:      "变量提交成功",
		SubmittedTo:  submittedTo,
		RemainingCDK: remainingCDK,
	}, nil
}

// PanelLoadInfo 面板负载信息
type PanelLoadInfo struct {
	PanelID   int64
	UsedSlots int32
}

// selectBestPanelForSubmit 选择最佳面板进行提交（负载均衡）
// 策略：选择该环境变量数量最少的面板，即负载最轻的面板
func (s *OpenService) selectBestPanelForSubmit(envID int64, panelIDs []int64) (int64, error) {
	// 获取环境变量信息
	env, err := repository.Envs.Where(repository.Envs.ID.Eq(envID)).Take()
	if err != nil {
		return 0, fmt.Errorf("查询环境变量失败: %w", err)
	}

	var panelLoads []PanelLoadInfo

	// 获取每个面板的负载情况
	for _, panelID := range panelIDs {
		// 创建青龙API实例
		qlAPI, err := s.panelService.CreateQlAPIWithAutoRefresh(panelID)
		if err != nil {
			config.Log.Warn(fmt.Sprintf("创建面板%d的API实例失败: %v", panelID, err))
			continue
		}

		// 获取面板中的所有环境变量
		envResponse, err := qlAPI.GetEnvs()
		if err != nil {
			config.Log.Warn(fmt.Sprintf("获取面板%d环境变量失败: %v", panelID, err))
			continue
		}

		if envResponse.Code != 200 {
			config.Log.Warn(fmt.Sprintf("获取面板%d环境变量失败，响应码: %d", panelID, envResponse.Code))
			continue
		}

		// 统计该面板中匹配的变量数量
		envCount := int32(0)
		for _, panelEnv := range envResponse.Data {
			if panelEnv.Name == env.Name {
				envCount++
			}
		}

		panelLoads = append(panelLoads, PanelLoadInfo{
			PanelID:   panelID,
			UsedSlots: envCount,
		})
	}

	if len(panelLoads) == 0 {
		return 0, errors.New("没有可用的面板")
	}

	// 按已使用位置数升序排序，选择负载最轻的面板
	sort.Slice(panelLoads, func(i, j int) bool {
		return panelLoads[i].UsedSlots < panelLoads[j].UsedSlots
	})

	bestPanel := panelLoads[0]
	config.Log.Info(fmt.Sprintf("选择面板%d进行提交，当前该变量数量: %d",
		bestPanel.PanelID, bestPanel.UsedSlots))

	return bestPanel.PanelID, nil
}

// submitAndAutoEnable 提交变量到最佳面板并根据配置自动启用
func (s *OpenService) submitAndAutoEnable(envID int64, panelIDs []int64, envName, processedValue, remarks string, isAutoEnable bool) error {
	// 选择最佳面板
	bestPanelID, err := s.selectBestPanelForSubmit(envID, panelIDs)
	if err != nil {
		return fmt.Errorf("选择最佳面板失败: %w", err)
	}

	// 提交到最佳面板
	panelEnvID, err := s.submitToPanel(bestPanelID, envName, processedValue, remarks)
	if err != nil {
		return fmt.Errorf("提交到面板%d失败: %w", bestPanelID, err)
	}

	// 如果需要自动启用，则启用环境变量
	if isAutoEnable && panelEnvID > 0 {
		err = s.enablePanelEnv(bestPanelID, panelEnvID)
		if err != nil {
			config.Log.Warn(fmt.Sprintf("自动启用面板%d变量%d失败: %v", bestPanelID, panelEnvID, err))
		}
	}

	return nil
}

// submitToPanel 提交变量到指定面板
func (s *OpenService) submitToPanel(panelID int64, name, value, remarks string) (int, error) {

	// 创建青龙API实例（使用自动token刷新功能）
	qlAPI, err := s.panelService.CreateQlAPIWithAutoRefresh(panelID)
	if err != nil {
		return 0, fmt.Errorf("创建青龙API实例失败: %w", err)
	}

	// 构建环境变量数据
	envData := []schema.PostEnvRequest{{
		Name:    name,
		Value:   value,
		Remarks: remarks,
	}}

	// 提交环境变量
	response, err := qlAPI.PostEnvs(envData)
	if err != nil {
		return 0, fmt.Errorf("提交环境变量失败: %w", err)
	}

	// 检查响应状态
	if response.Code != 200 {
		return 0, fmt.Errorf("提交失败，响应码: %d", response.Code)
	}

	// 获取创建的变量ID
	if len(response.Data) > 0 {
		createdEnvID := response.Data[0].Id
		config.Log.Info(fmt.Sprintf("成功提交变量到面板%d，变量ID: %d", panelID, createdEnvID))
		return createdEnvID, nil
	}

	return 0, errors.New("未获取到创建的变量ID")
}

// enablePanelEnv 启用面板中的环境变量
func (s *OpenService) enablePanelEnv(panelID int64, envID int) error {
	// 创建青龙API实例
	qlAPI, err := s.panelService.CreateQlAPIWithAutoRefresh(panelID)
	if err != nil {
		return fmt.Errorf("创建青龙API实例失败: %w", err)
	}

	// 启用环境变量
	enableRequest := schema.PutEnableEnvRequest{envID}
	response, err := qlAPI.PutEnableEnvs(enableRequest)
	if err != nil {
		return fmt.Errorf("启用环境变量失败: %w", err)
	}

	if response.Code != 200 {
		return fmt.Errorf("启用失败，响应码: %d", response.Code)
	}

	config.Log.Info(fmt.Sprintf("成功启用面板%d中的变量%d", panelID, envID))
	return nil
}

// updateExistingVariables 更新现有变量（更新模式）
func (s *OpenService) updateExistingVariables(panelIDs []int64, envName, regexPattern, newValue, remarks string) (int, []int64, error) {
	// 编译正则表达式
	regex, err := regexp.Compile(regexPattern)
	if err != nil {
		return 0, nil, fmt.Errorf("编译正则表达式失败: %w", err)
	}

	updatedCount := 0
	var updatedPanelIDs []int64

	// 遍历所有面板
	for _, panelID := range panelIDs {
		// 创建青龙API实例
		qlAPI, err := s.panelService.CreateQlAPIWithAutoRefresh(panelID)
		if err != nil {
			config.Log.Warn(fmt.Sprintf("创建面板%d的API实例失败: %v", panelID, err))
			continue
		}

		// 获取面板中的所有环境变量
		envResponse, err := qlAPI.GetEnvs()
		if err != nil {
			config.Log.Warn(fmt.Sprintf("获取面板%d环境变量失败: %v", panelID, err))
			continue
		}

		if envResponse.Code != 200 {
			config.Log.Warn(fmt.Sprintf("获取面板%d环境变量失败，响应码: %d", panelID, envResponse.Code))
			continue
		}

		// 查找匹配的环境变量
		panelUpdated := false
		for _, env := range envResponse.Data {
			// 检查变量名是否匹配
			if env.Name == envName {
				// 检查变量值是否匹配正则表达式
				if regex.MatchString(env.Value) {
					// 更新变量
					updateRequest := schema.PutEnvRequest{
						Id:      env.Id,
						Name:    env.Name,
						Value:   newValue,
						Remarks: remarks,
					}

					updateResponse, err := qlAPI.PutEnvs(updateRequest)
					if err != nil {
						config.Log.Warn(fmt.Sprintf("更新面板%d变量%d失败: %v", panelID, env.Id, err))
						continue
					}

					if updateResponse.Code != 200 {
						config.Log.Warn(fmt.Sprintf("更新面板%d变量%d失败，响应码: %d", panelID, env.Id, updateResponse.Code))
						continue
					}

					config.Log.Info(fmt.Sprintf("成功更新面板%d变量%d: %s", panelID, env.Id, env.Name))
					panelUpdated = true
				}
			}
		}

		if panelUpdated {
			updatedCount++
			updatedPanelIDs = append(updatedPanelIDs, panelID)
		}
	}

	return updatedCount, updatedPanelIDs, nil
}
