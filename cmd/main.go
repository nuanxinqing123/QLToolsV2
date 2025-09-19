// Package main QLToolsV2
// @title QLToolsV2
// @version 1.0
// @description 青龙ToolsV2
// @termsOfService https://swagger.io/terms/

// @contact.name API Support
// @contact.url https://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url https://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:1500
// @BasePath /api

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @description Bearer token for API authentication
package main

import "github.com/nuanxinqing123/QLToolsV2/internal/app"

func main() {
	app.Start()
}
