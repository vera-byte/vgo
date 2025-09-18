package cmd

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/vera-byte/vgo/v"
	"github.com/vera-byte/vgo/v/cmd"
)

func init() {

	baseCmd := &gcmd.Command{
		Name:  "server",
		Usage: "server",
		Brief: "start http server",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			// g.Dump(g.DB("test").GetConfig())
			if v.IsRedisMode {
				go v.ListenFunc(ctx)
			}

			s := g.Server()

			// 如果存在 data/v-admin-vue/dist 目录，则设置为主目录
			if gfile.IsDir("frontend/dist") {
				s.SetServerRoot("frontend/dist")
			}
			s.Run()
			return nil
		},
	}
	cmd.Root.AddCommand(baseCmd)

}
