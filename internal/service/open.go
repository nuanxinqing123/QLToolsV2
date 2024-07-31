package service

import (
	"errors"
	"fmt"

	"github.com/bluele/gcache"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"gorm.io/gorm"

	"QLToolsV2/config"
	_const "QLToolsV2/const"
	"QLToolsV2/internal/db"
	"QLToolsV2/internal/model"
	api "QLToolsV2/pkg/ql_api"
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
			return res.CodeGenericError, "用户名或密码错误"
		}

		// 记录日志
		config.GinLOG.Error(err.Error())
		return res.CodeServerBusy, _const.ServerBusy
	}

	// 判断密码是否正确
	if !m.BcryptCheck(p.Password) {
		return res.CodeGenericError, "用户名或密码错误"
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
		config.GinLOG.Debug("开始执行实时计算")

		result, err := api.GetOnlineService()
		if err != nil {
			return res.CodeServerBusy, _const.ServerBusy
		}

		return res.CodeSuccess, result
	} else {
		config.GinLOG.Debug(fmt.Sprintf("读取缓存成功, 缓存数据为: %s", cache.(string)))

		// 反序列化缓存数据
		var result []map[string]any
		err = json.Unmarshal([]byte(cache.(string)), &result)
		if err != nil {
			config.GinLOG.Error(err.Error())
			return res.CodeServerBusy, _const.ServerBusy
		}
		// 序列化缓存数据
		return res.CodeSuccess, result
	}
}

// SubmitService 提交服务
func SubmitService(p *model.Submit) (res.ResCode, any) {
	/*
		- 判断是否为空内容
		- 检查变量名是否存在并启用
		- 检查是否启用KEY, 并且用户提交的KEY是否有效
		- 校验正则, 判断是否满足提交条件
		- 执行实时计算, 判断是否还有空余提交位置
		- 判断是否启用插件, 并且执行插件处理[未开发逻辑]
		- 提交数据, 并自动自用改变量
		- 执行结束
	*/

	// 判断是否为空内容
	if p.Value == "" {
		return res.CodeInvalidParam, "提交内容不能为空"
	}

	// 检查变量名是否存在并启用
	env, err := db.GetEnvByName(p.Name)
	if err != nil {
		config.GinLOG.Error(err.Error())
		return res.CodeGenericError, "变量不存在或未启用"
	}
	if !env.IsEnable {
		return res.CodeGenericError, "变量未启用"
	}

	// 检查是否启用KEY, 并且用户提交的KEY是否有效
	var mKey db.CdKey
	if env.EnableKey {
		mKey, err = db.GetKeyByKey(p.Key)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return res.CodeGenericError, "卡密不存在"
			}
			config.GinLOG.Error(err.Error())
			return res.CodeServerBusy, _const.ServerBusy
		}
		if mKey.IsEnable == false {
			return res.CodeGenericError, "卡密已被禁用"
		}
		if mKey.Count <= 0 {
			return res.CodeGenericError, "卡密使用次数不足"
		}
	}

	// 校验正则, 判断是否满足提交条件
	if env.Regex != "" {
		reg, err := regexp.MustCompile(env.Regex, regexp.None).MatchString(p.Value)
		if err != nil {
			config.GinLOG.Error(err.Error())
			return res.CodeServerBusy, _const.ServerBusy
		}
		if reg == false {
			return res.CodeGenericError, "提交内容不符合规则要求"
		}
	}

	// 插件处理 TODO

	// 判断执行模式
	var pd api.ResQlNode
	fn := api.QlApiFn{
		Env: env,
	}
	switch env.Mode {
	case 1:
		// 新增模式
		config.GinLOG.Debug("新增模式")
		pd = fn.GetPanelByEnvMode1()
		if pd.PanelURL == "" {
			return res.CodeGenericError, "暂无空余位置"
		}

		ql := api.InitPanel(pd.PanelURL, pd.PanelToken, pd.PanelParams)

		var pe []api.PostEnv
		pe = append(pe, api.PostEnv{
			Name:    env.Name,
			Value:   p.Value,
			Remarks: p.Remark,
		})

		resp, err := ql.PostEnvs(pe)
		if err != nil {
			config.GinLOG.Error(err.Error())
			return res.CodeServerBusy, _const.ServerBusy
		}
		config.GinLOG.Debug(fmt.Sprintf("新增模式, 返回数据为: %v", resp.Code))
	case 2:
		// 合并模式
		config.GinLOG.Debug("合并模式")
		pd = fn.GetPanelByEnvMode2()
		if pd.PanelURL == "" {
			return res.CodeGenericError, "暂无空余位置"
		}

		ql := api.InitPanel(pd.PanelURL, pd.PanelToken, pd.PanelParams)

		// 判断新建、合并
		if pd.PanelEnvValue == "" {
			// 新建
			var pe []api.PostEnv
			pe = append(pe, api.PostEnv{
				Name:    env.Name,
				Value:   p.Value,
				Remarks: p.Remark,
			})

			_, err = ql.PostEnvs(pe)
			if err != nil {
				config.GinLOG.Error(err.Error())
				return res.CodeServerBusy, _const.ServerBusy
			}
		} else {
			// 合并
			var pe api.PutEnv
			pe.Name = env.Name
			pe.Value = pd.PanelEnvValue + env.Division + p.Value
			pe.Id = pd.PanelEnvId

			_, err = ql.PutEnvs(pe)
			if err != nil {
				config.GinLOG.Error(err.Error())
				return res.CodeServerBusy, _const.ServerBusy
			}
		}
	case 3:
		// 更新模式
		config.GinLOG.Debug("更新模式")
		pd = fn.GetPanelByEnvMode3(p.Value)
		if pd.PanelURL == "" {
			return res.CodeGenericError, "暂无空余位置"
		}

		ql := api.InitPanel(pd.PanelURL, pd.PanelToken, pd.PanelParams)

		// 判断新建、合并
		if pd.PanelEnvId == 0 {
			// 新建
			var pe []api.PostEnv
			pe = append(pe, api.PostEnv{
				Name:    env.Name,
				Value:   p.Value,
				Remarks: p.Remark,
			})

			_, err = ql.PostEnvs(pe)
			if err != nil {
				config.GinLOG.Error(err.Error())
				return res.CodeServerBusy, _const.ServerBusy
			}
		} else {
			// 更新
			var pe api.PutEnv
			pe.Name = env.Name
			pe.Value = pd.PanelEnvValue + env.Division + p.Value
			pe.Remarks = p.Remark
			pe.Id = pd.PanelEnvId

			_, err = ql.PutEnvs(pe)
			if err != nil {
				config.GinLOG.Error(err.Error())
				return res.CodeServerBusy, _const.ServerBusy
			}
		}
	default:
		// 未知模式
		return res.CodeGenericError, "未知内容, 拒绝提交"
	}

	// 检查是否启用KEY, 并且扣件库存
	if env.EnableKey {
		err = mKey.Deduction(1)
		if err != nil {
			config.GinLOG.Error(err.Error())
		}
	}

	return res.CodeSuccess, gin.H{
		"msg": "提交成功",
	}
}
