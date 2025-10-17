package initializer

import (
	"io/fs"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nuanxinqing123/QLToolsV2/web"
)

// SetupWebFrontend 注册前端静态资源
func SetupWebFrontend(r *gin.Engine) {
	// 创建一个子文件系统，根目录为 dist/assets
	assetsFS, err := fs.Sub(web.DistFS, "dist/assets")
	if err != nil {
		panic(err)
	}
	// 提供静态资源
	r.StaticFS("/assets", http.FS(assetsFS))

	// 对于所有未匹配的路由，返回 index.html（支持 Vue Router history 模式）
	r.NoRoute(func(c *gin.Context) {
		// 获取文件内容
		data, err := web.DistFS.ReadFile("dist/index.html")
		if err != nil {
			c.String(http.StatusInternalServerError, "index.html not found")
			return
		}

		// 判断是否请求静态资源（比如 .js/.css），直接 404
		if strings.Contains(c.Request.RequestURI, ".") {
			c.String(http.StatusNotFound, "404 page not found")
			return
		}

		// 返回 HTML
		c.Data(http.StatusOK, "text/html; charset=utf-8", data)
	})
}
