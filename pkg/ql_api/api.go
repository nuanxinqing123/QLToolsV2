package ql_api

import (
	"errors"
	"fmt"
	"io"
	"net/url"
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
		return cfRes, errors.New(fmt.Sprintf("连接面板失败, 原因: %s", err))
	}

	if err = json.Unmarshal(bytes, &cfRes); err != nil {
		return cfRes, err
	}

	return cfRes, nil
}

type QlApi struct {
	URL    string // 连接地址
	Token  string // token
	Params int    // params
}

func InitPanel(url, token string, params int) *QlApi {
	return &QlApi{
		URL:    url,
		Token:  token,
		Params: params,
	}
}

// GetEnvs 获取环境变量列表
func (api *QlApi) GetEnvs() (EnvRes, error) {
	var res EnvRes

	// http://127.0.0.1:5700/api/envs?searchValue=&t=1713865007052
	ads := fmt.Sprintf("%s/open/envs", api.URL)

	params := url.Values{}
	params.Add("t", strconv.Itoa(api.Params))

	// 发送请求
	response, err := requests.Requests("GET", ads, params, "", api.Token)
	if err != nil {
		return res, err
	}

	bytes, err := io.ReadAll(response.Body)
	if err != nil {
		return res, errors.New(fmt.Sprintf("连接面板失败, 原因: %s", err))
	}

	if err = json.Unmarshal(bytes, &res); err != nil {
		return res, err
	}

	return res, nil
}

// PostEnvs 添加环境变量
func (api *QlApi) PostEnvs(env []PostEnv) (PostEnvRes, error) {
	var res PostEnvRes

	// http://127.0.0.1:5700/api/envs?t=1713865007052
	ads := fmt.Sprintf("%s/open/envs", api.URL)

	params := url.Values{}
	params.Add("t", strconv.Itoa(api.Params))

	// 转换为String
	bytes, err := json.Marshal(env)
	if err != nil {
		return res, err
	}

	// 发送请求
	response, err := requests.Requests("POST", ads, params, string(bytes), api.Token)
	if err != nil {
		return res, err
	}

	bytes, err = io.ReadAll(response.Body)
	if err != nil {
		return res, errors.New(fmt.Sprintf("连接面板失败, 原因: %s", err))
	}

	if err = json.Unmarshal(bytes, &res); err != nil {
		return res, err
	}

	return res, nil
}

// PutEnvs 更新环境变量
func (api *QlApi) PutEnvs(env PutEnv) (PutEnvRes, error) {
	var res PutEnvRes

	// http://127.0.0.1:5700/api/envs?t=1713865007052
	ads := fmt.Sprintf("%s/open/envs", api.URL)

	params := url.Values{}
	params.Add("t", strconv.Itoa(api.Params))

	// 转换为String
	bytes, err := json.Marshal(env)
	if err != nil {
		return res, err
	}

	// 发送请求
	response, err := requests.Requests("PUT", ads, params, string(bytes), api.Token)
	if err != nil {
		return res, err
	}

	bytes, err = io.ReadAll(response.Body)
	if err != nil {
		return res, errors.New(fmt.Sprintf("连接面板失败, 原因: %s", err))
	}

	if err = json.Unmarshal(bytes, &res); err != nil {
		return res, err
	}

	return res, nil
}

// PutDisableEnvs 禁用环境变量
func (api *QlApi) PutDisableEnvs(env PutDisableEnv) (PutDisableEnvRes, error) {
	var res PutDisableEnvRes

	// http://127.0.0.1:5700/api/envs?t=1713865007052
	ads := fmt.Sprintf("%s/open/envs/disable", api.URL)

	params := url.Values{}
	params.Add("t", strconv.Itoa(api.Params))

	// 转换为String
	bytes, err := json.Marshal(env)
	if err != nil {
		return res, err
	}

	// 发送请求
	response, err := requests.Requests("PUT", ads, params, string(bytes), api.Token)
	if err != nil {
		return res, err
	}

	bytes, err = io.ReadAll(response.Body)
	if err != nil {
		return res, errors.New(fmt.Sprintf("连接面板失败, 原因: %s", err))
	}

	if err = json.Unmarshal(bytes, &res); err != nil {
		return res, err
	}

	return res, nil
}

// PutEnableEnvs 启用环境变量
func (api *QlApi) PutEnableEnvs(env PutEnableEnv) (PutEnableEnvRes, error) {
	var res PutEnableEnvRes

	// http://127.0.0.1:5700/api/envs?t=1713865007052
	ads := fmt.Sprintf("%s/open/envs/enable", api.URL)

	params := url.Values{}
	params.Add("t", strconv.Itoa(api.Params))

	// 转换为String
	bytes, err := json.Marshal(env)
	if err != nil {
		return res, err
	}

	// 发送请求
	response, err := requests.Requests("PUT", ads, params, string(bytes), api.Token)
	if err != nil {
		return res, err
	}

	bytes, err = io.ReadAll(response.Body)
	if err != nil {
		return res, errors.New(fmt.Sprintf("连接面板失败, 原因: %s", err))
	}

	if err = json.Unmarshal(bytes, &res); err != nil {
		return res, err
	}

	return res, nil
}

// DeleteEnvs 删除环境变量
func (api *QlApi) DeleteEnvs(env DeleteEnv) (DeleteEnvRes, error) {
	var res DeleteEnvRes

	// http://127.0.0.1:5700/api/envs?t=1713865007052
	ads := fmt.Sprintf("%s/open/envs", api.URL)

	params := url.Values{}
	params.Add("t", strconv.Itoa(api.Params))

	// 转换为String
	bytes, err := json.Marshal(env)
	if err != nil {
		return res, err
	}

	// 发送请求
	response, err := requests.Requests("DELETE", ads, params, string(bytes), api.Token)
	if err != nil {
		return res, err
	}

	bytes, err = io.ReadAll(response.Body)
	if err != nil {
		return res, errors.New(fmt.Sprintf("连接面板失败, 原因: %s", err))
	}

	if err = json.Unmarshal(bytes, &res); err != nil {
		return res, err
	}

	return res, nil
}
