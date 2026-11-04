package main

import (
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gres"
	"github.com/vera-byte/vgo/v/vcmd"
	_ "github.com/vera-byte/vgo/internal/loader"
	_ "github.com/vera-byte/vgo/modules"
)

func main() {
	gres.Dump()
	// 注册模块中的命令
	vcmd.GetRegistry().RegisterCommands()
	vcmd.Root.Run(gctx.New())
}
