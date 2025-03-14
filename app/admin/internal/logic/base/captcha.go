package base

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/mojocn/base64Captcha"
)

type sCaptcha struct {
	store base64Captcha.Store
}

var Captcha = &sCaptcha{
	store: base64Captcha.DefaultMemStore,
}

// 生成验证码
func (c *sCaptcha) GenerateCaptcha(ctx context.Context, width int, height int, color string) (id string, b64s string, answer string, err error) {
	g.Log().Debug(ctx, "生成验证码", width, height, color)
	driver := base64Captcha.NewDriverDigit(height, width, 4, 0.7, 80)
	captcha := base64Captcha.NewCaptcha(driver, c.store)
	id, b64s, answer, err = captcha.Generate()
	if err != nil {
		return "", "", "", err
	}
	return
}

// 验证验证码
func (c *sCaptcha) VerifyCaptcha(id, answer string) bool {
	return c.store.Verify(id, answer, true)
}
