package admin_sys

import (
	"context"

	"github.com/gogf/gf/v2/util/gconv"

	protobuf "github.com/vera-byte/vgo/app/admin/api/system/v1"
	v1 "github.com/vera-byte/vgo/app/gateway/api/admin_sys/v1"
)

func (c *ControllerV1) DepartmentList(ctx context.Context, req *v1.DepartmentListReq) (res *v1.DepartmentListRes, err error) {

	result, err := c.AdminBaseSysClient.DepartmentList(ctx, &protobuf.DepartmentListRpcInvoke{})
	if err != nil {
		return nil, err
	}
	gconv.Scan(result.Items, &res)
	return
}
