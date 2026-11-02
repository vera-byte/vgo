package middleware

import (
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/vera-byte/vgo/modules/user/service"
)

// userAuthPublicPaths 免登录路径（默认拒绝，登录类接口显式放行）
var userAuthPublicPaths = map[string]bool{
	"/app/user/login/mini":        true,
	"/app/user/login/mp":          true,
	"/app/user/login/wxApp":       true,
	"/app/user/login/phone":       true,
	"/app/user/login/uniPhone":    true,
	"/app/user/login/miniPhone":   true,
	"/app/user/login/captcha":     true,
	"/app/user/login/smsCode":     true,
	"/app/user/login/refreshToken": true,
	"/app/user/login/password":    true,
	"/app/user/comm/wxMpConfig":   true,
}

func init() {
	g.Server().BindMiddleware("/app/user/*", UserAuthMiddleware)
}

// UserAuthMiddleware C端用户鉴权（参照 cool-admin user 模块 UserMiddleware）
// 仅作用于 /app/user/*：校验 JWT 会话并注入当前用户；refreshToken 不可作访问令牌
func UserAuthMiddleware(r *ghttp.Request) {
	path := strings.Split(r.URL.Path, "?")[0]
	if !strings.HasPrefix(path, "/app/user/") {
		r.Middleware.Next()
		return
	}
	if userAuthPublicPaths[path] {
		r.Middleware.Next()
		return
	}

	token := r.GetHeader("Authorization")
	userId, err := service.UserLoginID(token)
	if err != nil {
		r.Response.WriteStatusExit(401, g.Map{
			"code":    1001,
			"message": "登录失效~",
		})
	}
	service.SetCurrentUser(r, userId, service.UserTenantFromSession(token))
	r.Middleware.Next()
}
