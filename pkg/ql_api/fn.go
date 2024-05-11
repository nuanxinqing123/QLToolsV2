package ql_api

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"QLToolsV2/config"
	"QLToolsV2/internal/db"
)

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

			// 变量总数
			envTotal := xx.Quantity * len(xx.Panels)
			config.GinLOG.Debug(fmt.Sprintf("变量总数: %d", envTotal))
			config.GinLOG.Debug(fmt.Sprintf("变量名称: %s, 变量模式: %d", xx.Name, xx.Mode))
			for _, y := range xx.Panels {
				config.GinLOG.Debug(fmt.Sprintf("面板名称: %s", y.Name))
				// 初始化面板
				ql := &QlApi{
					URL:    y.URL,
					Token:  y.Token,
					Params: y.Params,
				}
				// 获取面板所有变量数据
				getEnvs, err := ql.GetEnvs()
				if err != nil {
					// 减去失效的面板配额
					envTotal -= xx.Quantity
					config.GinLOG.Error(err.Error())
					continue
				}

				// 面板存在变量
				config.GinLOG.Debug(fmt.Sprintf("面板名称: %s, 面板存在变量: %d个", y.Name, len(getEnvs.Data)))
				if len(getEnvs.Data) > 0 {
					for _, z := range getEnvs.Data {
						config.GinLOG.Debug(fmt.Sprintf("变量名称: %s, 变量值: %s, 变量备注: %s", z.Name, z.Value, z.Remarks))
						if xx.Mode == 1 || xx.Mode == 3 {
							// 新建模式 || 更新模式
							if z.Name == xx.Name {
								envTotal--
							}
						} else {
							// 合并模式
							if z.Name == xx.Name {
								zL := strings.Split(z.Value, xx.Division)
								envTotal -= len(zL)
							}
						}
					}
				}
			}
			config.GinLOG.Debug(fmt.Sprintf("可使用变量总数: %d", envTotal))

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
