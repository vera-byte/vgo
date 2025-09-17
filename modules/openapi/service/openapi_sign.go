package service

import (
	"context"
	"encoding/json"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/grand"
	"github.com/vera-byte/vgo/modules/openapi/model"
	"github.com/vera-byte/vgo/modules/openapi/utils"
	"github.com/vera-byte/vgo/v"
)

// OpenapiSignService 开放平台签名服务
type OpenapiSignService struct {
	*v.Service
	rsaUtil *utils.RSASignUtil
}

// SignRequest 签名请求结构
type SignRequest struct {
	AppId     string      `json:"appId" v:"required#应用ID不能为空"`
	Timestamp int64       `json:"timestamp" v:"required#时间戳不能为空"`
	Nonce     string      `json:"nonce" v:"required#随机数不能为空"`
	Data      interface{} `json:"data"`
}

// SignResponse 签名响应结构
type SignResponse struct {
	Signature string `json:"signature"`
	RequestId string `json:"requestId"`
}

// CreateAppRequest 创建应用请求结构
type CreateAppRequest struct {
	AppName     string `json:"appName" v:"required#应用名称不能为空"`
	Description string `json:"description"`
	KeyBits     int    `json:"keyBits" v:"min:2048#密钥位数不能小于2048"`
}

// CreateAppResponse 创建应用响应结构
type CreateAppResponse struct {
	AppId      string `json:"appId"`
	AppSecret  string `json:"appSecret"`
	PublicKey  string `json:"publicKey"`
	PrivateKey string `json:"privateKey"`
}

// NewOpenapiSignService 创建开放平台签名服务实例
// 功能: 创建并初始化开放平台签名服务
// 返回值: *OpenapiSignService - 签名服务实例
func NewOpenapiSignService() *OpenapiSignService {
	return &OpenapiSignService{
		Service: &v.Service{Model: &model.OpenapiApp{}},
		rsaUtil: utils.NewRSASignUtil(),
	}
}

// CreateApp 创建应用
// 功能: 创建新的开放平台应用，生成RSA密钥对和应用密钥
// 参数: ctx - 上下文, req - 创建应用请求
// 返回值: *CreateAppResponse - 创建应用响应, error - 错误信息
func (s *OpenapiSignService) CreateApp(ctx context.Context, req *CreateAppRequest) (*CreateAppResponse, error) {
	// 参数验证
	if err := g.Validator().Data(req).Run(ctx); err != nil {
		return nil, gerror.Wrap(err, "参数验证失败")
	}

	// 设置默认密钥位数
	if req.KeyBits == 0 {
		req.KeyBits = 2048
	}

	// 生成应用ID和密钥
	appId := "app_" + grand.S(16)
	appSecret := grand.S(32)

	// 生成RSA密钥对
	privateKey, publicKey, err := s.rsaUtil.GenerateRSAKeyPair(req.KeyBits)
	if err != nil {
		return nil, gerror.Wrap(err, "生成RSA密钥对失败")
	}

	// 创建应用记录
	app := &model.OpenapiApp{
		Model:       v.NewModel(),
		AppId:       appId,
		AppSecret:   appSecret,
		AppName:     req.AppName,
		Description: req.Description,
		PublicKey:   publicKey,
		PrivateKey:  privateKey,
	}

	m := v.DBM(s.Model)
	_, err = m.Insert(app)
	if err != nil {
		return nil, gerror.Wrap(err, "创建应用失败")
	}

	return &CreateAppResponse{
		AppId:      appId,
		AppSecret:  appSecret,
		PublicKey:  publicKey,
		PrivateKey: privateKey,
	}, nil
}

// GenerateSign 生成签名
// 功能: 为指定应用生成RSA签名，并记录签名日志
// 参数: ctx - 上下文, req - 签名请求, clientIp - 客户端IP, userAgent - 用户代理
// 返回值: *SignResponse - 签名响应, error - 错误信息
func (s *OpenapiSignService) GenerateSign(ctx context.Context, req *SignRequest, clientIp, userAgent string) (*SignResponse, error) {
	// 参数验证
	if err := g.Validator().Data(req).Run(ctx); err != nil {
		return nil, gerror.Wrap(err, "参数验证失败")
	}

	// 验证时间戳（防止重放攻击）
	now := time.Now().Unix()
	if s.abs(now, req.Timestamp) > 300 { // 5分钟有效期
		return nil, gerror.New("时间戳无效，请求已过期")
	}

	// 查询应用信息
	m := v.DBM(&model.OpenapiApp{})
	appRecord, err := m.Where("app_id = ? AND status = 1", req.AppId).One()
	if err != nil || appRecord.IsEmpty() {
		return nil, gerror.New("应用不存在或已禁用")
	}

	// 检查随机数是否已使用（防止重放攻击）
	logM := v.DBM(&model.OpenapiSignLog{})
	existRecord, _ := logM.Where("app_id = ? AND nonce = ? AND timestamp = ?", req.AppId, req.Nonce, req.Timestamp).One()
	if !existRecord.IsEmpty() {
		return nil, gerror.New("随机数已使用，请使用新的随机数")
	}

	// 将请求数据转换为JSON
	jsonData, err := json.Marshal(req.Data)
	if err != nil {
		return nil, gerror.Wrap(err, "请求数据序列化失败")
	}

	// 生成签名
	privateKey := appRecord["private_key"].String()
	signature, err := s.rsaUtil.GenerateSignature(string(jsonData), req.AppId, req.Timestamp, req.Nonce, privateKey)
	if err != nil {
		// 记录失败日志
		s.logSignRequest(ctx, req, string(jsonData), "", clientIp, userAgent, 0, err.Error())
		return nil, gerror.Wrap(err, "生成签名失败")
	}

	// 生成请求ID
	requestId := "req_" + grand.S(16)

	// 记录成功日志
	s.logSignRequest(ctx, req, string(jsonData), signature, clientIp, userAgent, 1, "")

	return &SignResponse{
		Signature: signature,
		RequestId: requestId,
	}, nil
}

