package qinglong

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/go-resty/resty/v2"
	"github.com/nuanxinqing123/QLToolsV2/internal/app/config"
	"github.com/nuanxinqing123/QLToolsV2/internal/pkg/requests"
	"github.com/nuanxinqing123/QLToolsV2/internal/schema"
)

// TokenRefreshCallback token刷新回调函数类型
type TokenRefreshCallback func(panelID int64) (newToken string, err error)

// QlAPI QL API
type QlAPI struct {
	URL     string // 连接地址
	Token   string // token
	Params  int    // params
	PanelID int64  // 面板ID，用于token刷新

	client               *requests.Request    // client
	tokenRefreshCallback TokenRefreshCallback // token刷新回调函数
	isRefreshing         bool                 // 是否正在刷新token，防止并发刷新
	mutex                sync.Mutex           // 互斥锁，保护token刷新操作
}

// NewAPI 创建QL API
func NewAPI(url, token string, params int) *QlAPI {
	reqClient := requests.New()
	if token != "" {
		t := fmt.Sprintf("Bearer %s", token)
		reqClient.SetHeader("Authorization", t)
		reqClient.SetHeader("Token", t)
	}
	reqClient.SetHeader("User-Agent", "QLToolsV2")

	return &QlAPI{
		URL:     url,
		Token:   token,
		Params:  params,
		PanelID: 0, // 默认值，需要通过SetPanelID设置

		client:               reqClient,
		tokenRefreshCallback: nil, // 默认为空，需要通过SetTokenRefreshCallback设置
		isRefreshing:         false,
	}
}

// NewAPIWithPanel 创建带面板信息的QL API
func NewAPIWithPanel(url, token string, params int, panelID int64, callback TokenRefreshCallback) *QlAPI {
	api := NewAPI(url, token, params)
	api.PanelID = panelID
	api.tokenRefreshCallback = callback
	return api
}

// SetPanelID 设置面板ID
func (api *QlAPI) SetPanelID(panelID int64) {
	api.PanelID = panelID
}

// SetTokenRefreshCallback 设置token刷新回调函数
func (api *QlAPI) SetTokenRefreshCallback(callback TokenRefreshCallback) {
	api.tokenRefreshCallback = callback
}

// updateToken 更新API实例的token并重新设置请求头
func (api *QlAPI) updateToken(newToken string) {
	api.Token = newToken
	// 重新设置Authorization头
	authHeader := fmt.Sprintf("Bearer %s", newToken)
	api.client.SetHeader("Authorization", authHeader)
	api.client.SetHeader("Token", authHeader)
}

// refreshTokenIfNeeded 检查响应状态，如果是401则尝试刷新token
func (api *QlAPI) refreshTokenIfNeeded(response *resty.Response) error {
	// 检查是否为401未授权错误
	if response.StatusCode() != http.StatusUnauthorized {
		return nil
	}

	// 如果没有设置回调函数或面板ID，无法刷新token
	if api.tokenRefreshCallback == nil || api.PanelID == 0 {
		return fmt.Errorf("token已过期，但未设置刷新回调函数或面板ID")
	}

	// 使用互斥锁防止并发刷新
	api.mutex.Lock()
	defer api.mutex.Unlock()

	// 双重检查，防止在等待锁的过程中已经被其他goroutine刷新了
	if api.isRefreshing {
		return fmt.Errorf("token正在刷新中，请稍后重试")
	}

	api.isRefreshing = true
	defer func() {
		api.isRefreshing = false
	}()

	// 记录token刷新日志
	config.Log.Info(fmt.Sprintf("检测到401错误，开始刷新面板ID %d 的token", api.PanelID))

	// 调用回调函数刷新token
	newToken, err := api.tokenRefreshCallback(api.PanelID)
	if err != nil {
		config.Log.Error(fmt.Sprintf("刷新面板ID %d 的token失败: %v", api.PanelID, err))
		return fmt.Errorf("刷新token失败: %w", err)
	}

	// 更新token
	api.updateToken(newToken)
	config.Log.Info(fmt.Sprintf("面板ID %d 的token刷新成功", api.PanelID))

	return nil
}

// executeWithRetry 执行HTTP请求，如果遇到401错误则尝试刷新token后重试
func (api *QlAPI) executeWithRetry(requestFunc func() (*resty.Response, error)) (*resty.Response, error) {
	// 第一次尝试
	response, err := requestFunc()
	if err != nil {
		return response, err
	}

	// 检查是否需要刷新token
	if refreshErr := api.refreshTokenIfNeeded(response); refreshErr != nil {
		// 如果刷新失败，返回原始响应和刷新错误
		return response, refreshErr
	}

	// 如果刷新了token，重新执行请求
	if response.StatusCode() == http.StatusUnauthorized && api.tokenRefreshCallback != nil {
		config.Log.Debug("token刷新后重新执行请求")
		return requestFunc()
	}

	return response, nil
}

