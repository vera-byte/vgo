package middleware

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/sa-tokens/sa-token-go/stputil"
	"github.com/vera-byte/vgo/modules/base/service"
)

// 本类接口无需权限验证
func BaseAuthorityMiddlewareOpen(r *ghttp.Request) {
	r.SetCtxVar("AuthOpen", true)
	r.Middleware.Next()
}

// 本类接口无需权限验证,只需登录验证
func BaseAuthorityMiddlewareComm(r *ghttp.Request) {
	r.SetCtxVar("AuthComm", true)
	r.Middleware.Next()
}

// 其余接口需登录验证同时需要权限验证
func BaseAuthorityMiddleware(r *ghttp.Request) {
	// 无需登录验证
	if r.GetCtxVar("AuthOpen", false).Bool() {
		r.Middleware.Next()
		return
	}

	// 校验登录态（refresh token 不是登录 token，无法通过校验）
	token := r.GetHeader("Authorization")
	admin, err := service.AdminFromToken(token)
	if err != nil {
		r.Response.WriteStatusExit(401, g.Map{
			"code":    1001,
			"message": "登陆失效～",
		})
	}
	// 将用户信息放入上下文
	r.SetCtxVar("admin", admin)

	// 只验证登录不验证权限的接口
	if r.GetCtxVar("AuthComm", false).Bool() {
		r.Middleware.Next()
		return
	}

	// 权限校验：请求 URL 转为 "base:comm:person" 形式与权限标识比对
	url := r.URL.String()
	if parts := gstr.Split(url, "?"); len(parts) == 2 {
		url = parts[0]
	}
	urls := gstr.Split(url, "/")
	// 去除第一个空字符串和admin
	urls = urls[2:]
	url = gstr.Join(urls, ":")
	if !stputil.HasPermission(admin.UserId, url) {
		g.Log().Errorf(r.GetCtx(), "BaseAuthorityMiddleware perms invalid userId:%d url:%s", admin.UserId, url)
		r.Response.WriteStatusExit(403, g.Map{
			"code":    1001,
			"message": "登录失效或无权限访问~",
		})
	}
	r.Middleware.Next()
}
