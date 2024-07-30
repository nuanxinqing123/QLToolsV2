package ip

import (
	"errors"
	"fmt"
	"io"
	"net/url"

	jsoniter "github.com/json-iterator/go"

	"QLToolsV2/pkg/requests"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type IP struct{}

type ResIP struct {
	/*
		{
		    "country": "中国",
		    "short_name": "CN",
		    "province": "广东省",
		    "city": "广州市",
		    "area": "海珠区",
		    "isp": "电信",
		    "net": "",
		    "ip": "113.68.137.1",
		    "code": 200,
		    "desc": "success"
		}
	*/
	Country   string `json:"country"`
	ShortName string `json:"short_name"`
	Province  string `json:"province"`
	City      string `json:"city"`
	Area      string `json:"area"`
	Isp       string `json:"isp"`
	Net       string `json:"net"`
	Ip        string `json:"ip"`
	Code      int    `json:"code"`
	Desc      string `json:"desc"`
}

func InitIP() *IP {
	return &IP{}
}

func (i *IP) GetIP(ip string) (ResIP, error) {
	var res ResIP

	// 查询IP地址
	ads := fmt.Sprintf("https://ip.useragentinfIo.com/sp/TZb2y?ip=%s", ip)
	params := url.Values{}

	response, err := requests.Requests("POST", ads, params, "", "")
	if err != nil {
		return res, err
	}

	bytes, err := io.ReadAll(response.Body)
	if err != nil {
		return res, errors.New(fmt.Sprintf("查询数据, 原因: %s", err))
	}

	if err = json.Unmarshal(bytes, &res); err != nil {
		return res, err
	}

	return res, nil
}
