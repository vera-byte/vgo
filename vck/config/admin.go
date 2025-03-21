package vck_config

import (
	"context"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/golang-jwt/jwt/v4"
	"google.golang.org/grpc/metadata"
)

type Claims struct {
	IsRefresh       bool     `json:"isRefresh"`
	RoleIds         []string `json:"roleIds"`
	Username        string   `json:"username"`
	UserId          int64    `json:"userId"`
	PasswordVersion int      `json:"passwordVersion"`
	jwt.RegisteredClaims
}

type Admin struct {
	IsRefresh       bool     `json:"isRefresh"`
	RoleIds         []string `json:"roleIds"`
	Username        string   `json:"username"`
	UserId          int64    `json:"userId"`
	PasswordVersion int      `json:"passwordVersion"`
}

type AdminConfig struct {
	Jwt        *Jwt
	Middleware *Middleware
}

type Middleware struct {
	Authority *Authority
	Log       *Log
}

type Authority struct {
	Enable bool
}

type Log struct {
	Enable bool
}

type Token struct {
	Expire        int64 `json:"expire"`
	RefreshExpire int64 `json:"refreshExprire"`
}

type Jwt struct {
	Sso    bool   `json:"sso"`
	Secret string `json:"secret"`
	Token  *Token `json:"token"`
}

// NewAdminConfig new config
func NewAdminConfig() *AdminConfig {
	var (
		ctx g.Ctx
	)
	config := &AdminConfig{
		Jwt: &Jwt{
			Sso:    GetCfgWithDefault(ctx, "modules.base.jwt.sso", g.NewVar(false)).Bool(),
			Secret: GetCfgWithDefault(ctx, "modules.base.jwt.secret", g.NewVar("")).String(),
			Token: &Token{
				Expire:        GetCfgWithDefault(ctx, "modules.base.jwt.token.expire", g.NewVar(2*3600)).Int64(),
				RefreshExpire: GetCfgWithDefault(ctx, "modules.base.jwt.token.refreshExpire", g.NewVar(15*24*3600)).Int64(),
			},
		},
		Middleware: &Middleware{
			Authority: &Authority{
				Enable: GetCfgWithDefault(ctx, "modules.base.middleware.authority.enable", g.NewVar(true)).Bool(),
			},
			Log: &Log{
				Enable: GetCfgWithDefault(ctx, "modules.base.middleware.log.enable", g.NewVar(true)).Bool(),
			},
		},
	}

	return config
}

// 获取传入ctx 中的 admin 对象
func GetAdminAtGateway(ctx context.Context) *Admin {
	r := g.RequestFromCtx(ctx)
	if r == nil {
		g.Log().Error(ctx, "获取网关上下文用户失败,当前上下文为空")
		return nil
	}
	admin := &Admin{}
	err := gjson.New(r.GetCtxVar("admin").String()).Scan(admin)
	if err != nil {
		g.Log().Error(ctx, "获取网关上下文用户失败", err)
		return nil
	}
	return admin
}

// 在grpc中获取传入admin 对象
func GetAdminAtGrpcService(ctx context.Context) *Admin {
	// ✅ 读取 metadata
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil
	}

	// ✅ 获取 Admin 字段
	adminJSON := md.Get("Admin")
	if len(adminJSON) == 0 {
		g.Log().Error(ctx, "获取网关上下文用户失败", "当前上下文为空")

		return nil
	}
	admin := &Admin{}
	gconv.Scan(adminJSON[0], &admin)
	return admin
}
