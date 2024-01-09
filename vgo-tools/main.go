package main

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcfg"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/vera-byte/vgo/vgo-tools/internal/cmd"
	_ "github.com/vera-byte/vgo/vgo-tools/internal/packed"
	"github.com/vera-byte/vgo/vgo-tools/internal/utility/allyes"
	"github.com/vera-byte/vgo/vgo-tools/internal/utility/mlog"
)

const (
	cliFolderName = `hack`
)

func main() {
	// gres.Dump()
	// CLI configuration.
	if path, _ := gfile.Search(cliFolderName); path != "" {
		if adapter, ok := g.Cfg().GetAdapter().(*gcfg.AdapterFile); ok {
			if err := adapter.SetPath(path); err != nil {
				mlog.Fatal(err)
			}
		}
	}
	// -y option checks.
	allyes.Init()

	err := cmd.Main.RunWithError(gctx.New())
	if err != nil {
		println(err.Error())
	}
}
