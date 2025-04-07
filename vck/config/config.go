package vck_config

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/grand"
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

type VCKConfig struct {
	IsDebug      bool // 是否调试模式
	PprofAddress string
}

func NewVCKConfig() *VCKConfig {
	var (
		ctx = context.Background()
	)
	c := g.Cfg("vck.yaml")
	isDebug, _ := c.GetWithEnv(ctx, "IsDebug", false)
	pprofAddress, _ := c.GetWithEnv(ctx, "PprofAddress", fmt.Sprintf(":%d", grand.N(9000, 9500)))

	return &VCKConfig{
		IsDebug:      isDebug.Bool(),
		PprofAddress: pprofAddress.String(),
	}
}
