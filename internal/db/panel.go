package db

import (
	"QLToolsV2/config"
	"QLToolsV2/internal/model"
)

type Panel struct {
	model.Panel
}

// GetPanelByID ID 获取数据
func GetPanelByID(id int) (Panel, error) {
	var m Panel
	if err := config.GinDB.Model(&m).Where("id = ?", id).
		First(&m).Error; err != nil {
		return m, err
	}
	return m, nil
}

// GetPanels 分页查询
func GetPanels(page, pageSize int) ([]Panel, error) {
	var m []Panel
	if err := config.GinDB.Scopes(PaginateIdDesc(page, pageSize)).Find(&m).Error; err != nil {
		return m, err
	}
	return m, nil
}

// GetAllPanels 获取全部数据
func GetAllPanels() (model.AllPanel, error) {
	var m Panel
	var ms model.AllPanel
	if err := config.GinDB.Model(&m).Find(&ms).Error; err != nil {
		return ms, err
	}
	return ms, nil
}

// Create 创建数据
func (m *Panel) Create() error {
	if err := config.GinDB.Create(&m).Error; err != nil {
		return err
	}
	return nil
}

// Update 修改数据
func (m *Panel) Update(data map[string]any) error {
	if err := config.GinDB.Model(&m).Where("id = ?", m.ID).
		Updates(&data).Error; err != nil {
		return err
	}
	return nil
}

// Updates 批量修改数据
func (m *Panel) Updates(ids []int, data map[string]any) error {
	if err := config.GinDB.Model(&m).Where("id IN ?", ids).
		Updates(&data).Error; err != nil {
		return err
	}
	return nil
}

// Delete 删除数据
func (m *Panel) Delete(ids []int) error {
	if err := config.GinDB.Model(&m).Delete("id IN ?", ids).Error; err != nil {
		return err
	}
	return nil
}
