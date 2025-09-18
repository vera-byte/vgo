package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"strconv"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/vera-byte/vgo/modules/openapi/service"
)

// OpenapiSignVerifyMiddleware openapi签名验证中间件
// 功能: 验证openapi请求的数字签名，确保请求的安全性和完整性
// 参数: r - HTTP请求对象
// 返回值: 无返回值，验证失败时直接返回错误响应
func OpenapiSignVerifyMiddleware(r *ghttp.Request) {
	var (
		ctx = r.GetCtx()
	)

	// 检查是否为开放接口（无需签名验证）
	signOpen := r.GetCtxVar("SignOpen", false)
	if signOpen.Bool() {
		r.Middleware.Next()
		return
	}

	// 获取请求头中的签名信息
	appId := r.Header.Get("X-App-Id")
	timestamp := r.Header.Get("X-Timestamp")
	nonce := r.Header.Get("X-Nonce")
	signature := r.Header.Get("X-Signature")

	// 验证必要的签名参数
	if appId == "" {
		r.Response.WriteStatusExit(400, g.Map{
			"code":    4001,
			"message": "缺少应用ID1",
		})
		return
	}

	if timestamp == "" {
		r.Response.WriteStatusExit(400, g.Map{
			"code":    4002,
			"message": "缺少时间戳",
		})
		return
	}

	if nonce == "" {
		r.Response.WriteStatusExit(400, g.Map{
			"code":    4003,
			"message": "缺少随机数",
		})
		return
	}

	if signature == "" {
		r.Response.WriteStatusExit(400, g.Map{
			"code":    4004,
			"message": "缺少签名",
		})
		return
	}

	// 验证时间戳格式和有效性
	timestampInt, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		r.Response.WriteStatusExit(400, g.Map{
			"code":    4005,
			"message": "时间戳格式错误",
		})
		return
	}

	// 验证时间戳是否在有效范围内（5分钟内）
	currentTime := time.Now().Unix()
	if abs(currentTime, timestampInt) > 300 { // 5分钟 = 300秒
		r.Response.WriteStatusExit(400, g.Map{
			"code":    4006,
			"message": "请求已过期",
		})
		return
	}

	// 读取请求体数据
	var requestData interface{}
	if r.Method == "POST" || r.Method == "PUT" || r.Method == "PATCH" {
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			r.Response.WriteStatusExit(400, g.Map{
				"code":    4007,
				"message": "读取请求体失败",
			})
			return
		}

		// 重新设置请求体，以便后续处理可以正常读取
		r.Body = io.NopCloser(bytes.NewReader(bodyBytes))

		// 解析JSON数据
		if len(bodyBytes) > 0 {
			if err := json.Unmarshal(bodyBytes, &requestData); err != nil {
				// 如果不是JSON格式，使用原始字符串
				requestData = string(bodyBytes)
			}
		}
	}

	// 构建签名验证请求
	signService := service.NewOpenapiSignService()
	signReq := &service.SignRequest{
		AppId:     appId,
		Timestamp: timestampInt,
		Nonce:     nonce,
		Data:      requestData,
	}

	// 验证签名
	isValid, err := signService.VerifySign(ctx, signReq, signature)
	if err != nil {
		g.Log().Error(ctx, "签名验证失败:", err)
		r.Response.WriteStatusExit(401, g.Map{
			"code":    4008,
			"message": "签名验证失败: " + err.Error(),
		})
		return
	}

	if !isValid {
		r.Response.WriteStatusExit(401, g.Map{
			"code":    4009,
			"message": "签名验证不通过",
		})
		return
	}

	// 将应用ID存入上下文，供后续使用
	r.SetCtxVar("appId", appId)

	// 继续处理请求
	r.Middleware.Next()
}

// OpenapiSignVerifyMiddlewareOpen 开放接口中间件
// 功能: 标记接口为开放接口，无需签名验证
// 参数: r - HTTP请求对象
// 返回值: 无返回值
func OpenapiSignVerifyMiddlewareOpen(r *ghttp.Request) {
	r.SetCtxVar("SignOpen", true)
	r.Middleware.Next()
}

// abs 计算绝对值
// 功能: 计算两个int64数值的绝对值差
// 参数: a, b - 两个int64数值
// 返回值: int64 - 绝对值差
func abs(a, b int64) int64 {
	if a > b {
		return a - b
	}
	return b - a
}
