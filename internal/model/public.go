package model

// Pagination 分页
type Pagination struct {
	Page int `form:"page,default=1"`
	Size int `form:"size,default=20"`
}
