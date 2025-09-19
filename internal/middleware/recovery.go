package middleware

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/nuanxinqing123/QLToolsV2/internal/app/config"
)

// Recovery recover掉项目可能出现的panic，并使用zap记录相关日志
func Recovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered any) {
		var errs validator.ValidationErrors
		if errors.As(recovered.(validator.ValidationErrors), &errs) {
			c.String(http.StatusBadRequest, "参数校验失败")
			return
		}
		var err error
		if errors.As(recovered.(error), &err) {
			config.Log.Error(err.Error())
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		c.AbortWithStatus(http.StatusInternalServerError)
	})
}
