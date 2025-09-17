package config

// sConfig 开放平台模块配置结构
type sConfig struct {
	// RSA密钥配置
	RSA struct {
		DefaultKeyBits int `json:"defaultKeyBits" yaml:"defaultKeyBits"` // 默认RSA密钥位数
		MaxKeyBits     int `json:"maxKeyBits" yaml:"maxKeyBits"`         // 最大RSA密钥位数
	} `json:"rsa" yaml:"rsa"`

	// 签名配置
	Sign struct {
		ExpireTime    int  `json:"expireTime" yaml:"expireTime"`       // 签名有效期（秒）
		EnableReplay  bool `json:"enableReplay" yaml:"enableReplay"`   // 是否启用防重放攻击
		EnableLogging bool `json:"enableLogging" yaml:"enableLogging"` // 是否启用签名日志
	} `json:"sign" yaml:"sign"`

	// 应用配置
	App struct {
		MaxApps        int  `json:"maxApps" yaml:"maxApps"`               // 最大应用数量
		AutoGenSecret  bool `json:"autoGenSecret" yaml:"autoGenSecret"`   // 是否自动生成应用密钥
		SecretLength   int  `json:"secretLength" yaml:"secretLength"`     // 应用密钥长度
		EnableStatus   bool `json:"enableStatus" yaml:"enableStatus"`     // 是否启用应用状态管理
	} `json:"app" yaml:"app"`
}

// NewConfig 创建配置实例
// 功能: 创建并初始化开放平台模块配置
// 返回值: *sConfig - 配置实例
func NewConfig() *sConfig {
	config := &sConfig{}
	
	// 设置默认配置
	config.RSA.DefaultKeyBits = 2048
	config.RSA.MaxKeyBits = 4096
	
	config.Sign.ExpireTime = 300 // 5分钟
	config.Sign.EnableReplay = true
	config.Sign.EnableLogging = true
	
	config.App.MaxApps = 100
	config.App.AutoGenSecret = true
	config.App.SecretLength = 32
	config.App.EnableStatus = true
	
	return config
}

// GetRSAConfig 获取RSA配置
// 功能: 获取RSA密钥相关配置
// 返回值: RSA配置结构
func (c *sConfig) GetRSAConfig() struct {
	DefaultKeyBits int `json:"defaultKeyBits" yaml:"defaultKeyBits"`
	MaxKeyBits     int `json:"maxKeyBits" yaml:"maxKeyBits"`
} {
	return c.RSA
}

// GetSignConfig 获取签名配置
// 功能: 获取签名相关配置
// 返回值: 签名配置结构
func (c *sConfig) GetSignConfig() struct {
	ExpireTime    int  `json:"expireTime" yaml:"expireTime"`
	EnableReplay  bool `json:"enableReplay" yaml:"enableReplay"`
	EnableLogging bool `json:"enableLogging" yaml:"enableLogging"`
} {
	return c.Sign
}

// GetAppConfig 获取应用配置
// 功能: 获取应用管理相关配置
// 返回值: 应用配置结构
func (c *sConfig) GetAppConfig() struct {
	MaxApps        int  `json:"maxApps" yaml:"maxApps"`
	AutoGenSecret  bool `json:"autoGenSecret" yaml:"autoGenSecret"`
	SecretLength   int  `json:"secretLength" yaml:"secretLength"`
	EnableStatus   bool `json:"enableStatus" yaml:"enableStatus"`
} {
	return c.App
}