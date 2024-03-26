package utils

import (
	"strconv"
	"time"

	sf "github.com/bwmarrin/snowflake"
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
