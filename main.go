package main

import (
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gres"
	"github.com/vera-byte/vgo/internal/cmd"
	_ "github.com/vera-byte/vgo/internal/loader"
	_ "github.com/vera-byte/vgo/modules"
)

func main() {
	gres.Dump()
	cmd.GetRegistry().RegisterCommands()
	cmd.Root.Run(gctx.New())
}
