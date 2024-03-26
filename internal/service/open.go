package service

import (
	"errors"

	"gorm.io/gorm"

	"QLToolsV2/config"
	_const "QLToolsV2/const"
	"QLToolsV2/internal/db"
	"QLToolsV2/internal/model"
	res "QLToolsV2/pkg/response"
	"QLToolsV2/utils"
)

// var json = jsoniter.ConfigCompatibleWithStandardLibrary

// Login 用户登录
func Login(p *model.Login) (res.ResCode, any) {
	// 判断用户名是否存在
	m, err := db.GetUserByUsername(p.UserName)
	if err != nil {
		// 判断是否注册
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return res.CodeGenericError, "用户名不存在, 请先注册"
		}

		// 记录日志
		config.GinLOG.Error(err.Error())
		return res.CodeServerBusy, _const.ServerBusy
	}

	// 判断密码是否正确
	if !m.BcryptCheck(p.Password) {
		return res.CodeGenericError, "密码错误"
	}

	// 判断是否已封禁
	if m.Status == model.Inactive {
		return res.CodeGenericError, "该账号还未激活"
	} else if m.Status == model.Suspend {
		return res.CodeGenericError, "该账号已被暂停"
	}

	return res.CodeSuccess, ""
}

// Register 用户注册
func Register(p *model.Register) (res.ResCode, any) {
	var userCount int64

	if err := config.GinDB.Model(&db.User{}).Where("username = ?", p.UserName).
		Count(&userCount).Error; err != nil {
		config.GinLOG.Error(err.Error())
		return res.CodeServerBusy, _const.ServerBusy
	}
	if userCount > 0 {
		return res.CodeGenericError, "用户名已存在"
	}

	// 创建用户
	user := db.User{
		User: model.User{
			UserID:   utils.GenID(),
			UserName: p.UserName,
			Avatar:   "也许大概没有头像吧",
			Balance:  0,
			Role:     model.UserRole,
			Status:   model.Active,
		},
	}
	// 处理密码【根据自己对密码强度的需求进行修改】
	user.BcryptHash(p.Password)

	if err := config.GinDB.Create(&user).Error; err != nil {
		config.GinLOG.Error(err.Error())
		return res.CodeServerBusy, _const.ServerBusy
	}

	return res.CodeSuccess, ""
}
