package admin

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	v1 "github.com/vera-byte/vgo/modules/openapi/api/v1"
	"github.com/vera-byte/vgo/modules/openapi/service"
	"github.com/vera-byte/vgo/v"
)

// OpenapiSignController 开放平台签名控制器（管理端）
type OpenapiSignController struct {
	*v.Controller
}

func init() {
	var openapi_sign_controller = &OpenapiSignController{
		&v.Controller{
			Perfix:  "/admin/openapi/sign",
			Api:     []string{"CreateApp", "GetAppInfo", "TestSign", "VerifySign"},
			Service: service.NewOpenapiSignService(),
		},
	}
	// 注册路由
	v.RegisterController(openapi_sign_controller)
}

// CreateApp 创建开放平台应用
// 功能: 管理员创建新的开放平台应用
// 参数:
//   - ctx: 上下文
//   - req: 创建应用请求参数
//
// 返回值:
//   - res: 响应结果，包含应用信息和密钥对
//   - err: 错误信息
func (c *OpenapiSignController) CreateApp(ctx context.Context, req *v1.CreateAppReq) (res *v.BaseRes, err error) {
	signService := service.NewOpenapiSignService()

	createReq := &service.CreateAppRequest{
		AppName:     req.AppName,
		Description: req.Description,
	}

	result, err := signService.CreateApp(ctx, createReq)
	if err != nil {
		return v.Fail(err.Error()), nil
	}

	return v.Ok(result), nil
}

// GetAppInfo 获取应用信息
// 功能: 管理员获取指定应用的详细信息
// 参数: ctx - 上下文, req - 获取应用信息请求
// 返回值: res - 响应结果包含应用信息, err - 错误信息
func (c *OpenapiSignController) GetAppInfo(ctx context.Context, req *v1.GetAppInfoReq) (res *v.BaseRes, err error) {
	signService := service.NewOpenapiSignService()

	result, err := signService.GetAppInfo(ctx, req.AppId)
	if err != nil {
		return v.Fail(err.Error()), nil
	}

	return v.Ok(result), nil
}

// TestSignReq 测试签名请求结构
// TestSign 测试签名生成
// 功能: 管理员测试指定应用的签名生成功能
// 参数: ctx - 上下文, req - 测试签名请求
// 返回值: res - 响应结果包含签名信息, err - 错误信息
func (c *OpenapiSignController) TestSign(ctx context.Context, req *v1.TestSignReq) (res *v.BaseRes, err error) {
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
// VerifySign 验证签名
// 功能: 管理员验证指定应用的签名正确性
// 参数: ctx - 上下文, req - 验证签名请求
// 返回值: res - 响应结果包含验证结果, err - 错误信息
func (c *OpenapiSignController) VerifySign(ctx context.Context, req *v1.VerifySignReq) (res *v.BaseRes, err error) {
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
