package db

import (
	"gorm.io/gorm"

	"QLToolsV2/config"
	"QLToolsV2/internal/model"
)

type CdKey struct {
	model.CdKey
}

// GetKeyByID ID 获取数据
func GetKeyByID(id int) (CdKey, error) {
	var m CdKey
	if err := config.GinDB.Model(&m).Where("id = ?", id).
		First(&m).Error; err != nil {
		return m, err
	}
	return m, nil
}

// GetKeyByKey KEY 获取数据
func GetKeyByKey(key string) (CdKey, error) {
	var m CdKey
	if err := config.GinDB.Model(&m).Where("key = ?", key).
		First(&m).Error; err != nil {
		return m, err
	}
	return m, nil
}

// GetKeys 分页查询
func GetKeys(page, pageSize int) ([]CdKey, int64, int64, error) {
	var m []CdKey
	var count int64
	if err := config.GinDB.Model(&m).Count(&count).Scopes(PaginateIdDesc(page, pageSize)).Find(&m).Error; err != nil {
		return m, count, 0, err
	}

	pn := PaginateCount(count, pageSize)

	return m, count, pn, nil
}

// Create 创建数据
func (m *CdKey) Create() error {
	if err := config.GinDB.Create(&m).Error; err != nil {
		return err
	}
	return nil
}

// Update 修改数据
func (m *CdKey) Update(data map[string]any) error {
	if err := config.GinDB.Model(&m).Where("id = ?", m.ID).
		Updates(&data).Error; err != nil {
		return err
	}
	return nil
}

// Updates 批量修改数据
func (m *CdKey) Updates(ids []int, data map[string]any) error {
	if err := config.GinDB.Model(&m).Where("id IN ?", ids).
		Updates(&data).Error; err != nil {
		return err
	}
	return nil
}

// Delete 删除数据
func (m *CdKey) Delete(ids []int) error {
	if err := config.GinDB.Model(&m).Delete("id IN ?", ids).Error; err != nil {
		return err
	}
	return nil
}

// Deduction 扣减次数
func (m *CdKey) Deduction(count int) error {
	if err := config.GinDB.Model(&m).Where("id = ?", m.ID).
		Update("count", gorm.Expr("`count` - ?", count)).
		Error; err != nil {
		return err
	}
	return nil
}
