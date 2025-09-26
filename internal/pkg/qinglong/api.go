package qinglong

import (
	"fmt"

	"github.com/nuanxinqing123/QLToolsV2/internal/app/config"
	"github.com/nuanxinqing123/QLToolsV2/internal/pkg/requests"
	"github.com/nuanxinqing123/QLToolsV2/internal/schema"
)

// QlAPI QL API
type QlAPI struct {
	URL    string // 连接地址
	Token  string // token
	Params int    // params

	client *requests.Request // client
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
		URL:    url,
		Token:  token,
		Params: params,

		client: reqClient,
	}
}

// GetEnvs 获取环境变量列表
func (api *QlAPI) GetEnvs() (schema.EnvResponse, error) {
	var res schema.EnvResponse

	// http://127.0.0.1:5700/api/envs?searchValue=&t=1713865007052
	ads := fmt.Sprintf("%s/open/envs", api.URL)

	// 发送请求
	response, err := api.client.Get(ads, nil)
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

	// 转换为String
	bytes, err := config.JSON.Marshal(env)
	if err != nil {
		return res, err
	}

	// 发送请求
	response, err := api.client.Post(ads, string(bytes))
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

	// 转换为String
	bytes, err := config.JSON.Marshal(env)
	if err != nil {
		return res, err
	}

	// 发送请求
	response, err := api.client.Put(ads, string(bytes))
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

	// 转换为String
	bytes, err := config.JSON.Marshal(env)
	if err != nil {
		return res, err
	}

	// 发送请求
	response, err := api.client.Put(ads, string(bytes))
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

	// 转换为String
	bytes, err := config.JSON.Marshal(env)
	if err != nil {
		return res, err
	}

	// 发送请求
	response, err := api.client.Put(ads, string(bytes))
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

	// 转换为String
	bytes, err := config.JSON.Marshal(env)
	if err != nil {
		return res, err
	}

	// 发送请求
	response, err := api.client.Delete(ads, string(bytes))
	if err != nil {
		return res, err
	}

	if err = config.JSON.Unmarshal(response.Body(), &res); err != nil {
		return res, err
	}

	return res, nil
}
