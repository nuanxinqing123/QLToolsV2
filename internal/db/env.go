package db

import (
	"QLToolsV2/config"
	"QLToolsV2/internal/model"
)

type Env struct {
	model.Env
}

// GetEnvByID ID 获取数据
func GetEnvByID(id int) (Env, error) {
	var m Env
	if err := config.GinDB.Model(&m).Preload("Panel").Where("id = ?", id).
		First(&m).Error; err != nil {
		return m, err
	}
	return m, nil
}

// GetEnvs 分页查询
func GetEnvs(page, pageSize int) ([]Env, error) {
	var m []Env
	if err := config.GinDB.Scopes(Paginate(page, pageSize)).Find(&m).Error; err != nil {
		return m, err
	}
	return m, nil
}

// Create 创建数据
func (m *Env) Create() error {
	if err := config.GinDB.Create(&m).Error; err != nil {
		return err
	}
	return nil
}

// Update 修改数据
func (m *Env) Update(data map[string]any) error {
	if err := config.GinDB.Model(&m).Where("id = ?", m.ID).
		Updates(&data).Error; err != nil {
		return err
	}
	return nil
}

// Updates 批量修改数据
func (m *Env) Updates(ids []int, data map[string]any) error {
	if err := config.GinDB.Model(&m).Where("id IN ?", ids).
		Updates(&data).Error; err != nil {
		return err
	}
	return nil
}

// Delete 删除数据
func (m *Env) Delete(ids []int) error {
	if err := config.GinDB.Model(&m).Delete("id IN ?", ids).Error; err != nil {
		return err
	}
	return nil
}
