package middleware

import "github.com/gogf/gf/v2/frame/g"

func init() {
	// 添加签名验证中间件
	// group.Middleware(middleware.OpenapiSignVerifyMiddleware)
	// 添加响应处理中间件
	// group.Middleware(v.MiddlewareHandlerResponse)

	// 获取公钥接口无需签名验证

	g.Server().BindMiddleware("/openapi/*", OpenapiSignVerifyMiddleware)

}
