package service

import (
	"github.com/gin-gonic/gin"

	"QLToolsV2/config"
	_const "QLToolsV2/const"
	"QLToolsV2/internal/db"
	"QLToolsV2/internal/model"
	res "QLToolsV2/pkg/response"
)

// // Info 用户信息
// func Info(userId string) (res.ResCode, any) {
// 	// 判断用户名是否存在
// 	m, err := db.GetUserByUserID(userId)
// 	if err != nil {
// 		// 记录日志
// 		config.GinLOG.Error(err.Error())
// 		return res.CodeServerBusy, _const.ServerBusy
// 	}
//
// 	return res.CodeSuccess, m
// }

// LoginInfo 获取登录信息
func LoginInfo(p *model.Pagination) (res.ResCode, any) {
	ms, count, pn, err := db.GetLoginHistory(p.Page, p.Size)
	if err != nil {
		config.GinLOG.Error(err.Error())
		return res.CodeServerBusy, _const.ServerBusy
	}
	return res.CodeSuccess, gin.H{
		"data":   ms,
		"totals": count,
		"pages":  pn,
	}
}
