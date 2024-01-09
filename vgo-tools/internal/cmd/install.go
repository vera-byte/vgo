package cmd

import (
	"context"

	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/vera-byte/vgo/vgo-tools/internal/service"
)

var (
	Install = gcmd.Command{
		Name:  "install",
		Usage: "vgo-tools install",
		Brief: "Install vgo-tools to the system.",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			err = service.Install.Run(ctx)
			return
		},
	}
)

// init
func init() {
	Main.AddCommand(&Install)
}