// VerifySign 验证签名
// 功能: 验证RSA签名是否正确，包括时间戳验证和签名验证
// 参数: ctx - 上下文, req - 签名请求, signature - 待验证的签名
// 返回值: bool - 验证结果, error - 错误信息
func (s *OpenapiSignService) VerifySign(ctx context.Context, req *SignRequest, signature string) (bool, error) {
	// 参数验证
	if err := g.Validator().Data(req).Run(ctx); err != nil {
		return false, gerror.Wrap(err, "参数验证失败")
	}

	if signature == "" {
		return false, gerror.New("签名不能为空")
	}

	// 查询应用信息
	m := v.DBM(&model.OpenapiApp{})
	appRecord, err := m.Where("app_id = ? AND status = 1", req.AppId).One()
	if err != nil || appRecord.IsEmpty() {
		return false, gerror.New("应用不存在或已禁用")
	}

	// 将请求数据转换为JSON
	jsonData, err := json.Marshal(req.Data)
	if err != nil {
		return false, gerror.Wrap(err, "请求数据序列化失败")
	}

	// 验证签名
	publicKey := appRecord["public_key"].String()
	isValid, err := s.rsaUtil.VerifySignature(string(jsonData), req.AppId, req.Timestamp, req.Nonce, signature, publicKey)
	if err != nil {
		return false, gerror.Wrap(err, "验证签名失败")
	}

	return isValid, nil
}

// GetAppInfo 获取应用信息
// 功能: 根据应用ID获取应用的基本信息（不包含私钥）
// 参数: ctx - 上下文, appId - 应用ID
// 返回值: g.Map - 应用信息映射, error - 错误信息
func (s *OpenapiSignService) GetAppInfo(ctx context.Context, appId string) (g.Map, error) {
	if appId == "" {
		return nil, gerror.New("应用ID不能为空")
	}

	m := v.DBM(&model.OpenapiApp{})
	appRecord, err := m.Fields("id, app_id, app_name, description, public_key, status, create_time, update_time").
		Where("app_id = ?", appId).One()
	if err != nil || appRecord.IsEmpty() {
		return nil, gerror.New("应用不存在")
	}

	return appRecord.Map(), nil
}

// logSignRequest 记录签名请求日志
// 功能: 记录签名请求的详细信息到数据库
// 参数: ctx - 上下文, req - 签名请求, requestBody - 请求体, signature - 签名, clientIp - 客户端IP, userAgent - 用户代理, status - 状态, errorMsg - 错误信息
func (s *OpenapiSignService) logSignRequest(ctx context.Context, req *SignRequest, requestBody, signature, clientIp, userAgent string, status int32, errorMsg string) {
	log := &model.OpenapiSignLog{
		Model:       v.NewModel(),
		AppId:       req.AppId,
		RequestId:   "req_" + grand.S(16),
		Timestamp:   req.Timestamp,
		Nonce:       req.Nonce,
		RequestBody: requestBody,
		Signature:   signature,
		ClientIp:    clientIp,
		UserAgent:   userAgent,
		Status:      &status,
		ErrorMsg:    errorMsg,
	}

	// 异步记录日志，不影响主流程
	go func() {
		logM := v.DBM(&model.OpenapiSignLog{})
		if _, err := logM.Insert(log); err != nil {
			g.Log().Error(ctx, "记录签名日志失败:", err)
		}
	}()
}

// abs 计算绝对值
// 功能: 计算两个int64数值的绝对差值
// 参数: a - 第一个数值, b - 第二个数值
// 返回值: int64 - 绝对差值
func (s *OpenapiSignService) abs(a, b int64) int64 {
	if a > b {
		return a - b
	}
	return b - a
}