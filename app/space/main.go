package main

import (
	_ "vgo/app/space/internal/packed"

	"github.com/gogf/gf/v2/os/gctx"

	"vgo/app/space/internal/cmd"
)

func main() {
	cmd.Main.Run(gctx.GetInitCtx())
}
