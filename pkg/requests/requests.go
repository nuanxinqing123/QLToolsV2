package requests

import (
	"QLToolsV2/config"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
)

// Requests 封装HTTP请求
//func Requests(method, url string, params url.Values, data, token string) (*http.Response, error) {
//	// 创建HTTP实例
//	client := &http.Client{}
//
//	// 添加请求数据
//	var ReqData = strings.NewReader(data)
//	req, err := http.NewRequest(method, url, ReqData)
//
//	// 添加Params数据
//	req.URL.RawQuery = params.Encode()
//
//	// 添加请求Token
//	if token != "" {
//		Token := fmt.Sprintf("Bearer %s", token)
//		req.Header.Set("Authorization", Token)
//		req.Header.Set("Token", token)
//	}
//	req.Header.Set("Content-Type", "application/json")
//	req.Header.Set("User-Agent", "Easy-Gin")
//	// 发送请求
//	resp, err := client.Do(req)
//	if err != nil {
//		return resp, err
//	}
//	return resp, nil
//}

// Request HTTP请求客户端结构体
type Request struct {
	client *resty.Client // resty客户端实例
}

// New 创建一个新的HTTP请求客户端
func New(token string, params map[string]string) *Request {
	client := resty.New()

	// 设置DeBug模式
	// 当应用程序运行在debug模式时，开启HTTP请求的调试输出
	if config.GinConfig.App.Mode == "debug" {
		client.SetDebug(true)
	}

	// 添加Headers
	if token != "" {
		client.SetAuthToken(token)
	}
	client.SetHeader("User-Agent", "QLToolsV2")

	// 设置请求重试机制
	client.SetRetryCount(3)                  // 最大重试次数为3次
	client.SetRetryWaitTime(1 * time.Second) // 重试等待时间为1秒

	// 设置URL查询参数
	if params != nil {
		config.GinLOG.Debug(fmt.Sprintf("params: %v", params))
		client.SetQueryParams(params)
	}

	return &Request{
		client: client,
	}
}

// Get 发送GET请求
// url: 请求地址
// params: URL查询参数
func (r *Request) Get(url string) (*resty.Response, error) {
	// 记录请求日志
	config.GinLOG.Debug("GET: " + url)

	// 发送GET请求
	return r.client.R().
		Get(url)
}

// Post 发送POST请求
// url: 请求地址
// body: 请求体数据
func (r *Request) Post(url string, body any) (*resty.Response, error) {
	// 记录请求日志
	config.GinLOG.Debug("POST: " + url)
	config.GinLOG.Debug(fmt.Sprintf("body: %v", body))

	// 发送POST请求
	return r.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Post(url)
}

// Put 发送PUT请求
// url: 请求地址
// body: 请求体数据
func (r *Request) Put(url string, body any) (*resty.Response, error) {
	// 记录请求日志
	config.GinLOG.Debug("PUT: " + url)
	config.GinLOG.Debug(fmt.Sprintf("body: %v", body))

	// 发送PUT请求
	return r.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Put(url)
}

// Delete 发送DELETE请求
// url: 请求地址
// body: 请求体数据
func (r *Request) Delete(url string, body any) (*resty.Response, error) {
	// 记录请求日志
	config.GinLOG.Debug("DELETE: " + url)
	config.GinLOG.Debug(fmt.Sprintf("body: %v", body))

	// 发送DELETE请求
	return r.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Delete(url)
}
