package ql_api

import (
	"time"
)

// TokenRes 获取 Token【返回】
type TokenRes struct {
	Code int `json:"code"`
	Data struct {
		Token      string `json:"token"`
		TokenType  string `json:"token_type"`
		Expiration int    `json:"expiration"`
	} `json:"data"`
	Message string
}

// EnvRes 获取环境变量【返回】
type EnvRes struct {
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

// PostEnv 创建环境变量
type PostEnv struct {
	Name    string `json:"name"`
	Value   string `json:"value"`
	Remarks string `json:"remarks,omitempty"`
}

// PostEnvRes 创建环境变量【返回】
type PostEnvRes struct {
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

// PutEnv 更新环境变量
type PutEnv struct {
	Name    string `json:"name"`
	Value   string `json:"value"`
	Remarks string `json:"remarks"`
	Id      int    `json:"id"`
}

// PutEnvRes 更新环境变量【返回】
type PutEnvRes struct {
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

// PutDisableEnv 禁用环境变量
type PutDisableEnv []int

// PutDisableEnvRes 禁用环境变量【返回】
type PutDisableEnvRes struct {
	Code int `json:"code"`
}

// PutEnableEnv 启用环境变量
type PutEnableEnv []int

// PutEnableEnvRes 启用环境变量【返回】
type PutEnableEnvRes struct {
	Code int `json:"code"`
}

// DeleteEnv 删除环境变量
type DeleteEnv []int

// DeleteEnvRes 删除环境变量【返回】
type DeleteEnvRes struct {
	Code int `json:"code"`
}
