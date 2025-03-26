package cmd

import (
	"context"
	_ "vgo/app/gateway/internal/middleware"

	"vgo/app/gateway/internal/controller/admin_comm"
	"vgo/app/gateway/internal/controller/admin_open"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/gogf/gf/v2/util/gconv"
	vck "github.com/vera-byte/vgo/vgo_core_kit"
)

var (
	Main = gcmd.Command{
		Name:  "main",
		Usage: "main",
		Brief: "start http gateway server",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			s := g.Server()
			/// admin-base
			s.Group("/v1/admin/base", func(admin *ghttp.RouterGroup) {
				admin.Group("/open", func(open *ghttp.RouterGroup) {
					open.Bind(admin_open.NewV1())
				})
				admin.Group("/comm", func(comm *ghttp.RouterGroup) {
					comm.Bind(admin_comm.NewV1())
				})

			})
			// admin-dict
			s.Group("/v1/admin/dict", func(admin *ghttp.RouterGroup) {

			})
			_, err = vck.EtcdManager.Client.Put(ctx, "admin", gconv.String(vck.GetAdminConfig))
			if err == nil {
				data, _ := vck.EtcdManager.GetConfig("admin")
				g.Dump(gconv.Map(data))

			}
			s.Run()
			return nil
		},
	}
)
