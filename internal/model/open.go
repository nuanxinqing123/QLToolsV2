package model

type Submit struct {
	Name   string `json:"name" binding:"required"`  // 变量名
	Value  string `json:"value" binding:"required"` // 变量值
	Remark string `json:"remark"`                   // 备注
	Key    string `json:"key"`                      // KEY
}