// GetEnvs 获取环境变量列表
func (api *QlAPI) GetEnvs() (schema.EnvResponse, error) {
	var res schema.EnvResponse

	// http://127.0.0.1:5700/api/envs?searchValue=&t=1713865007052
	ads := fmt.Sprintf("%s/open/envs", api.URL)

	// 使用executeWithRetry发送请求，自动处理401错误和token刷新
	response, err := api.executeWithRetry(func() (*resty.Response, error) {
		return api.client.Get(ads, nil)
	})
	if err != nil {
		return res, err
	}

	if err = config.JSON.Unmarshal(response.Body(), &res); err != nil {
		return res, err
	}

	return res, nil
}

// PostEnvs 添加环境变量
func (api *QlAPI) PostEnvs(env []schema.PostEnvRequest) (schema.PostEnvResponse, error) {
	var res schema.PostEnvResponse

	// http://127.0.0.1:5700/api/envs?t=1713865007052
	ads := fmt.Sprintf("%s/open/envs", api.URL)

	// 使用executeWithRetry发送请求，自动处理401错误和token刷新
	response, err := api.executeWithRetry(func() (*resty.Response, error) {
		return api.client.Post(ads, env)
	})
	if err != nil {
		return res, err
	}

	if err = config.JSON.Unmarshal(response.Body(), &res); err != nil {
		return res, err
	}

	return res, nil
}

// PutEnvs 更新环境变量
func (api *QlAPI) PutEnvs(env schema.PutEnvRequest) (schema.PutEnvResponse, error) {
	var res schema.PutEnvResponse

	// http://127.0.0.1:5700/api/envs?t=1713865007052
	ads := fmt.Sprintf("%s/open/envs", api.URL)

	// 使用executeWithRetry发送请求，自动处理401错误和token刷新
	response, err := api.executeWithRetry(func() (*resty.Response, error) {
		return api.client.Put(ads, env)
	})
	if err != nil {
		return res, err
	}

	if err = config.JSON.Unmarshal(response.Body(), &res); err != nil {
		return res, err
	}

	return res, nil
}

// PutDisableEnvs 禁用环境变量
func (api *QlAPI) PutDisableEnvs(env schema.PutDisableEnvRequest) (schema.PutDisableEnvResponse, error) {
	var res schema.PutDisableEnvResponse

	// http://127.0.0.1:5700/api/envs?t=1713865007052
	ads := fmt.Sprintf("%s/open/envs/disable", api.URL)

	// 使用executeWithRetry发送请求，自动处理401错误和token刷新
	response, err := api.executeWithRetry(func() (*resty.Response, error) {
		return api.client.Put(ads, env)
	})
	if err != nil {
		return res, err
	}

	if err = config.JSON.Unmarshal(response.Body(), &res); err != nil {
		return res, err
	}

	return res, nil
}

// PutEnableEnvs 启用环境变量
func (api *QlAPI) PutEnableEnvs(env schema.PutEnableEnvRequest) (schema.PutEnableEnvResponse, error) {
	var res schema.PutEnableEnvResponse

	// http://127.0.0.1:5700/api/envs?t=1713865007052
	ads := fmt.Sprintf("%s/open/envs/enable", api.URL)

	// 使用executeWithRetry发送请求，自动处理401错误和token刷新
	response, err := api.executeWithRetry(func() (*resty.Response, error) {
		return api.client.Put(ads, env)
	})
	if err != nil {
		return res, err
	}

	if err = config.JSON.Unmarshal(response.Body(), &res); err != nil {
		return res, err
	}

	return res, nil
}

// DeleteEnvs 删除环境变量
func (api *QlAPI) DeleteEnvs(env schema.DeleteEnvRequest) (schema.DeleteEnvResponse, error) {
	var res schema.DeleteEnvResponse

	// http://127.0.0.1:5700/api/envs?t=1713865007052
	ads := fmt.Sprintf("%s/open/envs", api.URL)

	// 使用executeWithRetry发送请求，自动处理401错误和token刷新
	response, err := api.executeWithRetry(func() (*resty.Response, error) {
		return api.client.Delete(ads, env)
	})
	if err != nil {
		return res, err
	}

	if err = config.JSON.Unmarshal(response.Body(), &res); err != nil {
		return res, err
	}

	return res, nil
}
