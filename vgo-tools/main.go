package main

import (
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/vera-byte/vgo/vgo-tools/internal/cmd/vgocmd"
	"github.com/vera-byte/vgo/vgo-tools/utility/mlog"
)

func main() {
	var (
		ctx          = gctx.GetInitCtx()
		command, err = vgocmd.GetCommand(ctx)
	)

	if err != nil {
		mlog.Fatalf(`%+v`, err)
	}
	if command == nil {
		panic(gerror.New(`retrieve root command failed for "vgo"`))
	}
	command.Run(ctx)
}
