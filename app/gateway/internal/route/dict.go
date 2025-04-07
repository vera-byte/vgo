package route

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

func init() {
	g.Log().Info(context.Background(), "注册dict模块路由")
	s := g.Server()
	s.Group("/v1/admin/dict", func(admin *ghttp.RouterGroup) {

	})
}
