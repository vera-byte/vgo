package route

import (
	"context"

	"github.com/vera-byte/vgo/app/gateway/internal/controller/admin_comm"
	"github.com/vera-byte/vgo/app/gateway/internal/controller/admin_open"
	"github.com/vera-byte/vgo/app/gateway/internal/controller/admin_sys"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

func init() {
	g.Log().Info(context.Background(), "注册base模块路由")
	s := g.Server()
	/// admin-base
	s.Group("/v1/admin/base", func(admin *ghttp.RouterGroup) {
		admin.Group("/open", func(open *ghttp.RouterGroup) {
			open.Bind(admin_open.NewV1())
		})
		admin.Group("/comm", func(comm *ghttp.RouterGroup) {
			comm.Bind(admin_comm.NewV1())
		})
		admin.Group("/sys", func(comm *ghttp.RouterGroup) {
			comm.Bind(admin_sys.NewV1())
		})

	})
	s.Run()

}
