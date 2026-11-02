package config

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/vera-byte/vgo/v"
)

// Token C端用户会话配置（sa-token 会话时长；旧键名 modules.user.jwt.* 仍兼容）
type Token struct {
	Expire        uint `json:"expire"`        // token 过期时间（秒）
	RefreshExpire uint `json:"refreshExpire"` // 刷新 token 有效期提示（秒）
}

// Sms 短信验证码配置（参照 module.user.sms）
type Sms struct {
	Timeout uint `json:"timeout"` // 验证码有效期（秒）
	Debug   bool `json:"debug"`   // 调试模式：无短信插件时验证码写日志
}

// Wx 微信开放平台配置
type Wx struct {
	AppId     string `json:"appId"`
	AppSecret string `json:"appSecret"`
}

// UserConfig 用户模块配置
type UserConfig struct {
	Token *Token
	Sms   *Sms
	Wx    *Wx
}

func newUserConfig() *UserConfig {
	ctx := context.Background()
	// 优先 modules.user.token.*，兼容旧 modules.user.jwt.*
	expandUint := func(key, legacy string, def uint) uint {
		if val := v.GetCfgWithDefault(ctx, key, g.NewVar(0)).Uint(); val != 0 {
			return val
		}
		return v.GetCfgWithDefault(ctx, legacy, g.NewVar(def)).Uint()
	}
	return &UserConfig{
		Token: &Token{
			Expire:        expandUint("modules.user.token.expire", "modules.user.jwt.expire", 24*3600),
			RefreshExpire: expandUint("modules.user.token.refreshExpire", "modules.user.jwt.refreshExpire", 30*24*3600),
		},
		Sms: &Sms{
			Timeout: v.GetCfgWithDefault(ctx, "modules.user.sms.timeout", g.NewVar(180)).Uint(),
			Debug:   v.GetCfgWithDefault(ctx, "modules.user.sms.debug", g.NewVar(false)).Bool(),
		},
		Wx: &Wx{
			AppId:     v.GetCfgWithDefault(ctx, "modules.user.wx.appId", g.NewVar("")).String(),
			AppSecret: v.GetCfgWithDefault(ctx, "modules.user.wx.appSecret", g.NewVar("")).String(),
		},
	}
}

// Config 用户模块全局配置
var Config = newUserConfig()
