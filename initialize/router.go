package initialize

import (
	"html/template"
	"strings"

	"github.com/gin-gonic/gin"

	assetfs "github.com/elazarl/go-bindata-assetfs"

	"QLToolsV2/internal/middleware"
	"QLToolsV2/internal/router"
	res "QLToolsV2/pkg/response"
	"QLToolsV2/web/bindata"
)

func Routers() *gin.Engine {
	Router := gin.New()
	Router.Use(middleware.Logger())
	Router.Use(middleware.Recovery())

	// 允许跨域
	Router.Use(middleware.Cors())

	// (可选项)
	// PID 限流基于实例的 CPU 使用率，通过拒绝一定比例的流量, 将实例的 CPU 使用率稳定在设定的阈值上。
	// 地址: https://github.com/bytedance/pid_limits
	// Router.Use(adaptive.PlatoMiddlewareGinDefault(0.8))

	// 前端静态文件
	{
		// 加载模板文件
		t, err := loadTemplate()
		if err != nil {
			panic(err)
		}
		Router.SetHTMLTemplate(t)

		// 加载静态文件
		fs := &assetfs.AssetFS{
			Asset:     bindata.Asset,
			AssetDir:  bindata.AssetDir,
			AssetInfo: nil,
			Prefix:    "assets",
		}
		Router.StaticFS("/assets", fs)

		Router.GET("/", func(c *gin.Context) {
			c.HTML(200, "index.html", nil)
		})
	}

	// 存活检测
	Router.GET("/ping", func(c *gin.Context) {
		res.ResSuccess(c, "pong")
	})

	ApiGroupOpen := Router.Group("/api")
	router.InitRouterOpen(ApiGroupOpen)

	ApiGroupAdmin := Router.Group("/api/admin")
	router.InitRouterAdmin(ApiGroupAdmin)

	return Router
}

// loadTemplate 加载模板文件
func loadTemplate() (*template.Template, error) {
	t := template.New("")
	for _, name := range bindata.AssetNames() {
		if !strings.HasSuffix(name, ".html") {
			continue
		}
		asset, err := bindata.Asset(name)
		if err != nil {
			continue
		}
		name = strings.Replace(name, "assets/", "", 1)
		t, err = t.New(name).Parse(string(asset))
		if err != nil {
			return nil, err
		}
	}
	return t, nil
}
