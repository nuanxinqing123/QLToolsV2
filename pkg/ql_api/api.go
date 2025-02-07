package ql_api

import (
	"fmt"
	"strconv"

	jsoniter "github.com/json-iterator/go"

	"QLToolsV2/pkg/requests"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type QlConfig struct {
	URL          string
	ClientID     string
	ClientSecret string
}

func InitConfig(url, clientId, clientSecret string) *QlConfig {
	return &QlConfig{
		URL:          url,
		ClientID:     clientId,
		ClientSecret: clientSecret,
	}
}

func (api *QlConfig) GetConfig() (TokenRes, error) {
	var res TokenRes

	ads := fmt.Sprintf("%s/open/auth/token", api.URL)

	params := map[string]string{
		"client_id":     api.ClientID,
		"client_secret": api.ClientSecret,
	}

	// 发送请求
	client := requests.New("", params)
	response, err := client.Get(ads)
	if err != nil {
		return res, err
	}

	if err = json.Unmarshal(response.Body(), &res); err != nil {
		return res, err
	}

	return res, nil
}

type QlApi struct {
	URL    string // 连接地址
	Token  string // token
	Params int    // params

	client *requests.Request // client
}

func InitPanel(url, token string, params int) *QlApi {
	return &QlApi{
		URL:    url,
		Token:  token,
		Params: params,

		client: requests.New(token, map[string]string{"t": strconv.Itoa(params)}),
	}
}

// GetEnvs 获取环境变量列表
func (api *QlApi) GetEnvs() (EnvRes, error) {
	var res EnvRes

	// http://127.0.0.1:5700/api/envs?searchValue=&t=1713865007052
	ads := fmt.Sprintf("%s/open/envs", api.URL)

	// 发送请求
	response, err := api.client.Get(ads)
	if err != nil {
		return res, err
	}

	if err = json.Unmarshal(response.Body(), &res); err != nil {
		return res, err
	}

	return res, nil
}

// PostEnvs 添加环境变量
func (api *QlApi) PostEnvs(env []PostEnv) (PostEnvRes, error) {
	var res PostEnvRes

	// http://127.0.0.1:5700/api/envs?t=1713865007052
	ads := fmt.Sprintf("%s/open/envs", api.URL)

	// 转换为String
	bytes, err := json.Marshal(env)
	if err != nil {
		return res, err
	}

	// 发送请求
	response, err := api.client.Post(ads, string(bytes))
	if err != nil {
		return res, err
	}

	if err = json.Unmarshal(response.Body(), &res); err != nil {
		return res, err
	}

	return res, nil
}

// PutEnvs 更新环境变量
func (api *QlApi) PutEnvs(env PutEnv) (PutEnvRes, error) {
	var res PutEnvRes

	// http://127.0.0.1:5700/api/envs?t=1713865007052
	ads := fmt.Sprintf("%s/open/envs", api.URL)

	// 转换为String
	bytes, err := json.Marshal(env)
	if err != nil {
		return res, err
	}

	// 发送请求
	response, err := api.client.Put(ads, string(bytes))
	if err != nil {
		return res, err
	}

	if err = json.Unmarshal(response.Body(), &res); err != nil {
		return res, err
	}

	return res, nil
}

// PutDisableEnvs 禁用环境变量
func (api *QlApi) PutDisableEnvs(env PutDisableEnv) (PutDisableEnvRes, error) {
	var res PutDisableEnvRes

	// http://127.0.0.1:5700/api/envs?t=1713865007052
	ads := fmt.Sprintf("%s/open/envs/disable", api.URL)

	// 转换为String
	bytes, err := json.Marshal(env)
	if err != nil {
		return res, err
	}

	// 发送请求
	response, err := api.client.Put(ads, string(bytes))
	if err != nil {
		return res, err
	}

	if err = json.Unmarshal(response.Body(), &res); err != nil {
		return res, err
	}

	return res, nil
}

// PutEnableEnvs 启用环境变量
func (api *QlApi) PutEnableEnvs(env PutEnableEnv) (PutEnableEnvRes, error) {
	var res PutEnableEnvRes

	// http://127.0.0.1:5700/api/envs?t=1713865007052
	ads := fmt.Sprintf("%s/open/envs/enable", api.URL)

	// 转换为String
	bytes, err := json.Marshal(env)
	if err != nil {
		return res, err
	}

	// 发送请求
	response, err := api.client.Put(ads, string(bytes))
	if err != nil {
		return res, err
	}

	if err = json.Unmarshal(response.Body(), &res); err != nil {
		return res, err
	}

	return res, nil
}

// DeleteEnvs 删除环境变量
func (api *QlApi) DeleteEnvs(env DeleteEnv) (DeleteEnvRes, error) {
	var res DeleteEnvRes

	// http://127.0.0.1:5700/api/envs?t=1713865007052
	ads := fmt.Sprintf("%s/open/envs", api.URL)

	// 转换为String
	bytes, err := json.Marshal(env)
	if err != nil {
		return res, err
	}

	// 发送请求
	response, err := api.client.Delete(ads, string(bytes))
	if err != nil {
		return res, err
	}

	if err = json.Unmarshal(response.Body(), &res); err != nil {
		return res, err
	}

	return res, nil
}
