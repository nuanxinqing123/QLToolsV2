package requests

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// Requests 封装HTTP请求
func Requests(method, url string, params url.Values, data, token string) (*http.Response, error) {
	// 创建HTTP实例
	client := &http.Client{}

	// 添加请求数据
	var ReqData = strings.NewReader(data)
	req, err := http.NewRequest(method, url, ReqData)

	// 添加Params数据
	req.URL.RawQuery = params.Encode()

	// 添加请求Token
	if token != "" {
		Token := fmt.Sprintf("Bearer %s", token)
		req.Header.Set("Authorization", Token)
		req.Header.Set("Token", token)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Easy-Gin")
	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return resp, err
	}
	return resp, nil
}
