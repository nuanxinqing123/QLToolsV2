package middleware

import (
	"bytes"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"QLToolsV2/config"
	_const "QLToolsV2/const"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}
func (w bodyLogWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		bodyLogWriter := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = bodyLogWriter

		// 开始时间
		startTime := time.Now()

		b, _ := c.Copy().GetRawData()

		c.Request.Body = io.NopCloser(bytes.NewReader(b))

		// 处理请求
		c.Next()

		// 结束时间
		endTime := time.Now()

		config.GinLOG.Info("请求响应",
			zap.String("client_ip", c.ClientIP()),
			zap.String("request_time", TimeFormat(startTime)),
			zap.String("response_time", TimeFormat(endTime)),
			zap.String("cost_time", endTime.Sub(startTime).String()),
		)
	}
}

func TimeFormat(t time.Time) string {
	return t.Format(_const.TimeFormatAll)
}
