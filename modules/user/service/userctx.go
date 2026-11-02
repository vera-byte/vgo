package service

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/vera-byte/vgo/v"
)

// AppUser C端当前登录用户（由 app 鉴权中间件注入）
type AppUser struct {
	ID       uint `json:"id"`
	TenantID uint `json:"tenantId"`
}

// CurrentUser 从上下文获取当前登录用户（/app/user/* 中间件注入）
func CurrentUser(ctx context.Context) *AppUser {
	r := g.RequestFromCtx(ctx)
	user := &AppUser{}
	_ = gconv.Struct(r.GetCtxVar("user").Map(), user)
	return user
}

// SetCurrentUser 中间件注入当前登录用户
func SetCurrentUser(r *ghttp.Request, id, tenantID uint) {
	r.SetCtxVar("user", &AppUser{ID: id, TenantID: tenantID})
}

// 保持 v 包引用（DBM 扩展依赖）
var _ = v.DBM
