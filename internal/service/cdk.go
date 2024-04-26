package service

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/segmentio/ksuid"
	"gorm.io/gorm"

	"QLToolsV2/config"
	_const "QLToolsV2/const"
	"QLToolsV2/internal/db"
	"QLToolsV2/internal/model"
	res "QLToolsV2/pkg/response"
)

// KeyCheck 检查KEY
func KeyCheck(p *model.KeyCheck) (res.ResCode, any) {
	m, err := db.GetKeyByKey(p.Key)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return res.CodeGenericError, "卡密不存在"
		}
		config.GinLOG.Error(err.Error())
		return res.CodeServerBusy, _const.ServerBusy
	}
	if m.IsEnable == false {
		return res.CodeGenericError, "卡密已被禁用"
	}
	if m.Count <= 0 {
		return res.CodeGenericError, "卡密使用次数不足"
	}

	return res.CodeSuccess, m
}

// CDKList 获取KEY列表
func CDKList(p *model.Pagination) (res.ResCode, any) {
	ms, count, pn, err := db.GetKeys(p.Page, p.Size)
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

// AddCDK 添加卡密
func AddCDK(p *model.AddCDK) (res.ResCode, any) {
	m := db.CdKey{
		CdKey: model.CdKey{
			Key:      p.Key,
			Count:    p.Count,
			IsEnable: true,
		},
	}

	// 写入数据
	if err := m.Create(); err != nil {
		config.GinLOG.Error(err.Error())
		return res.CodeServerBusy, _const.ServerBusy
	}
	return res.CodeSuccess, m
}

// BatchAddCDK 批量添加卡密
func BatchAddCDK(p *model.BatchAddCDK) (res.ResCode, any) {
	// 判断本地是否还有遗留文件
	// _, err := os.Stat(_const.KeyFile)
	// if err == nil {
	// 	// 删除旧文件
	// 	err = os.Remove(_const.KeyFile)
	// 	if err != nil {
	// 		zap.L().Error(err.Error())
	// 		return res.CodeServerBusy, _const.ServerBusy
	// 	}
	// }

	// 创建KEY列表
	var keys []string

	// 获取生成数量
	for i := 0; i < p.AddCount; i++ {
		// 生成KEY
		uid := ksuid.New()

		// 写入数据
		m := db.CdKey{
			CdKey: model.CdKey{
				Key:      p.Prefix + uid.String(),
				Count:    p.Count,
				IsEnable: true,
			},
		}
		if err := m.Create(); err != nil {
			config.GinLOG.Error(err.Error())
			continue
		}

		// 加入数组
		keys = append(keys, m.Key)
	}

	// 创建文件并写入数据
	// file, err := os.OpenFile(_const.KeyFile, os.O_WRONLY|os.O_CREATE, 0)
	// if err != nil {
	// 	zap.L().Error(err.Error())
	// 	return res.CodeServerBusy, _const.ServerBusy
	// }
	// defer func(file *os.File) {
	// 	err = file.Close()
	// 	if err != nil {
	// 		config.GinLOG.Error(err.Error())
	// 	}
	// }(file)
	//
	// // 写入数据
	// writer := bufio.NewWriter(file)
	// keyString := strings.Join(keys, "\n")
	// _, err = writer.WriteString(keyString)
	// if err != nil {
	// 	config.GinLOG.Error(err.Error())
	// 	return res.CodeServerBusy, _const.ServerBusy
	// }
	//
	// // 刷新缓冲区
	// if writer.Flush() != nil {
	// 	config.GinLOG.Error(err.Error())
	// 	return res.CodeServerBusy, _const.ServerBusy
	// }

	return res.CodeSuccess, keys
}

// BatchOperationCDK 批量操作[启用/禁用]
func BatchOperationCDK(p *model.BatchOperationCDK) (res.ResCode, any) {
	// 更新数据
	m := db.CdKey{}
	if err := m.Updates(p.IDs, map[string]any{"is_enable": p.IsEnable}); err != nil {
		config.GinLOG.Error(err.Error())
		return res.CodeServerBusy, _const.ServerBusy
	}
	return res.CodeSuccess, "操作成功"
}

// UpdateCDK 修改
func UpdateCDK(p *model.UpdateCDK) (res.ResCode, any) {
	m, err := db.GetKeyByID(p.ID)
	if err != nil {
		config.GinLOG.Error(err.Error())
		return res.CodeServerBusy, _const.ServerBusy
	}

	// 更新数据
	if err = m.Update(map[string]any{
		"key":       p.Key,
		"count":     p.Count,
		"is_enable": p.IsEnable,
	}); err != nil {
		config.GinLOG.Error(err.Error())
		return res.CodeServerBusy, _const.ServerBusy
	}
	return res.CodeSuccess, "修改成功"
}

// DeleteCDK 删除
func DeleteCDK(p *model.DeleteCDK) (res.ResCode, any) {
	// 更新数据
	m := db.CdKey{}
	if err := m.Delete(p.IDs); err != nil {
		config.GinLOG.Error(err.Error())
		return res.CodeServerBusy, _const.ServerBusy
	}
	return res.CodeSuccess, "删除成功"
}
