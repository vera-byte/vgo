package middleware

import (
	"fmt"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/vera-byte/vgo/v"
)

// RateLimitMiddleware 按客户端 IP 限速（时间窗口内计数），用于登录/刷新等敏感接口。
func RateLimitMiddleware(limit int, window time.Duration) ghttp.HandlerFunc {
	return func(r *ghttp.Request) {
		ctx := r.GetCtx()
		bucket := time.Now().Unix() / int64(window.Seconds())
		key := fmt.Sprintf("ratelimit:%s:%s:%d", r.URL.Path, r.GetClientIp(), bucket)
		val, _ := v.CacheManager.Get(ctx, key)
		count := val.Int64() + 1
		v.CacheManager.Set(ctx, key, count, window)
		if count > int64(limit) {
			r.Response.WriteStatusExit(429, g.Map{
				"code":    1001,
				"message": "请求过于频繁，请稍后再试",
			})
		}
		r.Middleware.Next()
	}
}
