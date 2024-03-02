package main

import (
	_ "QLToolsV2/internal/packed"

	"github.com/gogf/gf/v2/os/gctx"

	"QLToolsV2/internal/cmd"
)

func main() {
	cmd.Main.Run(gctx.GetInitCtx())
}
