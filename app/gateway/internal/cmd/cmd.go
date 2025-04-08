package cmd

import (
	"context"

	_ "github.com/vera-byte/vgo/app/gateway/internal/middleware"
	_ "github.com/vera-byte/vgo/app/gateway/internal/route"

	"github.com/gogf/gf/v2/os/gcmd"
)

var (
	Main = gcmd.Command{
		Name:  "main",
		Usage: "main",
		Brief: "start http gateway server",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {

			return nil
		},
	}
)
