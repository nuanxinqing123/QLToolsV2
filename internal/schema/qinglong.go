package schema

import "time"

// TokenResponse 获取 Token【返回】
type TokenResponse struct {
	Code int `json:"code"`
	Data struct {
		Token      string `json:"token"`
		TokenType  string `json:"token_type"`
		Expiration int    `json:"expiration"`
	} `json:"data"`
	Message string
}

// EnvResponse 获取环境变量【返回】
type EnvResponse struct {
	Code int `json:"code"`
	Data []struct {
		Id        int       `json:"id"`
		Value     string    `json:"value"`
		Timestamp string    `json:"timestamp"`
		Status    int       `json:"status"`
		Position  int64     `json:"position"`
		Name      string    `json:"name"`
		Remarks   string    `json:"remarks"`
		CreatedAt time.Time `json:"createdAt"`
		UpdatedAt time.Time `json:"updatedAt"`
	} `json:"data"`
}

// PostEnvRequest 创建环境变量
type PostEnvRequest struct {
	Name    string `json:"name"`
	Value   string `json:"value"`
	Remarks string `json:"remarks,omitempty"`
}

// PostEnvResponse 创建环境变量【返回】
type PostEnvResponse struct {
	Code int `json:"code"`
	Data []struct {
		Id        int       `json:"id"`
		Value     string    `json:"value"`
		Status    int       `json:"status"`
		Timestamp string    `json:"timestamp"`
		Position  int64     `json:"position"`
		Name      string    `json:"name"`
		Remarks   string    `json:"remarks"`
		UpdatedAt time.Time `json:"updatedAt"`
		CreatedAt time.Time `json:"createdAt"`
	} `json:"data"`
}

// PutEnvRequest 更新环境变量
type PutEnvRequest struct {
	Name    string `json:"name"`
	Value   string `json:"value"`
	Remarks string `json:"remarks"`
	Id      int    `json:"id"`
}

// PutEnvResponse 更新环境变量【返回】
type PutEnvResponse struct {
	Code int `json:"code"`
	Data struct {
		Id        int       `json:"id"`
		Value     string    `json:"value"`
		Timestamp string    `json:"timestamp"`
		Status    int       `json:"status"`
		Position  int64     `json:"position"`
		Name      string    `json:"name"`
		Remarks   string    `json:"remarks"`
		CreatedAt time.Time `json:"createdAt"`
		UpdatedAt time.Time `json:"updatedAt"`
	} `json:"data"`
}

// PutDisableEnvRequest 禁用环境变量
type PutDisableEnvRequest []int

// PutDisableEnvResponse 禁用环境变量【返回】
type PutDisableEnvResponse struct {
	Code int `json:"code"`
}

// PutEnableEnvRequest 启用环境变量
type PutEnableEnvRequest []int

// PutEnableEnvResponse 启用环境变量【返回】
type PutEnableEnvResponse struct {
	Code int `json:"code"`
}

// DeleteEnvRequest 删除环境变量
type DeleteEnvRequest []int

// DeleteEnvResponse 删除环境变量【返回】
type DeleteEnvResponse struct {
	Code int `json:"code"`
}
