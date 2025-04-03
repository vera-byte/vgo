package cmd

import (
	"context"
	_ "vgo/app/gateway/internal/middleware"
	_ "vgo/app/gateway/internal/route"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcmd"
)

var (
	Main = gcmd.Command{
		Name:  "main",
		Usage: "main",
		Brief: "start http gateway server",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			s := g.Server()
			// /// admin-base
			// s.Group("/v1/admin/base", func(admin *ghttp.RouterGroup) {
			// 	admin.Group("/open", func(open *ghttp.RouterGroup) {
			// 		open.Bind(admin_open.NewV1())
			// 	})
			// 	admin.Group("/comm", func(comm *ghttp.RouterGroup) {
			// 		comm.Bind(admin_comm.NewV1())
			// 	})
			// 	admin.Group("/sys", func(comm *ghttp.RouterGroup) {
			// 		comm.Bind(admin_sys.NewV1())
			// 	})

			// })
			// // admin-dict
			// s.Group("/v1/admin/dict", func(admin *ghttp.RouterGroup) {

			// })

			s.Run()
			return nil
		},
	}
)
