package web

import (
	"embed"
)

// DistFS 嵌入前端构建产物
//
//go:embed all:dist
var DistFS embed.FS

// AssetsFS 嵌入前端构建产物的 assets 目录
//
//go:embed all:dist/assets
var AssetsFS embed.FS
