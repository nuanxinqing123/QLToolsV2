package requests

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/nuanxinqing123/QLToolsV2/internal/app/config"
)

// Request HTTP请求客户端结构体
type Request struct {
	client *resty.Client // resty客户端实例
}

// New 创建一个新的HTTP请求客户端
func New() *Request {
	client := resty.New()

	// 设置DeBug模式
	// 当应用程序运行在debug模式时，开启HTTP请求的调试输出
	if config.Config.App.Mode == gin.DebugMode {
		client.SetDebug(true)
	}

	// 设置请求重试机制
	client.SetRetryCount(3)                  // 最大重试次数为3次
	client.SetRetryWaitTime(1 * time.Second) // 重试等待时间为1秒

	return &Request{
		client: client,
	}
}

// SetHeader 设置请求头
func (r *Request) SetHeader(header, value string) {
	r.client.SetHeader(header, value)
}

// Get 发送GET请求
// url: 请求地址
// params: URL查询参数
func (r *Request) Get(url string, params map[string]string) (*resty.Response, error) {
	// 记录请求日志
	config.Log.Debug("GET: " + url)
	config.Log.Debug(fmt.Sprintf("params: %v", params))

	// 发送GET请求
	return r.client.R().
		SetQueryParams(params).
		Get(url)
}

// Post 发送POST请求
// url: 请求地址
// body: 请求体数据
func (r *Request) Post(url string, body any) (*resty.Response, error) {
	// 记录请求日志
	config.Log.Debug("POST: " + url)
	config.Log.Debug(fmt.Sprintf("body: %v", body))

	// 发送POST请求
	return r.client.R().
		SetBody(body).
		Post(url)
}

// Put 发送PUT请求
// url: 请求地址
// body: 请求体数据
func (r *Request) Put(url string, body any) (*resty.Response, error) {
	// 记录请求日志
	config.Log.Debug("PUT: " + url)
	config.Log.Debug(fmt.Sprintf("body: %v", body))

	return r.client.R().
		SetBody(body).
		Put(url)
}

// Delete 发送DELETE请求
// url: 请求地址
// body: 请求体数据
func (r *Request) Delete(url string, body any) (*resty.Response, error) {
	// 记录请求日志
	config.Log.Debug("DELETE: " + url)
	config.Log.Debug(fmt.Sprintf("body: %v", body))

	// 发送DELETE请求
	return r.client.R().
		SetBody(body).
		Delete(url)
}
