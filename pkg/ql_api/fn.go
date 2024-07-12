package ql_api

import (
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"QLToolsV2/config"
	"QLToolsV2/internal/db"
	"QLToolsV2/internal/model"
)

type QlApiFn struct {
	db.Env
}

type ResOS struct {
	Name          string `json:"name"`           // 名称
	Remarks       string `json:"remarks"`        // 备注
	Remainder     int    `json:"remainder"`      // 剩余数量
	IsPrompt      bool   `json:"is_prompt"`      // 是否需要提示
	PromptLevel   string `json:"prompt_level"`   // 提示等级
	PromptContent string `json:"prompt_content"` // 提示内容
	EnableKey     bool   `json:"enable_key"`     // 启用CD-KEY
}

type ResQlNode struct {
	PanelId    int `json:"panel_id"`     // 面板ID
	Count      int `json:"count"`        // 可用数量
	PanelEnvId int `json:"panel_env_id"` // 面板内的变量ID位置[更新模式]
}

// GetOnlineService 获取在线服务
func GetOnlineService() ([]ResOS, error) {
	// 获取启用的所有变量以及绑定的面板
	envs, err := db.GetAllEnvs()
	if err != nil {
		config.GinLOG.Error(err.Error())
		return []ResOS{}, err
	}

	// 计算结果
	var result []ResOS

	// 开启并发处理
	var wg sync.WaitGroup
	var mu sync.Mutex

	// 计算数据
	for _, x := range envs {
		wg.Add(1)
		go func(xx db.Env) {
			defer wg.Done()

			// 计算剩余位置数
			fn := QlApiFn{
				xx,
			}
			envTotal := fn.GetEnvRemainder()

			// 开启互斥锁
			mu.Lock()
			result = append(result, ResOS{
				Name:          xx.Name,
				Remarks:       xx.Remarks,
				Remainder:     envTotal,
				IsPrompt:      xx.IsPrompt,
				PromptLevel:   xx.PromptLevel,
				PromptContent: xx.PromptContent,
				EnableKey:     xx.EnableKey,
			})
			mu.Unlock() // 解锁
		}(x)
	}

	// 等待完成
	wg.Wait()

	// 序列化缓存数据
	byteData, err := json.Marshal(result)
	if err != nil {
		config.GinLOG.Error(err.Error())
		return []ResOS{}, err
	}

	// 写入缓存
	err = config.GinCache.SetWithExpire("onlineService", string(byteData), time.Hour*24)
	if err != nil {
		config.GinLOG.Error(err.Error())
	}
	return result, nil
}

// GetEnvRemainder 计算变量剩余位置及配额
func (api *QlApiFn) GetEnvRemainder() int {
	envTotal := api.Quantity * len(api.Panels)
	config.GinLOG.Debug(fmt.Sprintf("变量总数: %d", envTotal))
	config.GinLOG.Debug(fmt.Sprintf("变量名称: %s, 变量模式: %d", api.Name, api.Mode))
	for _, p := range api.Panels {
		config.GinLOG.Debug(fmt.Sprintf("面板名称: %s", p.Name))
		// 初始化面板
		ql := &QlApi{
			URL:    p.URL,
			Token:  p.Token,
			Params: p.Params,
		}
		// 获取面板所有变量数据
		getEnvs, err := ql.GetEnvs()
		if err != nil {
			// 减去失效的面板配额
			envTotal -= api.Quantity
			config.GinLOG.Error(err.Error())
			continue
		}

		// 面板存在变量
		config.GinLOG.Debug(fmt.Sprintf("面板名称: %s, 面板存在变量: %d个", p.Name, len(getEnvs.Data)))
		if len(getEnvs.Data) > 0 {
			for _, z := range getEnvs.Data {
				config.GinLOG.Debug(fmt.Sprintf("变量名称: %s, 变量值: %s, 变量备注: %s", z.Name, z.Value, z.Remarks))
				if api.Mode == 1 || api.Mode == 3 {
					// 新建模式 || 更新模式
					if z.Name == api.Name {
						envTotal--
					}
				} else {
					// 合并模式
					if z.Name == api.Name {
						zL := strings.Split(z.Value, api.Division)
						envTotal -= len(zL)
					}
				}
			}
		}
	}
	config.GinLOG.Debug(fmt.Sprintf("可使用变量总数: %d", envTotal))
	return envTotal
}

