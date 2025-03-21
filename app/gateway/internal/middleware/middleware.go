package middleware

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	vck "github.com/vera-byte/vgo/vgo_core_kit"
)

func init() {
	s := g.Server()
	if vck.GetAdminConfig.Middleware.Authority.Enable {
		s.BindMiddleware("/*", ghttp.MiddlewareHandlerResponse)

		s.BindMiddleware("/v1/admin/*/open/*", BaseAuthorityMiddlewareOpen)
		s.BindMiddleware("/v1/admin/*/comm/*", BaseAuthorityMiddlewareComm)
		s.BindMiddleware("/v1/admin/*", BaseAuthorityMiddleware)

	}
	if vck.GetAdminConfig.Middleware.Log.Enable {
		s.BindMiddleware("/v1/admin/*", BaseLog)
	}

}
