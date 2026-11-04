package main

import (
	"github.com/gogf/gf/os/gres"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/vera-byte/vgo/v/vcmd"

	_ "vgo-simple/internal/loader"
	_ "vgo-simple/modules"
)

func main() {
	gres.Dump()
	// 注册模块中的命令
	vcmd.GetRegistry().RegisterCommands()
	vcmd.Root.Name = "vgo-simple"
	vcmd.Root.Run(gctx.New())
}
