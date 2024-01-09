package v

import "github.com/vera-byte/vgo/cool/vconfig"

var (
	Config            = vconfig.Config            // 配置中的cool节相关配置
	GetCfgWithDefault = vconfig.GetCfgWithDefault // GetCfgWithDefault 获取配置，如果配置不存在，则使用默认值
)
