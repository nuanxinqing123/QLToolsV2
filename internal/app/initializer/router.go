package initializer

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/nuanxinqing123/QLToolsV2/docs"
	"github.com/nuanxinqing123/QLToolsV2/internal/app/config"
	"github.com/nuanxinqing123/QLToolsV2/internal/controller"
	"github.com/nuanxinqing123/QLToolsV2/internal/middleware"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	ginprometheus "github.com/zsais/go-gin-prometheus"
)

func Routers() *gin.Engine {
	Router := gin.New()
	Router.Use(middleware.Logger())
	Router.Use(middleware.Recovery())

	// 跨域配置
	Router.Use(cors.New(middleware.CorsConfig))

	// (可选项)
	// PID 限流基于实例的 CPU 使用率，通过拒绝一定比例的流量, 将实例的 CPU 使用率稳定在设定的阈值上。
	// 地址: https://github.com/bytedance/pid_limits
	// Router.Use(adaptive.PlatoMiddlewareGinDefault(0.8))

	// 初始化 Prometheus 中间件
	p := ginprometheus.NewPrometheus("gin")
	p.Use(Router)

	// 存活检测
	Router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	// 注册Swagger
	if config.Config.App.Mode == gin.DebugMode {
		Router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	api := Router.Group("/api")

	// 健康检查
	HealthyGroup := api.Group("/")
	HealthyCon := controller.NewHealthyController()
	HealthyCon.HealthyRouter(HealthyGroup)

	// 公开服务（应用令牌桶限速）
	OpenGroup := api.Group("/open")
	OpenGroup.Use(middleware.OpenAPIRateLimit()) // 应用令牌桶限速中间件
	OpenCon := controller.NewOpenController()
	OpenCon.OpenRouter(OpenGroup)

	// 认证校验
	AuthGroup := api.Group("/auth")
	AuthCon := controller.NewAuthController()
	AuthCon.AuthRouter(AuthGroup)

	// 认证接口
	authAPI := api.Group("")
	authAPI.Use(middleware.JWTAuth()) // 校验请求认证
	{
		// 认证
		AuthRequiredGroup := authAPI.Group("/auth")
		AuthRequiredCon := controller.NewAuthRequiredController()
		AuthRequiredCon.AuthRequiredRouter(AuthRequiredGroup)

		// 仪表盘
		DashboardGroup := authAPI.Group("/dashboard")
		DashboardCon := controller.NewDashboardController()
		DashboardCon.DashboardRouter(DashboardGroup)

		// 面板管理
		PanelGroup := authAPI.Group("/panel")
		PanelCon := controller.NewPanelController()
		PanelCon.PanelRouter(PanelGroup)

		// 变量管理
		EnvGroup := authAPI.Group("/env")
		EnvCon := controller.NewEnvController()
		EnvCon.EnvRouter(EnvGroup)

		// CDK管理
		CDKGroup := authAPI.Group("/cdk")
		CDKCon := controller.NewCDKController()
		CDKCon.CDKRouter(CDKGroup)

		// 插件管理
		PluginGroup := authAPI.Group("/plugin")
		PluginCon := controller.NewPluginController()
		PluginCon.PluginRouter(PluginGroup)
	}

	return Router
}
