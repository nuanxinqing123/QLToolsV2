package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/nuanxinqing123/QLToolsV2/internal/app/config"
	"github.com/nuanxinqing123/QLToolsV2/internal/data/ent"
	"github.com/nuanxinqing123/QLToolsV2/internal/data/ent/cdkey"
	"github.com/nuanxinqing123/QLToolsV2/internal/schema"
	"github.com/segmentio/ksuid"
)

type CDKService struct{}

// NewCDKService 创建 CDKService
func NewCDKService() *CDKService {
	return &CDKService{}
}

// AddCDK 添加CDK
func (s *CDKService) AddCDK(req schema.AddCDKRequest) (*schema.AddCDKResponse, error) {
	ctx := context.Background()
	// 检查CDK密钥是否已存在
	exists, err := config.Ent.CdKey.Query().
		Where(cdkey.KeyEQ(req.Key)).
		Exist(ctx)
	if err != nil {
		return nil, fmt.Errorf("查询CDK失败: %w", err)
	}
	if exists {
		return nil, errors.New("CDK密钥已存在")
	}

	// 创建CDK记录
	cdk, err := config.Ent.CdKey.Create().
		SetKey(req.Key).
		SetCount(req.Count).
		SetIsEnable(true).
		SetCreatedAt(time.Now()).
		SetUpdatedAt(time.Now()).
		Save(ctx)
	if err != nil {
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

	ctx := context.Background()
	var builders []*ent.CdKeyCreate
	var keys []string
	now := time.Now()

	// 批量生成CDK
	for i := int32(0); i < req.Count; i++ {
		key := ksuid.New().String()

		// 检查密钥是否重复
		exists, _ := config.Ent.CdKey.Query().Where(cdkey.KeyEQ(key)).Exist(ctx)
		if exists {
			i--
			continue
		}

		builders = append(builders, config.Ent.CdKey.Create().
			SetKey(key).
			SetCount(req.UseCount).
			SetIsEnable(true).
			SetCreatedAt(now).
			SetUpdatedAt(now))
		keys = append(keys, key)
	}

	// 批量创建CDK记录
	if err := config.Ent.CdKey.CreateBulk(builders...).Exec(ctx); err != nil {
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
	ctx := context.Background()
	// 查询CDK是否存在
	cdk, err := config.Ent.CdKey.Get(ctx, req.ID)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.New("CDK不存在")
		}
		return nil, fmt.Errorf("查询CDK失败: %w", err)
	}

	// 检查密钥是否与其他CDK冲突
	if req.Key != cdk.Key {
		exists, err := config.Ent.CdKey.Query().
			Where(cdkey.And(
				cdkey.KeyEQ(req.Key),
				cdkey.IDNEQ(req.ID),
			)).
			Exist(ctx)
		if err != nil {
			return nil, fmt.Errorf("查询CDK失败: %w", err)
		}
		if exists {
			return nil, errors.New("CDK密钥已存在")
		}
	}

	// 执行更新
	updater := config.Ent.CdKey.UpdateOneID(req.ID).
		SetKey(req.Key).
		SetCount(req.Count).
		SetUpdatedAt(time.Now())

	if req.IsEnable != nil {
		updater.SetIsEnable(*req.IsEnable)
	}

	if err := updater.Exec(ctx); err != nil {
		return nil, fmt.Errorf("更新CDK失败: %w", err)
	}

	return &schema.UpdateCDKResponse{
		Message: "CDK更新成功",
	}, nil
}

// GetCDKList 获取CDK列表
func (s *CDKService) GetCDKList(req schema.GetCDKListRequest) (*schema.GetCDKListResponse, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	ctx := context.Background()
	query := config.Ent.CdKey.Query()

	if req.Key != "" {
		query.Where(cdkey.KeyContains(req.Key))
	}

	if req.IsEnable != nil {
		query.Where(cdkey.IsEnableEQ(*req.IsEnable))
	}

	total, err := query.Count(ctx)
	if err != nil {
		return nil, fmt.Errorf("查询CDK总数失败: %w", err)
	}

	offset := (req.Page - 1) * req.PageSize
	cdks, err := query.Order(ent.Desc(cdkey.FieldCreatedAt)).
		Limit(req.PageSize).
		Offset(offset).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("查询CDK列表失败: %w", err)
	}

	var list []schema.GetCDKResponse
	for _, c := range cdks {
		list = append(list, schema.GetCDKResponse{
			ID:        c.ID,
			Key:       c.Key,
			Count:     c.Count,
			IsEnable:  c.IsEnable,
			CreatedAt: c.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: c.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return &schema.GetCDKListResponse{
		Total: int64(total),
		List:  list,
	}, nil
}

// DeleteCDK 删除CDK
func (s *CDKService) DeleteCDK(req schema.DeleteCDKRequest) (*schema.DeleteCDKResponse, error) {
	ctx := context.Background()
	cdk, err := config.Ent.CdKey.Get(ctx, req.ID)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.New("CDK不存在")
		}
		return nil, fmt.Errorf("查询CDK失败: %w", err)
	}

	if err := config.Ent.CdKey.DeleteOneID(req.ID).Exec(ctx); err != nil {
		return nil, fmt.Errorf("删除CDK失败: %w", err)
	}

	return &schema.DeleteCDKResponse{
		Message: fmt.Sprintf("CDK \"%s\" 删除成功", cdk.Key),
	}, nil
}

// ToggleCDKStatus 切换CDK启用状态
func (s *CDKService) ToggleCDKStatus(req schema.ToggleCDKStatusRequest) (*schema.ToggleCDKStatusResponse, error) {
	ctx := context.Background()
	cdk, err := config.Ent.CdKey.Get(ctx, req.ID)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.New("CDK不存在")
		}
		return nil, fmt.Errorf("查询CDK失败: %w", err)
	}

	if err := config.Ent.CdKey.UpdateOneID(req.ID).
		SetIsEnable(req.IsEnable).
		SetUpdatedAt(time.Now()).
		Exec(ctx); err != nil {
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
