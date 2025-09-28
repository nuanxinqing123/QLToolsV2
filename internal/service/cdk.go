package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/nuanxinqing123/QLToolsV2/internal/model"
	"github.com/nuanxinqing123/QLToolsV2/internal/repository"
	"github.com/nuanxinqing123/QLToolsV2/internal/schema"
	"github.com/segmentio/ksuid"
	"gorm.io/gorm"
)

type CDKService struct{}

// NewCDKService 创建 CDKService
func NewCDKService() *CDKService {
	return &CDKService{}
}

// AddCDK 添加CDK
func (s *CDKService) AddCDK(req schema.AddCDKRequest) (*schema.AddCDKResponse, error) {
	// 检查CDK密钥是否已存在
	existingCDK, err := repository.CdKeys.Where(
		repository.CdKeys.Key.Eq(req.Key),
	).Take()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("查询CDK失败: %w", err)
	}
	if existingCDK != nil {
		return nil, errors.New("CDK密钥已存在")
	}

	now := time.Now()
	cdk := &model.CdKeys{
		CreatedAt: now,
		UpdatedAt: now,
		Key:       req.Key,
		Count_:    req.Count,
		IsEnable:  true, // 默认启用
	}

	// 创建CDK记录
	if err = repository.CdKeys.WithContext(context.Background()).Create(cdk); err != nil {
		return nil, fmt.Errorf("创建CDK失败: %w", err)
	}

	return &schema.AddCDKResponse{
		ID:      cdk.ID,
		Message: "CDK添加成功",
	}, nil
}

// AddCDKBatch 批量添加CDK
func (s *CDKService) AddCDKBatch(req schema.AddCDKBatchRequest) (*schema.AddCDKBatchResponse, error) {
	if req.Count <= 0 || req.Count > 1000 {
		return nil, errors.New("生成数量必须在1-1000之间")
	}

	var cdks []*model.CdKeys
	var keys []string
	now := time.Now()

	// 批量生成CDK
	for i := int32(0); i < req.Count; i++ {
		// 使用ksuid生成唯一CDK密钥
		key := ksuid.New().String()

		// 检查密钥是否重复（ksuid几乎不会重复，但为了保险起见）
		existingCDK, err := repository.CdKeys.Where(
			repository.CdKeys.Key.Eq(key),
		).Take()
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("查询CDK失败: %w", err)
		}
		if existingCDK != nil {
			// 重新生成（极小概率）
			i--
			continue
		}

		cdk := &model.CdKeys{
			CreatedAt: now,
			UpdatedAt: now,
			Key:       key,
			Count_:    req.UseCount,
			IsEnable:  true,
		}

		cdks = append(cdks, cdk)
		keys = append(keys, key)
	}

	// 批量创建CDK记录
	if err := repository.CdKeys.WithContext(context.Background()).CreateInBatches(cdks, 100); err != nil {
		return nil, fmt.Errorf("批量创建CDK失败: %w", err)
	}

	return &schema.AddCDKBatchResponse{
		Count:   int32(len(keys)),
		Keys:    keys,
		Message: fmt.Sprintf("成功生成%d个CDK", len(keys)),
	}, nil
}

// UpdateCDK 更新CDK
func (s *CDKService) UpdateCDK(req schema.UpdateCDKRequest) (*schema.UpdateCDKResponse, error) {
	// 查询CDK是否存在
	cdk, err := repository.CdKeys.Where(
		repository.CdKeys.ID.Eq(req.ID),
	).Take()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("CDK不存在")
		}
		return nil, fmt.Errorf("查询CDK失败: %w", err)
	}

	// 检查密钥是否与其他CDK冲突
	if req.Key != cdk.Key {
		existingCDK, err := repository.CdKeys.Where(
			repository.CdKeys.Key.Eq(req.Key),
			repository.CdKeys.ID.Neq(req.ID),
		).Take()
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("查询CDK失败: %w", err)
		}
		if existingCDK != nil {
			return nil, errors.New("CDK密钥已存在")
		}
	}

	// 构建更新数据
	updates := map[string]interface{}{
		"key":        req.Key,
		"count":      req.Count,
		"updated_at": time.Now(),
	}

	// 如果指定了启用状态，则更新
	if req.IsEnable != nil {
		updates["is_enable"] = *req.IsEnable
	}

	// 执行更新
	if _, err = repository.CdKeys.Where(
		repository.CdKeys.ID.Eq(req.ID),
	).Updates(updates); err != nil {
		return nil, fmt.Errorf("更新CDK失败: %w", err)
	}

	return &schema.UpdateCDKResponse{
		Message: "CDK更新成功",
	}, nil
}

