package qinglong

import (
	"fmt"

	"github.com/nuanxinqing123/QLToolsV2/internal/app/config"
	"github.com/nuanxinqing123/QLToolsV2/internal/pkg/requests"
	"github.com/nuanxinqing123/QLToolsV2/internal/schema"
)

// QlConfig QL配置
type QlConfig struct {
	URL          string
	ClientID     string
	ClientSecret string

	client *requests.Request // client
}

// NewConfig 创建QL配置
func NewConfig(url, clientId, clientSecret string) *QlConfig {
	reqClient := requests.New()

	return &QlConfig{
		URL:          url,
		ClientID:     clientId,
		ClientSecret: clientSecret,

		client: reqClient,
	}
}

func (cfg *QlConfig) GetConfig() (schema.TokenResponse, error) {
	var res schema.TokenResponse

	ads := fmt.Sprintf("%s/open/auth/token", cfg.URL)

	params := map[string]string{
		"client_id":     cfg.ClientID,
		"client_secret": cfg.ClientSecret,
	}

	// 发送请求
	response, err := cfg.client.Get(ads, params)
	if err != nil {
		return res, err
	}

	if err = config.JSON.Unmarshal(response.Body(), &res); err != nil {
		return res, err
	}

	return res, nil
}