// GetPanelByEnvMode1 新增模式
func (api *QlApiFn) GetPanelByEnvMode1() ResQlNode {
	var ps []ResQlNode

	// 开启并发处理
	var wg sync.WaitGroup
	// 互斥锁
	var mu sync.Mutex

	for _, p := range api.Panels {
		wg.Add(1)

		go func(panel model.Panel) {
			defer wg.Done()
			// 初始化面板
			ql := &QlApi{
				URL:    panel.URL,
				Token:  panel.Token,
				Params: panel.Params,
			}
			// 获取面板所有变量数据
			getEnvs, err := ql.GetEnvs()
			if err != nil {
				// 失效面板
				config.GinLOG.Error(err.Error())
				return
			}

			// 判断面板存在变量
			if len(getEnvs.Data) <= 0 {
				// 加锁
				mu.Lock()
				ps = append(ps, ResQlNode{
					PanelId:    int(panel.ID),
					Count:      api.Quantity,
					PanelEnvId: 0,
				})
				// 解锁
				mu.Unlock()
			} else {
				// 面板存在数据
				count := api.Quantity
				for _, x := range getEnvs.Data {
					if x.Name == api.Name {
						count--
					}
				}

				// 加锁
				mu.Lock()
				ps = append(ps, ResQlNode{
					PanelId:    int(panel.ID),
					Count:      count,
					PanelEnvId: 0,
				})
				// 解锁
				mu.Unlock()
			}
		}(p)
	}

	// 等待执行结束
	wg.Wait()

	// 判断是否有可用面板
	if len(ps) <= 0 {
		return ResQlNode{}
	}

	// 根据map中的count进行排序【降序】
	sort.Slice(ps, func(i, j int) bool {
		return ps[i].Count > ps[j].Count
	})
	return ps[0]
}

// GetPanelByEnvMode2 合并模式
func (api *QlApiFn) GetPanelByEnvMode2() ResQlNode {
	var ps []ResQlNode

	// 开启并发处理
	var wg sync.WaitGroup
	// 互斥锁
	var mu sync.Mutex

	for _, p := range api.Panels {
		wg.Add(1)

		go func(panel model.Panel) {
			defer wg.Done()

			// 初始化面板
			ql := &QlApi{
				URL:    panel.URL,
				Token:  panel.Token,
				Params: panel.Params,
			}
			// 获取面板所有变量数据
			getEnvs, err := ql.GetEnvs()
			if err != nil {
				// 失效面板
				config.GinLOG.Error(err.Error())
				return
			}

			// 判断面板存在变量
			if len(getEnvs.Data) <= 0 {
				// 加锁
				mu.Lock()
				ps = append(ps, ResQlNode{
					PanelId:    int(panel.ID),
					Count:      api.Quantity,
					PanelEnvId: 0,
				})
				// 解锁
				mu.Unlock()
			} else {
				// 面板存在数据
				count := api.Quantity
				for _, x := range getEnvs.Data {
					if x.Name == api.Name {
						// 根据合并分隔符分割变量值
						count -= len(strings.Split(x.Value, api.Division))
					}
				}

				// 加锁
				mu.Lock()
				ps = append(ps, ResQlNode{
					PanelId:    int(panel.ID),
					Count:      count,
					PanelEnvId: 0,
				})
				// 解锁
				mu.Unlock()
			}
		}(p)
	}

	// 等待执行结束
	wg.Wait()

	// 判断是否有可用面板
	if len(ps) <= 0 {
		return ResQlNode{}
	}

	// 根据map中的count进行排序【降序】
	sort.Slice(ps, func(i, j int) bool {
		return ps[i].Count > ps[j].Count
	})
	return ps[0]
}

// GetPanelByEnvMode3 更新模式
func (api *QlApiFn) GetPanelByEnvMode3() ResQlNode {
	/*
		1、协程查找所有绑定面板的变量存在位置
		2、如果找到, 则直接返回面板的ID, 以及变量的位置
		3、如果找不到, 则判断计算还有没有空余位置可以上传
	*/
	var ps []ResQlNode
	return ps[0]
}
