package db

import (
	"gorm.io/gorm"
)

// PaginateIdDesc 根据 ID 倒序分页查询
func PaginateIdDesc(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 20
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize).Order("id desc")
	}
}

// PaginateCount 计算页码数
func PaginateCount(count int64, pageSize int) int64 {
	switch {
	case pageSize > 100:
		pageSize = 100
	case pageSize <= 0:
		pageSize = 20
	}

	// 计算页码数(pn:Page Number)
	pn := count / int64(pageSize)
	if count%int64(pageSize) != 0 {
		pn++
	}

	return pn
}
