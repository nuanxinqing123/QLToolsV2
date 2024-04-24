package service

import (
	"gorm.io/gorm"

	"QLToolsV2/config"
	_const "QLToolsV2/const"
	"QLToolsV2/internal/db"
	"QLToolsV2/internal/model"
	res "QLToolsV2/pkg/response"
)

// EnvList 获取变量列表
func EnvList(p *model.Pagination) (res.ResCode, any) {
	ms, err := db.GetEnvs(p.Page, p.Size)
	if err != nil {
		config.GinLOG.Error(err.Error())
		return res.CodeServerBusy, _const.ServerBusy
	}
	return res.CodeSuccess, ms
}

// AddEnv 添加变量
func AddEnv(p *model.AddEnv) (res.ResCode, any) {
	m := db.Env{
		Env: model.Env{
			Name:        p.Name,
			Remarks:     p.Remarks,
			Quantity:    p.Quantity,
			Regex:       p.Regex,
			Mode:        p.Mode,
			Division:    p.Division,
			RegexUpdate: p.RegexUpdate,
			EnableKey:   p.EnableKey,
			IsEnable:    p.IsEnable,
		},
	}

	// 写入数据
	if err := m.Create(); err != nil {
		config.GinLOG.Error(err.Error())
		return res.CodeServerBusy, _const.ServerBusy
	}
	return res.CodeSuccess, m
}

// BatchOperationEnv 批量操作[启用/禁用]
func BatchOperationEnv(p *model.BatchOperationEnv) (res.ResCode, any) {
	// 更新数据
	m := db.Env{}
	if err := m.Updates(p.IDs, map[string]any{"is_enable": p.IsEnable}); err != nil {
		config.GinLOG.Error(err.Error())
		return res.CodeServerBusy, _const.ServerBusy
	}
	return res.CodeSuccess, "操作成功"
}

// UpdateEnv 修改
func UpdateEnv(p *model.UpdateEnv) (res.ResCode, any) {
	m, err := db.GetEnvByID(p.ID)
	if err != nil {
		config.GinLOG.Error(err.Error())
		return res.CodeServerBusy, _const.ServerBusy
	}

	// 更新数据
	if err = m.Update(map[string]any{
		"name":         p.Name,
		"remarks":      p.Remarks,
		"quantity":     p.Quantity,
		"regex":        p.Regex,
		"mode":         p.Mode,
		"division":     p.Division,
		"regex_update": p.RegexUpdate,
		"enable_key":   p.EnableKey,
		"is_enable":    p.IsEnable,
	}); err != nil {
		config.GinLOG.Error(err.Error())
		return res.CodeServerBusy, _const.ServerBusy
	}
	return res.CodeSuccess, "修改成功"
}

// BindPanel 绑定面板
func BindPanel(p *model.BindPanel) (res.ResCode, any) {
	m, err := db.GetEnvByID(p.EnvID)
	if err != nil {
		config.GinLOG.Error(err.Error())
		return res.CodeServerBusy, _const.ServerBusy
	}

	for _, id := range p.PanelIDs {
		m.Panels = append(m.Panels, model.Panel{
			Model: gorm.Model{
				ID: uint(id),
			},
		})
	}

	// 更新数据
	if err = m.Save(); err != nil {
		config.GinLOG.Error(err.Error())
		return res.CodeServerBusy, _const.ServerBusy
	}

	return res.CodeSuccess, "绑定成功"
}

// DeleteEnv 删除
func DeleteEnv(p *model.DeleteEnv) (res.ResCode, any) {
	// 更新数据
	m := db.Env{}
	if err := m.Delete(p.IDs); err != nil {
		config.GinLOG.Error(err.Error())
		return res.CodeServerBusy, _const.ServerBusy
	}
	return res.CodeSuccess, "删除成功"
}
