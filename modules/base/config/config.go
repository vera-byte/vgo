package config

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/vera-byte/vgo/v"
)

// sConfig 配置
type sConfig struct {
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
	Expire        uint `json:"expire"`
	RefreshExpire uint `json:"refreshExprire"`
}

type Jwt struct {
	Sso    bool   `json:"sso"`
	Secret string `json:"secret"`
	Token  *Token `json:"token"`
}

// NewConfig new config
func NewConfig() *sConfig {
	var (
		ctx g.Ctx
	)
	config := &sConfig{
		Jwt: &Jwt{
			Sso:    v.GetCfgWithDefault(ctx, "modules.base.jwt.sso", g.NewVar(false)).Bool(),
			Secret: v.GetCfgWithDefault(ctx, "modules.base.jwt.secret", g.NewVar(v.ProcessFlag)).String(),
			Token: &Token{
				Expire:        v.GetCfgWithDefault(ctx, "modules.base.jwt.token.expire", g.NewVar(2*3600)).Uint(),
				RefreshExpire: v.GetCfgWithDefault(ctx, "modules.base.jwt.token.refreshExpire", g.NewVar(15*24*3600)).Uint(),
			},
		},
		Middleware: &Middleware{
			Authority: &Authority{
				Enable: v.GetCfgWithDefault(ctx, "modules.base.middleware.authority.enable", g.NewVar(true)).Bool(),
			},
			Log: &Log{
				Enable: v.GetCfgWithDefault(ctx, "modules.base.middleware.log.enable", g.NewVar(true)).Bool(),
			},
		},
	}

	return config
}

// Config config
var Config = NewConfig()
