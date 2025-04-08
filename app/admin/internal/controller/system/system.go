package system

import (
	"context"

	v1 "github.com/vera-byte/vgo/app/admin/api/system/v1"
	"github.com/vera-byte/vgo/app/admin/internal/service"

	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/util/gconv"
	vck_request "github.com/vera-byte/vgo/vgo_core_kit/request"
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

func (*Controller) DepartmentList(ctx context.Context, req *v1.DepartmentListRpcInvoke) (res *v1.DepartmentListRpcRes, err error) {
	var (
		baseSysDepartmentService = service.BaseSysDepartmentLogic()
	)
	departments, err := baseSysDepartmentService.List(ctx)
	if err != nil {
		return
	}
	res = &v1.DepartmentListRpcRes{}
	for _, v := range departments {
		res.Items = append(res.Items, &v1.DepartmentResultItem{
			Id:         int64(v.Id),
			Name:       v.Name,
			OrderNum:   int64(v.OrderNum),
			ParentId:   int64(v.ParentId),
			TenantId:   int64(v.TenantId),
			CreateTime: v.CreatedAt.Unix(),
			UpdateTime: v.UpdatedAt.Unix(),
		})
	}
	return
}

func (*Controller) UserPage(ctx context.Context, req *v1.UserPageRpcInvoke) (res *v1.UserPageRpcRes, err error) {
	var (
		baseSysUserService = service.BaseSysUserLogic()
	)
	user, page, err := baseSysUserService.Page(ctx, &vck_request.PageReq{
		Page:  int(req.PageReq.Page),
		Size:  int(req.PageReq.Size),
		Sort:  req.PageReq.Sort,
		Order: req.PageReq.Order,
	}, req.DepartmentIds)
	if err != nil {
		return
	}
	res = &v1.UserPageRpcRes{}
	gconv.Scan(user, &res.List)
	gconv.Scan(page, &res.Pagination)
	return res, nil
}
