package ql_api

import (
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"QLToolsV2/config"
	"QLToolsV2/internal/db"
)

type QlApiFn struct {
	db.Env
}

// GetOnlineService 获取在线服务
func GetOnlineService() (any, error) {
	// 获取启用的所有变量以及绑定的面板
	envs, err := db.GetAllEnvs()
	if err != nil {
		config.GinLOG.Error(err.Error())
		return nil, err
	}

	// 计算结果
	var result []map[string]any

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
			result = append(result, map[string]any{
				"name":           xx.Name,
				"remarks":        xx.Remarks,
				"remainder":      envTotal,
				"is_prompt":      xx.IsPrompt,
				"prompt_level":   xx.PromptLevel,
				"prompt_content": xx.PromptContent,
				"enable_key":     xx.EnableKey,
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
		return nil, err
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
func (api *QlApiFn) GetPanelByEnvMode1() map[string]any {
	var ps []map[string]any

	for _, p := range api.Panels {
		// 初始化面板
		ql := &QlApi{
			URL:    p.URL,
			Token:  p.Token,
			Params: p.Params,
		}
		// 获取面板所有变量数据
		getEnvs, err := ql.GetEnvs()
		if err != nil {
			// 失效面板
			config.GinLOG.Error(err.Error())
			continue
		}

		// 判断面板存在变量
		if len(getEnvs.Data) <= 0 {
			ps = append(ps, map[string]any{
				"id":    p.ID,
				"count": api.Quantity,
			})
		} else {
			// 面板存在数据
			count := api.Quantity
			for _, x := range getEnvs.Data {
				if x.Name == api.Name {
					count--
				}
			}

			ps = append(ps, map[string]any{
				"id":    p.ID,
				"count": count,
			})
		}
	}

	// 判断是否有可用面板
	if len(ps) <= 0 {
		return nil
	}

	// 根据map中的count进行排序【降序】
	sort.Slice(ps, func(i, j int) bool {
		return ps[i]["count"].(int) > ps[j]["count"].(int)
	})
	return ps[0]
}

// GetPanelByEnvMode2 合并模式
func (api *QlApiFn) GetPanelByEnvMode2() map[string]any {
	var ps []map[string]any

	for _, p := range api.Panels {
		// 初始化面板
		ql := &QlApi{
			URL:    p.URL,
			Token:  p.Token,
			Params: p.Params,
		}
		// 获取面板所有变量数据
		getEnvs, err := ql.GetEnvs()
		if err != nil {
			// 失效面板
			config.GinLOG.Error(err.Error())
			continue
		}

		// 判断面板存在变量
		if len(getEnvs.Data) <= 0 {
			ps = append(ps, map[string]any{
				"id":    p.ID,
				"count": api.Quantity,
			})
		} else {
			// 面板存在数据
			count := api.Quantity
			for _, x := range getEnvs.Data {
				if x.Name == api.Name {
					// 根据合并分隔符分割变量值
					count -= len(strings.Split(x.Value, api.Division))
				}
			}

			ps = append(ps, map[string]any{
				"id":    p.ID,
				"count": count,
			})
		}
	}

	// 判断是否有可用面板
	if len(ps) <= 0 {
		return nil
	}

	// 根据map中的count进行排序【降序】
	sort.Slice(ps, func(i, j int) bool {
		return ps[i]["count"].(int) > ps[j]["count"].(int)
	})
	return ps[0]
}

// GetPanelByEnvMode3 更新模式
func (api *QlApiFn) GetPanelByEnvMode3() map[string]any {
	return nil
}
