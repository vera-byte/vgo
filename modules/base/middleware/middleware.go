package middleware

import (
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/vera-byte/vgo/modules/base/config"
)

func init() {
	if config.Config.Middleware.Authority.Enable {
		g.Server().BindMiddleware("/admin/*/open/*", BaseAuthorityMiddlewareOpen)
		g.Server().BindMiddleware("/admin/*/comm/*", BaseAuthorityMiddlewareComm)
		g.Server().BindMiddleware("/admin/*", BaseAuthorityMiddleware)
		g.Server().BindMiddleware("/*", AutoI18n)

		// 敏感接口限速：登录/刷新/验证码与开放平台签名接口
		authLimit := RateLimitMiddleware(20, time.Minute)
		g.Server().BindMiddleware("/admin/*/open/login", authLimit)
		g.Server().BindMiddleware("/admin/*/open/refreshToken", authLimit)
		g.Server().BindMiddleware("/admin/*/open/captcha", authLimit)
		g.Server().BindMiddleware("/openapi/*", RateLimitMiddleware(60, time.Minute))
	}
	if config.Config.Middleware.Log.Enable {
		g.Server().BindMiddleware("/admin/*", BaseLog)
	}

}
