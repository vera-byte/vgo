package middleware

import (
	"github.com/gogf/gf/v2/net/ghttp"
)

func BaseLog(r *ghttp.Request) {
	var (
	// ctx = r.GetCtx()
	// BaseSysLogService = service.NewBaseSysLogService()
	)
	// BaseSysLogService.Record(ctx)
	r.Middleware.Next()
}
