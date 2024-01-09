package middleware

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/vera-byte/vgo/modules/base/config"
)

func init() {
	if config.Config.Middleware.Authority.Enable {
		g.Server().BindMiddleware("/admin/*/open/*", BaseAuthorityMiddlewareOpen)
		g.Server().BindMiddleware("/admin/*/comm/*", BaseAuthorityMiddlewareComm)
		g.Server().BindMiddleware("/admin/*", BaseAuthorityMiddleware)
		g.Server().BindMiddleware("/*", AutoI18n)

	}
	if config.Config.Middleware.Log.Enable {
		g.Server().BindMiddleware("/admin/*", BaseLog)
	}

}
