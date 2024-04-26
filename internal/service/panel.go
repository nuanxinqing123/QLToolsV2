package service

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"QLToolsV2/config"
	_const "QLToolsV2/const"
	"QLToolsV2/internal/db"
	"QLToolsV2/internal/model"
	api "QLToolsV2/pkg/ql_api"
	res "QLToolsV2/pkg/response"
)

// PanelList 获取面板列表
func PanelList(p *model.Pagination) (res.ResCode, any) {
	ms, count, pn, err := db.GetPanels(p.Page, p.Size)
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

// PanelAllList 获取全部面板列表
func PanelAllList() (res.ResCode, any) {
	ms, err := db.GetAllPanels()
	if err != nil {
		config.GinLOG.Error(err.Error())
		return res.CodeServerBusy, _const.ServerBusy
	}
	return res.CodeSuccess, ms
}

// AddPanel 添加面板
func AddPanel(p *model.AddPanel) (res.ResCode, any) {
	// 获取面板Token
	cf := api.QlConfig{
		URL:          p.URL,
		ClientID:     p.ClientID,
		ClientSecret: p.ClientSecret,
	}
	cfRes, err := cf.GetConfig()
	if err != nil {
		config.GinLOG.Error(err.Error())
		return res.CodeGenericError, fmt.Sprintf("连接面板失败, 原因: %s", err)
	}

	// 写入数据
	m := db.Panel{
		Panel: model.Panel{
			Name:         p.Name,
			URL:          p.URL,
			ClientID:     p.ClientID,
			ClientSecret: p.ClientSecret,
			IsEnable:     p.IsEnable,
			Token:        cfRes.Data.Token,
			Params:       cfRes.Data.Expiration,
		},
	}
	if err = m.Create(); err != nil {
		config.GinLOG.Error(err.Error())
		return res.CodeServerBusy, _const.ServerBusy
	}
	return res.CodeSuccess, m
}

// BatchOperationPanel 批量操作[启用/禁用]
func BatchOperationPanel(p *model.BatchOperationPanel) (res.ResCode, any) {
	// 更新数据
	m := db.Panel{}
	if err := m.Updates(p.IDs, map[string]any{"is_enable": p.IsEnable}); err != nil {
		config.GinLOG.Error(err.Error())
		return res.CodeServerBusy, _const.ServerBusy
	}
	return res.CodeSuccess, "操作成功"
}

// UpdatePanel 修改
func UpdatePanel(p *model.UpdatePanel) (res.ResCode, any) {
	// 获取面板Token
	cf := api.QlConfig{
		URL:          p.URL,
		ClientID:     p.ClientID,
		ClientSecret: p.ClientSecret,
	}
	cfRes, err := cf.GetConfig()
	if err != nil {
		config.GinLOG.Error(err.Error())
		return res.CodeGenericError, fmt.Sprintf("连接面板失败, 原因: %s", err)
	}

	m, err := db.GetPanelByID(p.ID)
	if err != nil {
		config.GinLOG.Error(err.Error())
		return res.CodeServerBusy, _const.ServerBusy
	}

	// 更新数据
	if err = m.Update(map[string]any{
		"name":          p.Name,
		"url":           p.URL,
		"client_id":     p.ClientID,
		"client_secret": p.ClientSecret,
		"is_enable":     p.IsEnable,
		"token":         cfRes.Data.Token,
		"params":        cfRes.Data.Expiration,
	}); err != nil {
		config.GinLOG.Error(err.Error())
		return res.CodeServerBusy, _const.ServerBusy
	}
	return res.CodeSuccess, "修改成功"
}

// DeletePanel 删除
func DeletePanel(p *model.DeletePanel) (res.ResCode, any) {
	// 更新数据
	m := db.Panel{}
	if err := m.Delete(p.IDs); err != nil {
		config.GinLOG.Error(err.Error())
		return res.CodeServerBusy, _const.ServerBusy
	}
	return res.CodeSuccess, "删除成功"
}
