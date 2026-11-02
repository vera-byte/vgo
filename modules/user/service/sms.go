package service

import (
	"context"
	"fmt"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/grand"
	pluginService "github.com/vera-byte/vgo/modules/plugin/service"
	"github.com/vera-byte/vgo/modules/user/config"
	"github.com/vera-byte/vgo/v"
)

// UserSmsService 短信验证码（参照 cool-admin user 模块 UserSmsService）
// 发送经插件分发（sms-tx / sms-ali），调试模式下验证码写日志便于联调
type UserSmsService struct{}

// NewUserSmsService 创建短信服务实例
func NewUserSmsService() *UserSmsService {
	return &UserSmsService{}
}

// SendSms 发送短信验证码（随机四位，缓存 sms:手机号，有效期 Sms.Timeout 秒）
func (s *UserSmsService) SendSms(ctx context.Context, phone string) error {
	code := grand.Digits(4)
	sent := false
	var lastErr error
	for _, key := range []string{"sms-tx", "sms-ali"} {
		if _, err := pluginService.Invoke(ctx, key, "send", []string{phone}, code); err == nil {
			sent = true
			break
		} else {
			lastErr = err
		}
	}
	if !sent {
		if config.Config.Sms.Debug {
			g.Log().Warningf(ctx, "[短信调试] 手机号 %s 验证码: %s（未配置短信插件）", phone, code)
		} else {
			return gerror.Wrap(lastErr, "未配置短信插件，请安装短信插件（如 sms-tx / sms-ali）")
		}
	}
	ttl := time.Duration(config.Config.Sms.Timeout) * time.Second
	return v.CacheManager.Set(ctx, fmt.Sprintf("sms:%s", phone), code, ttl)
}

// CheckCode 校验短信验证码（通过后失效）
func (s *UserSmsService) CheckCode(ctx context.Context, phone, code string) bool {
	key := fmt.Sprintf("sms:%s", phone)
	val, err := v.CacheManager.Get(ctx, key)
	if err != nil || val.IsNil() || val.String() != code {
		return false
	}
	v.CacheManager.Remove(ctx, key)
	return true
}