// GetCDKList 获取CDK列表
func (s *CDKService) GetCDKList(req schema.GetCDKListRequest) (*schema.GetCDKListResponse, error) {
	// 设置默认分页参数
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	// 构建查询条件
	query := repository.CdKeys.WithContext(context.Background())

	// 按密钥模糊搜索
	if req.Key != "" {
		query = query.Where(repository.CdKeys.Key.Like("%" + req.Key + "%"))
	}

	// 按启用状态筛选
	if req.IsEnable != nil {
		query = query.Where(repository.CdKeys.IsEnable.Is(*req.IsEnable))
	}

	// 获取总数
	total, err := query.Count()
	if err != nil {
		return nil, fmt.Errorf("查询CDK总数失败: %w", err)
	}

	// 分页查询
	offset := (req.Page - 1) * req.PageSize
	cdks, err := query.Order(repository.CdKeys.CreatedAt.Desc()).
		Limit(req.PageSize).
		Offset(offset).
		Find()
	if err != nil {
		return nil, fmt.Errorf("查询CDK列表失败: %w", err)
	}

	// 转换为响应格式
	var list []schema.GetCDKResponse
	for _, cdk := range cdks {
		list = append(list, schema.GetCDKResponse{
			ID:        cdk.ID,
			Key:       cdk.Key,
			Count:     cdk.Count_,
			IsEnable:  cdk.IsEnable,
			CreatedAt: cdk.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: cdk.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return &schema.GetCDKListResponse{
		Total: total,
		List:  list,
	}, nil
}

// DeleteCDK 删除CDK
func (s *CDKService) DeleteCDK(req schema.DeleteCDKRequest) (*schema.DeleteCDKResponse, error) {
	// 查询CDK是否存在
	cdk, err := repository.CdKeys.Where(
		repository.CdKeys.ID.Eq(req.ID),
	).Take()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("CDK不存在")
		}
		return nil, fmt.Errorf("查询CDK失败: %w", err)
	}

	// 执行软删除
	if _, err = repository.CdKeys.Where(
		repository.CdKeys.ID.Eq(req.ID),
	).Delete(); err != nil {
		return nil, fmt.Errorf("删除CDK失败: %w", err)
	}

	return &schema.DeleteCDKResponse{
		Message: fmt.Sprintf("CDK \"%s\" 删除成功", cdk.Key),
	}, nil
}

// ToggleCDKStatus 切换CDK启用状态
func (s *CDKService) ToggleCDKStatus(req schema.ToggleCDKStatusRequest) (*schema.ToggleCDKStatusResponse, error) {
	// 查询CDK是否存在
	cdk, err := repository.CdKeys.Where(
		repository.CdKeys.ID.Eq(req.ID),
	).Take()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("CDK不存在")
		}
		return nil, fmt.Errorf("查询CDK失败: %w", err)
	}

	// 更新启用状态
	if _, err = repository.CdKeys.Where(
		repository.CdKeys.ID.Eq(req.ID),
	).Updates(map[string]interface{}{
		"is_enable":  req.IsEnable,
		"updated_at": time.Now(),
	}); err != nil {
		return nil, fmt.Errorf("更新CDK状态失败: %w", err)
	}

	status := "禁用"
	if req.IsEnable {
		status = "启用"
	}

	return &schema.ToggleCDKStatusResponse{
		Message: fmt.Sprintf("CDK \"%s\" %s成功", cdk.Key, status),
	}, nil
}
