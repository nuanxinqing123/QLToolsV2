package db

import (
	"QLToolsV2/config"
	"QLToolsV2/internal/model"
)

type LoginHistory struct {
	model.LoginHistory
}

// GetLoginHistory 分页查询
func GetLoginHistory(page, pageSize int) ([]LoginHistory, int64, int64, error) {
	var m []LoginHistory
	var count int64
	if err := config.GinDB.Model(&m).Count(&count).Scopes(PaginateIdDesc(page, pageSize)).Find(&m).Error; err != nil {
		return m, count, 0, err
	}

	pn := PaginateCount(count, pageSize)

	return m, count, pn, nil
}
