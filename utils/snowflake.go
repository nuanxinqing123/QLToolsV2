package utils

import (
	"strconv"
	"time"

	sf "github.com/bwmarrin/snowflake"
	"golang.org/x/exp/rand"
)

var node *sf.Node

func InitSnowflake() (err error) {
	var st time.Time
	st, err = time.Parse("2006-01-02", "2023-03-15")
	if err != nil {
		return
	}

	sf.Epoch = st.UnixNano() / 1000000
	node, err = sf.NewNode(1)
	return
}

// GenID 生成唯一ID
func GenID() string {
	return strconv.FormatInt(node.Generate().Int64(), 10)
}

// GenRandomString 生成指定长度的随机字符串
func GenRandomString(length int) string {
	rand.Seed(uint64(time.Now().Unix()))
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	result := make([]rune, length)
	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}
	return string(result)
}
