package ql_api

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	jsoniter "github.com/json-iterator/go"

	"QLToolsV2/config"
	"QLToolsV2/pkg/requests"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type QlConfig struct {
	URL          string
	ClientID     string
	ClientSecret string
}

func (api QlConfig) InitConfig(url, clientId, clientSecret string) *QlConfig {
	api.URL = url
	api.ClientID = clientId
	api.ClientSecret = clientSecret
	return &api
}

func (api QlConfig) GetConfig() (TokenRes, error) {
	var cfRes TokenRes

	ads := fmt.Sprintf("%s/open/auth/token", api.URL)

	params := url.Values{}
	params.Add("client_id", api.ClientID)
	params.Add("client_secret", api.ClientSecret)

	// 发送请求
	response, err := requests.Requests("GET", ads, params, "", "")
	if err != nil {
		return cfRes, err
	}

	bytes, err := io.ReadAll(response.Body)
	if err != nil {
		config.GinLOG.Error(err.Error())
		return cfRes, errors.New(fmt.Sprintf("连接面板失败, 原因: %s", err))
	}

	if err = json.Unmarshal(bytes, &cfRes); err != nil {
		config.GinLOG.Error(err.Error())
		return cfRes, err
	}

	return cfRes, nil
}

type QlApi struct {
	URL    string // 连接地址
	Token  string // token
	Params string // params
}

func (api QlApi) InitPanel(url, token, params string) *QlApi {
	api.URL = url
	api.Token = token
	api.Params = params
	return &api
}

// GetEnvs 获取环境变量列表
func (api QlApi) GetEnvs() (*http.Response, error) {
	// http://127.0.0.1:5700/api/envs?searchValue=&t=1713865007052
	ads := fmt.Sprintf("%s/open/envs", api.URL)

	params := url.Values{}
	params.Add("t", api.Params)

	// 发送请求
	response, err := requests.Requests("GET", ads, params, "", api.Token)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// PostEnvs 添加环境变量
func (api QlApi) PostEnvs(env PostEnv) (*http.Response, error) {
	// http://127.0.0.1:5700/api/envs?t=1713865007052
	ads := fmt.Sprintf("%s/open/envs", api.URL)

	params := url.Values{}
	params.Add("t", api.Params)

	// 转换为String
	bytes, err := json.Marshal(env)
	if err != nil {
		return nil, err
	}

	// 发送请求
	response, err := requests.Requests("POST", ads, params, string(bytes), api.Token)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// PutEnvs 更新环境变量
func (api QlApi) PutEnvs(env PutEnv) (*http.Response, error) {
	// http://127.0.0.1:5700/api/envs?t=1713865007052
	ads := fmt.Sprintf("%s/open/envs", api.URL)

	params := url.Values{}
	params.Add("t", api.Params)

	// 转换为String
	bytes, err := json.Marshal(env)
	if err != nil {
		return nil, err
	}

	// 发送请求
	response, err := requests.Requests("PUT", ads, params, string(bytes), api.Token)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// PutDisableEnvs 禁用环境变量
func (api QlApi) PutDisableEnvs(env PutDisableEnv) (*http.Response, error) {
	// http://127.0.0.1:5700/api/envs?t=1713865007052
	ads := fmt.Sprintf("%s/open/envs/disable", api.URL)

	params := url.Values{}
	params.Add("t", api.Params)

	// 转换为String
	bytes, err := json.Marshal(env)
	if err != nil {
		return nil, err
	}

	// 发送请求
	response, err := requests.Requests("PUT", ads, params, string(bytes), api.Token)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// PutEnableEnvs 启用环境变量
func (api QlApi) PutEnableEnvs(env PutEnableEnv) (*http.Response, error) {
	// http://127.0.0.1:5700/api/envs?t=1713865007052
	ads := fmt.Sprintf("%s/open/envs/enable", api.URL)

	params := url.Values{}
	params.Add("t", api.Params)

	// 转换为String
	bytes, err := json.Marshal(env)
	if err != nil {
		return nil, err
	}

	// 发送请求
	response, err := requests.Requests("PUT", ads, params, string(bytes), api.Token)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// DeleteEnvs 删除环境变量
func (api QlApi) DeleteEnvs(env DeleteEnv) (*http.Response, error) {
	// http://127.0.0.1:5700/api/envs?t=1713865007052
	ads := fmt.Sprintf("%s/open/envs", api.URL)

	params := url.Values{}
	params.Add("t", api.Params)

	// 转换为String
	bytes, err := json.Marshal(env)
	if err != nil {
		return nil, err
	}

	// 发送请求
	response, err := requests.Requests("DELETE", ads, params, string(bytes), api.Token)
	if err != nil {
		return nil, err
	}

	return response, nil
}
