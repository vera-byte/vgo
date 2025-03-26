package system

import (
	"context"
	v1 "vgo/app/admin/api/system/v1"
	"vgo/app/admin/internal/service"

	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
)

type Controller struct {
	v1.UnimplementedBaseSysServer
}

func Register(s *grpcx.GrpcServer) {
	v1.RegisterBaseSysServer(s.Server, &Controller{})
}

func (*Controller) SystemLogPage(ctx context.Context, req *v1.SystemLogPageRpcInvoke) (res *v1.SystemLogPageRpcRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}

func (*Controller) SystemLogGateway(ctx context.Context, req *v1.SystemLogGatewayRpcInvoke) (res *v1.SystemLogGatewayRpcRes, err error) {
	var (
		baseSysLogService = service.BaseSysLogLogic()
	)

	baseSysLogService.RecordLog(ctx,
		req.UserId,
		req.Action,
		req.Ip,
		req.Params,
		req.TenantId,
		req.TraceId,
	)

	return
}
