package app

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/vera-byte/vgo/modules/openapi/service"
	"github.com/vera-byte/vgo/v"
)

// OpenapiSignController 开放平台签名控制器（应用端）
type OpenapiSignController struct {
	*v.Controller
}

func init() {
	var openapi_sign_controller = &OpenapiSignController{
		&v.Controller{
			Perfix:  "/app/openapi/sign",
			Api:     []string{"GenerateSign", "VerifySign", "GetPublicKey"},
			Service: service.NewOpenapiSignService(),
		},
	}
	// 注册路由
	v.RegisterController(openapi_sign_controller)
}

// GenerateSignReq 生成签名请求结构
type GenerateSignReq struct {
	g.Meta    `path:"/generate" method:"POST"`
	AppId     string      `json:"appId" v:"required#应用ID不能为空"`
	Timestamp int64       `json:"timestamp" v:"required#时间戳不能为空"`
	Nonce     string      `json:"nonce" v:"required#随机数不能为空"`
	Data      interface{} `json:"data"`
}

// GenerateSign 生成签名
// 功能: 为第三方应用生成RSA签名，用于API调用认证
// 参数: ctx - 上下文, req - 生成签名请求
// 返回值: res - 响应结果包含签名信息, err - 错误信息
func (c *OpenapiSignController) GenerateSign(ctx context.Context, req *GenerateSignReq) (res *v.BaseRes, err error) {
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

// VerifySignReq 验证签名请求结构
type VerifySignReq struct {
	g.Meta    `path:"/verify" method:"POST"`
	AppId     string      `json:"appId" v:"required#应用ID不能为空"`
	Timestamp int64       `json:"timestamp" v:"required#时间戳不能为空"`
	Nonce     string      `json:"nonce" v:"required#随机数不能为空"`
	Data      interface{} `json:"data"`
	Signature string      `json:"signature" v:"required#签名不能为空"`
}

// VerifySign 验证签名
// 功能: 第三方应用验证RSA签名的正确性
// 参数: ctx - 上下文, req - 验证签名请求
// 返回值: res - 响应结果包含验证结果, err - 错误信息
func (c *OpenapiSignController) VerifySign(ctx context.Context, req *VerifySignReq) (res *v.BaseRes, err error) {
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

// GetPublicKeyReq 获取公钥请求结构
type GetPublicKeyReq struct {
	g.Meta `path:"/publicKey" method:"GET"`
	AppId  string `json:"appId" v:"required#应用ID不能为空"`
}

// GetPublicKey 获取公钥
// 功能: 第三方应用获取指定应用的RSA公钥，用于签名验证
// 参数: ctx - 上下文, req - 获取公钥请求
// 返回值: res - 响应结果包含公钥信息, err - 错误信息
func (c *OpenapiSignController) GetPublicKey(ctx context.Context, req *GetPublicKeyReq) (res *v.BaseRes, err error) {
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