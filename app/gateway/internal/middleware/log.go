package middleware

import (
	protobuf "vgo/app/admin/api/system/v1"
	"vgo/app/gateway/internal/controller/admin_sys"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gctx"
	vck "github.com/vera-byte/vgo/vgo_core_kit"
)

func BaseSysLog(r *ghttp.Request) {
	var (
		ctx     = r.GetCtx()
		req     = g.RequestFromCtx(ctx)
		rpc     = admin_sys.NewV1()
		body    = r.GetBodyString() // 先缓存 Body 内容
		ip      = r.GetClientIp()
		traceId = gctx.CtxId(ctx)
		action  = req.Method + ":" + r.URL.Path
		admin   = vck.GetAdminAtGateway(ctx)
	)
	go func() {
		if admin == nil {
			rpc.AdminBaseSysClient.SystemLogGateway(ctx, &protobuf.SystemLogGatewayRpcInvoke{
				Action:  action,
				Params:  body, // 直接使用缓存的 body
				Ip:      ip,
				TraceId: traceId,
			})
		} else {
			rpc.AdminBaseSysClient.SystemLogGateway(ctx, &protobuf.SystemLogGatewayRpcInvoke{
				UserId:  &admin.UserId,
				Action:  action,
				Params:  body, // 直接使用缓存的 body
				Ip:      ip,
				TraceId: traceId,
			})
		}

	}()

	r.Middleware.Next()
}
