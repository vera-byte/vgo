package vconfig

import "github.com/gogf/gf/v2/frame/g"

// v config
type sConfig struct {
	AutoMigrate bool  `json:"auto_migrate,omitempty"` // 是否自动创建表
	Eps         bool  `json:"eps,omitempty"`          // 是否开启eps
	File        *file `json:"file,omitempty"`         // 文件上传配置
}

// OSS相关配置
type oss struct {
	Endpoint        string `json:"endpoint"`
	AccessKeyID     string `json:"accessKeyID"`
	SecretAccessKey string `json:"secretAccessKey"`
	UseSSL          bool   `json:"useSSL"`
	BucketName      string `json:"bucketName"`
	Location        string `json:"location"`
}

// 文件上传配置
type file struct {
	Mode   string `json:"mode"`   // 模式 local oss
	Domain string `json:"domain"` // 域名 http://
	Oss    *oss   `json:"oss,omitempty"`
}

// NewConfig new config
func newConfig() *sConfig {
	var ctx g.Ctx
	config := &sConfig{
		AutoMigrate: GetCfgWithDefault(ctx, "v.autoMigrate", g.NewVar(false)).Bool(),
		Eps:         GetCfgWithDefault(ctx, "v.eps", g.NewVar(false)).Bool(),
		File: &file{
			Mode:   GetCfgWithDefault(ctx, "v.file.mode", g.NewVar("none")).String(),
			Domain: GetCfgWithDefault(ctx, "v.file.domain", g.NewVar("http://127.0.0.1:8300")).String(),
			Oss: &oss{
				Endpoint:        GetCfgWithDefault(ctx, "v.file.oss.endpoint", g.NewVar("127.0.0.1:9000")).String(),
				AccessKeyID:     GetCfgWithDefault(ctx, "v.file.oss.accessKeyID", g.NewVar("")).String(),
				SecretAccessKey: GetCfgWithDefault(ctx, "v.file.oss.secretAccessKey", g.NewVar("")).String(),
				UseSSL:          GetCfgWithDefault(ctx, "v.file.oss.useSSL", g.NewVar(false)).Bool(),
				BucketName:      GetCfgWithDefault(ctx, "v.file.oss.bucketName", g.NewVar("vgo")).String(),
				Location:        GetCfgWithDefault(ctx, "v.file.oss.location", g.NewVar("us-east-1")).String(),
			},
		},
	}
	return config
}

// Config config
var Config = newConfig()

// GetCfgWithDefault get config with default value
func GetCfgWithDefault(ctx g.Ctx, key string, defaultValue *g.Var) *g.Var {
	value, err := g.Cfg().GetWithEnv(ctx, key)
	if err != nil {
		return defaultValue
	}
	if value.IsEmpty() || value.IsNil() {
		return defaultValue
	}
	return value
}
