package vconfig

import (
	"context"
	"sync"

	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/frame/g"
)

var (
	// configManager 全局配置管理器实例
	configManager *ConfigManager
	// initOnce 确保配置管理器只初始化一次
	initOnce sync.Once
)

// initConfigManager 初始化配置管理器
// 使用sync.Once确保只初始化一次，支持并发安全
func initConfigManager() {
	initOnce.Do(func() {
		configManager = NewConfigManager()

		// 注册文件配置适配器作为主配置源
		fileAdapter, err := NewFileAdapter(&FileAdapterConfig{
			Path:     "manifest/config",
			FileName: "config.yaml",
			Watch:    true,
		})
		if err == nil {
			configManager.RegisterAdapter("file", fileAdapter)
			configManager.SetPrimary("file")
		}

		// 可以在这里动态注册其他配置源
		// 例如：consul、kubecm等
		// 这些配置源可以通过环境变量或配置文件进行配置
	})
}

// sConfig v框架配置结构体
// 支持从多种配置源获取配置：file、consul、kubecm等
type sConfig struct {
	AutoMigrate bool  `json:"auto_migrate,omitempty"` // 是否自动创建表
	Eps         bool  `json:"eps,omitempty"`          // 是否开启eps
	File        *file `json:"file,omitempty"`         // 文件上传配置
}

// oss OSS相关配置结构体
type oss struct {
	Endpoint        string `json:"endpoint"`        // OSS服务端点
	AccessKeyID     string `json:"accessKeyID"`     // 访问密钥ID
	SecretAccessKey string `json:"secretAccessKey"` // 访问密钥
	UseSSL          bool   `json:"useSSL"`          // 是否使用SSL
	BucketName      string `json:"bucketName"`      // 存储桶名称
	Location        string `json:"location"`        // 区域位置
}

// file 文件上传配置结构体
type file struct {
	Mode   string `json:"mode"`          // 模式 local oss
	Domain string `json:"domain"`        // 域名 http://
	Oss    *oss   `json:"oss,omitempty"` // OSS配置
}

// newConfig 创建新的配置实例
// 支持从多种配置源获取配置：file、consul、kubecm等
// 返回: 配置实例指针
func newConfig() *sConfig {
	var ctx = context.Background()
	config := &sConfig{
		AutoMigrate: GetCfgWithDefault(ctx, "v.autoMigrate", gvar.New(false)).Bool(),
		Eps:         GetCfgWithDefault(ctx, "v.eps", gvar.New(false)).Bool(),
		File: &file{
			Mode:   GetCfgWithDefault(ctx, "v.file.mode", gvar.New("none")).String(),
			Domain: GetCfgWithDefault(ctx, "v.file.domain", gvar.New("http://127.0.0.1:8300")).String(),
			Oss: &oss{
				Endpoint:        GetCfgWithDefault(ctx, "v.file.oss.endpoint", gvar.New("127.0.0.1:9000")).String(),
				AccessKeyID:     GetCfgWithDefault(ctx, "v.file.oss.accessKeyID", gvar.New("")).String(),
				SecretAccessKey: GetCfgWithDefault(ctx, "v.file.oss.secretAccessKey", gvar.New("")).String(),
				UseSSL:          GetCfgWithDefault(ctx, "v.file.oss.useSSL", gvar.New(false)).Bool(),
				BucketName:      GetCfgWithDefault(ctx, "v.file.oss.bucketName", gvar.New("vgo")).String(),
				Location:        GetCfgWithDefault(ctx, "v.file.oss.location", gvar.New("us-east-1")).String(),
			},
		},
	}
	return config
}

// Config 全局配置实例
// 支持从多种配置源获取配置：file、consul、kubecm等
var Config = newConfig()

// GetCfgWithDefault 获取配置值，支持多配置源和默认值
// ctx: 上下文
// key: 配置键
// defaultValue: 默认值
// 返回: 配置值
func GetCfgWithDefault(ctx context.Context, key string, defaultValue *gvar.Var) *gvar.Var {
	// 确保配置管理器已初始化
	initConfigManager()

	// 尝试从配置管理器获取配置
	if configManager != nil {
		if value, err := configManager.Get(ctx, key); err == nil && value != nil && !value.IsEmpty() {
			return value
		}
	}

	// 回退到GoFrame默认配置获取方式
	value, err := g.Cfg().GetWithEnv(ctx, key)
	if err != nil {
		return defaultValue
	}
	if value.IsEmpty() || value.IsNil() {
		return defaultValue
	}
	return value
}
