package service

import (
	"errors"

	"github.com/bluele/gcache"
	jsoniter "github.com/json-iterator/go"
	"gorm.io/gorm"

	"QLToolsV2/config"
	_const "QLToolsV2/const"
	"QLToolsV2/internal/db"
	"QLToolsV2/internal/model"
	res "QLToolsV2/pkg/response"
	"QLToolsV2/utils"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

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

	// 初始化 JWT
	j := utils.NewJWT()

	// 已存在用户, 生成授权 Token
	claims := j.CreateClaims(utils.BaseClaims{
		UserID: m.UserID,
	})

	token, err := j.CreateToken(claims)
	if err != nil {
		config.GinLOG.Error("[生成 Token]失败，原因：" + err.Error())
		return res.CodeServerBusy, "系统繁忙，请稍候再试"
	}

	return res.CodeSuccess, token
}

// Register 用户注册
func Register(p *model.Register) (res.ResCode, any) {
	var userCount int64

	if err := config.GinDB.Model(&db.User{}).Count(&userCount).Error; err != nil {
		config.GinLOG.Error(err.Error())
		return res.CodeServerBusy, _const.ServerBusy
	}
	if userCount > 0 {
		return res.CodeGenericError, "管理员已存在, 已自动关闭注册功能"
	}

	// 创建用户
	m := db.User{
		User: model.User{
			UserID:   utils.GenID(),
			UserName: p.UserName,
		},
	}
	// 处理密码【根据自己对密码强度的需求进行修改】
	m.BcryptHash(p.Password)

	// 写入数据
	if err := m.Create(); err != nil {
		config.GinLOG.Error(err.Error())
		return res.CodeServerBusy, _const.ServerBusy
	}

	// 初始化 JWT
	j := utils.NewJWT()

	// 已存在用户, 生成授权 Token
	claims := j.CreateClaims(utils.BaseClaims{
		UserID: m.UserID,
	})

	token, err := j.CreateToken(claims)
	if err != nil {
		config.GinLOG.Error("[生成 Token]失败，原因：" + err.Error())
		return res.CodeServerBusy, "系统繁忙，请稍候再试"
	}

	return res.CodeSuccess, token
}

// OnlineService 在线服务
func OnlineService() (res.ResCode, any) {
	// 查询缓存是否存在
	cache, err := config.GinCache.GetIFPresent("onlineService")
	if err != nil {
		// 如果不是缓存不存在的错误
		if !errors.Is(gcache.KeyNotFoundError, err) {
			config.GinLOG.Error(err.Error())
			return res.CodeServerBusy, _const.ServerBusy
		}

		/*
			执行实时计算
		*/
		// 获取启用的所有变量以及绑定的面板
		envs, err := db.GetAllEnvs()
		if err != nil {
			config.GinLOG.Error(err.Error())
			return res.CodeServerBusy, _const.ServerBusy
		}
		// for _, x := range envs {
		// 	// 变量总数
		// 	envTotal := x.Quantity * len(x.Panels)
		// 	// 变量剩余总数
		// 	envTotalRemaining := 0
		// 	for _, y := range x.Panels {
		// 		// 获取面板所有变量数据
		// 	}
		//
		// }

		return res.CodeSuccess, envs
	} else {
		// 序列化缓存数据
		return res.CodeSuccess, cache
	}
}
