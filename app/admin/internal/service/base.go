// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"
)

type (
	ICaptcha interface {
		// 生成验证码
		GenerateCaptcha(ctx context.Context, width int, height int, color string) (id string, b64s string, answer string, err error)
		// 验证验证码
		VerifyCaptcha(id string, answer string) bool
	}
)

var (
	localCaptcha ICaptcha
)

func Captcha() ICaptcha {
	if localCaptcha == nil {
		panic("implement not found for interface ICaptcha, forgot register?")
	}
	return localCaptcha
}

func RegisterCaptcha(i ICaptcha) {
	localCaptcha = i
}
