package vck_config

import (
	"context"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
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
		ctx    g.Ctx
		config *AdminConfig
		etcd   = NewChainableEtcdClient()
	)
	if etcd != nil {
		g.Log().Warning(ctx, "Admin配置为分布式配置")
		adminConfig, err := etcd.GetConfig("admin")
		if err != nil {
			err = gerror.New("当前实例未配置admin配置")
			panic(err)
		}
		adminConfig.Scan(&config)
	} else {
		g.Log().Warning(ctx, "Admin配置当前为单接点配置,如分布式模式请配置etcd")
		config = GetAdminConfigAtFile()
	}

	return config
}

func GetAdminConfigAtFile() *AdminConfig {
	var (
		ctx    g.Ctx
		config *AdminConfig
	)
	config = &AdminConfig{
		Jwt: &Jwt{
			Sso:    GetCfgWithDefault(ctx, "admin.jwt.sso", false).Bool(),
			Secret: GetCfgWithDefault(ctx, "admin.jwt.secret", "").String(),
			Token: &Token{
				Expire:        GetCfgWithDefault(ctx, "admin.jwt.token.expire", 2*3600).Int64(),
				RefreshExpire: GetCfgWithDefault(ctx, "admin.jwt.token.refreshExpire", 15*24*3600).Int64(),
			},
		},
		Middleware: &Middleware{
			Authority: &Authority{
				Enable: GetCfgWithDefault(ctx, "admin.middleware.authority.enable", true).Bool(),
			},
			Log: &Log{
				Enable: GetCfgWithDefault(ctx, "admin.middleware.log.enable", true).Bool(),
			},
		},
	}
	return config
}

func putAdminAtEtcd() {
	var (
		config = GetAdminConfigAtFile()
		etcd   = NewChainableEtcdClient()
	)

	// 执行 Put 操作
	err := etcd.PutConfig("admin", config)
	if err != nil {
		// 处理 etcd 操作错误
		g.Log().Errorf(context.Background(), "Failed to put data in etcd: %v", err)
	}
	g.Dump(NewAdminConfig())
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
	gatewayAdminToken := md.Get("AdminToken")
	if len(gatewayAdminToken) == 0 {
		g.Log().Info(ctx, "当前上下文为空")
		return nil
	}
	claims, err := ParseToken(ctx, gatewayAdminToken[0])
	if err != nil {
		g.Log().Error(ctx, "获取网关上下文用户失败", err)
		return nil
	}
	admin := &Admin{
		IsRefresh:       claims.IsRefresh,
		RoleIds:         claims.RoleIds,
		Username:        claims.Username,
		UserId:          claims.UserId,
		PasswordVersion: claims.PasswordVersion,
	}

	return admin
}

// / 解析token
func ParseToken(ctx context.Context, tokenString string) (claims *Claims, err error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(NewAdminConfig().Jwt.Secret), nil
	})
	if err != nil {
		g.Log().Error(ctx, "BaseAuthorityMiddleware", err)
		return nil, err
	}
	if !token.Valid {
		g.Log().Error(ctx, "BaseAuthorityMiddleware", "token invalid")
		return nil, gerror.New("token invalid")
	}
	claims = token.Claims.(*Claims)
	return
}
