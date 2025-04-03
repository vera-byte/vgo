package admin_sys

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"

	"vgo/app/admin/api/pbentity"
	protobuf "vgo/app/admin/api/system/v1"
	v1 "vgo/app/gateway/api/admin_sys/v1"
)

func (c *ControllerV1) UserPage(ctx context.Context, req *v1.UserPageReq) (res *v1.UserPageRes, err error) {
	_result, err := c.AdminBaseSysClient.UserPage(ctx, &protobuf.UserPageRpcInvoke{
		DepartmentIds: req.DepartmentIds,
		PageReq: &pbentity.PageRpcInvoke{
			Page:  int32(req.Page),
			Size:  int32(req.Size),
			Sort:  req.Sort,
			Order: req.Order,
		},
	})
	if err != nil {
		return nil, err
	}
	res = &v1.UserPageRes{}
	gconv.Scan(_result.List, &res.List)
	gconv.Scan(_result.Pagination, &res.Pagination)
	g.Dump(res)
	return res, nil
}
