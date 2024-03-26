package middleware

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"QLToolsV2/config"
	res "QLToolsV2/pkg/response"
)

// Recovery recover掉项目可能出现的panic，并使用zap记录相关日志
func Recovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		var errs validator.ValidationErrors
		if errors.As(recovered.(validator.ValidationErrors), &errs) {
			res.ResErrorWithMsg(c, res.CodeInvalidParam, "入参校验失败", errs)
			return
		}
		var err error
		if errors.As(recovered.(error), &err) {
			config.GinLOG.Error(err.Error())
			res.ResError(c, res.CodeServerBusy)
			return
		}
		c.AbortWithStatus(http.StatusInternalServerError)
	})
}
