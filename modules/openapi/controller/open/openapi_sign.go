package app

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	v1 "github.com/vera-byte/vgo/modules/openapi/api/v1"
	"github.com/vera-byte/vgo/modules/openapi/service"
	"github.com/vera-byte/vgo/v"
)

// OpenapiSignController 开放平台签名控制器（应用端）
type OpenapiSignController struct {
	*v.ControllerSimple
}

func init() {
	var openapi_sign_controller = &OpenapiSignController{
		&v.ControllerSimple{
			Perfix: "/openapi/sign",
		},
	}
	v.RegisterControllerSimple(openapi_sign_controller)
}

// GenerateSign 生成签名
// 功能: 为应用生成数字签名
// 参数: ctx - 上下文, req - 生成签名请求
// 返回值: res - 响应结果包含生成的签名, err - 错误信息
func (c *OpenapiSignController) GenerateSign(ctx context.Context, req *v1.GenerateSignReq) (res *v.BaseRes, err error) {
	signService := service.NewOpenapiSignService()

	signReq := &service.SignRequest{
		AppId:     req.AppId,
		Timestamp: req.Timestamp,
		Nonce:     req.Nonce,
		Data:      req.Data,
	}

	// 获取客户端信息
	clientIp := g.RequestFromCtx(ctx).GetClientIp()
	userAgent := g.RequestFromCtx(ctx).Header.Get("User-Agent")

	result, err := signService.GenerateSign(ctx, signReq, clientIp, userAgent)
	if err != nil {
		return v.Fail(err.Error()), nil
	}

	return v.Ok(result), nil
}

// VerifySign 验证签名
// 功能: 验证应用提交的数字签名是否正确
// 参数: ctx - 上下文, req - 验证签名请求
// 返回值: res - 响应结果包含验证结果, err - 错误信息
func (c *OpenapiSignController) VerifySign(ctx context.Context, req *v1.AppVerifySignReq) (res *v.BaseRes, err error) {
	signService := service.NewOpenapiSignService()

	signReq := &service.SignRequest{
		AppId:     req.AppId,
		Timestamp: req.Timestamp,
		Nonce:     req.Nonce,
		Data:      req.Data,
	}

	isValid, err := signService.VerifySign(ctx, signReq, req.Signature)
	if err != nil {
		return v.Fail(err.Error()), nil
	}

	result := g.Map{
		"valid": isValid,
	}

	return v.Ok(result), nil
}

// GetPublicKey 获取公钥
// 功能: 获取指定应用的RSA公钥，用于签名验证
// 参数: ctx - 上下文, req - 获取公钥请求
// 返回值: res - 响应结果包含公钥信息, err - 错误信息
func (c *OpenapiSignController) GetPublicKey(ctx context.Context, req *v1.GetPublicKeyReq) (res *v.BaseRes, err error) {
	signService := service.NewOpenapiSignService()

	appInfo, err := signService.GetAppInfo(ctx, req.AppId)
	if err != nil {
		return v.Fail(err.Error()), nil
	}

	result := g.Map{
		"appId":     req.AppId,
		"publicKey": appInfo["public_key"],
	}

	return v.Ok(result), nil
}
