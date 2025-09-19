package middleware

import (
	"bytes"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nuanxinqing123/QLToolsV2/internal/app/config"
	_const "github.com/nuanxinqing123/QLToolsV2/internal/const"
	"go.uber.org/zap"
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
	return func(ctx *gin.Context) {
		bodyLogWriter := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: ctx.Writer}
		ctx.Writer = bodyLogWriter

		// 开始时间
		startTime := time.Now()

		b, _ := ctx.Copy().GetRawData()

		ctx.Request.Body = io.NopCloser(bytes.NewReader(b))

		// 处理请求
		ctx.Next()

		// 结束时间
		endTime := time.Now()

		config.Log.Info("请求响应",
			zap.Int("status", ctx.Writer.Status()),
			zap.String("method", ctx.Request.Method),
			zap.String("url", ctx.Request.URL.String()),
			zap.String("client_ip", ctx.ClientIP()),
			zap.String("request_time", TimeFormat(startTime)),
			zap.String("response_time", TimeFormat(endTime)),
			zap.String("cost_time", endTime.Sub(startTime).String()),
		)
	}
}

func TimeFormat(t time.Time) string {
	return t.Format(_const.TimeFormatAll)
}
