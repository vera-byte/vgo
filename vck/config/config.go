package vck_config

import (
	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/frame/g"
)

// GetCfgWithDefault get config with default value
func GetCfgWithDefault(ctx g.Ctx, key string, defaultValue interface{}) *g.Var {
	value, err := g.Cfg().GetWithEnv(ctx, key)
	if err != nil {
		return gvar.New(defaultValue)
	}
	if value.IsEmpty() || value.IsNil() {
		return gvar.New(defaultValue)
	}
	return value
}
