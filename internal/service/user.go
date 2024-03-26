package service

import (
	"QLToolsV2/config"
	_const "QLToolsV2/const"
	"QLToolsV2/internal/db"
	res "QLToolsV2/pkg/response"
)

// Info 用户信息
func Info(userId string) (res.ResCode, any) {
	// 判断用户名是否存在
	m, err := db.GetUserByUserID(userId)
	if err != nil {
		// 记录日志
		config.GinLOG.Error(err.Error())
		return res.CodeServerBusy, _const.ServerBusy
	}

	return res.CodeSuccess, m
}
