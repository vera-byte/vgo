package admin

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
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

// CreateAppReq 创建应用请求结构
type CreateAppReq struct {
	g.Meta      `path:"/createApp" method:"POST"`
	AppName     string `json:"appName" v:"required#应用名称不能为空"`
	Description string `json:"description"`
	KeyBits     int    `json:"keyBits" v:"min:2048#密钥位数不能小于2048"`
}

// CreateApp 创建开放平台应用
// 功能: 管理员创建新的开放平台应用
// 参数:
//   - ctx: 上下文
//   - req: 创建应用请求参数
// 返回值:
//   - res: 响应结果，包含应用信息和密钥对
//   - err: 错误信息
func (c *OpenapiSignController) CreateApp(ctx context.Context, req *CreateAppReq) (res *v.BaseRes, err error) {
	signService := service.NewOpenapiSignService()
	
	createReq := &service.CreateAppRequest{
		AppName:     req.AppName,
		Description: req.Description,
		KeyBits:     req.KeyBits,
	}

	result, err := signService.CreateApp(ctx, createReq)
	if err != nil {
		return v.Fail(err.Error()), nil
	}

	return v.Ok(result), nil
}

// GetAppInfoReq 获取应用信息请求结构
type GetAppInfoReq struct {
	g.Meta `path:"/getAppInfo" method:"GET"`
	AppId  string `json:"appId" v:"required#应用ID不能为空"`
}

// GetAppInfo 获取应用信息
// 功能: 管理员获取指定应用的详细信息
// 参数: ctx - 上下文, req - 获取应用信息请求
// 返回值: res - 响应结果包含应用信息, err - 错误信息
func (c *OpenapiSignController) GetAppInfo(ctx context.Context, req *GetAppInfoReq) (res *v.BaseRes, err error) {
	signService := service.NewOpenapiSignService()

	result, err := signService.GetAppInfo(ctx, req.AppId)
	if err != nil {
		return v.Fail(err.Error()), nil
	}

	return v.Ok(result), nil
}

// TestSignReq 测试签名请求结构
type TestSignReq struct {
	g.Meta    `path:"/testSign" method:"POST"`
	AppId     string      `json:"appId" v:"required#应用ID不能为空"`
	Timestamp int64       `json:"timestamp" v:"required#时间戳不能为空"`
	Nonce     string      `json:"nonce" v:"required#随机数不能为空"`
	Data      interface{} `json:"data"`
}

// TestSign 测试签名生成
// 功能: 管理员测试指定应用的签名生成功能
// 参数: ctx - 上下文, req - 测试签名请求
// 返回值: res - 响应结果包含生成的签名, err - 错误信息
func (c *OpenapiSignController) TestSign(ctx context.Context, req *TestSignReq) (res *v.BaseRes, err error) {
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
	g.Meta    `path:"/verifySign" method:"POST"`
	AppId     string      `json:"appId" v:"required#应用ID不能为空"`
	Timestamp int64       `json:"timestamp" v:"required#时间戳不能为空"`
	Nonce     string      `json:"nonce" v:"required#随机数不能为空"`
	Data      interface{} `json:"data"`
	Signature string      `json:"signature" v:"required#签名不能为空"`
}

// VerifySign 验证签名
// 功能: 管理员验证指定签名是否正确
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