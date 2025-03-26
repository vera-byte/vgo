package admin_sys

import (
	"context"

	v1 "vgo/app/gateway/api/admin_sys/v1"
)

func (c *ControllerV1) LogPage(ctx context.Context, req *v1.LogPageReq) (res *v1.LogPageRes, err error) {
	// c.AdminBaseSysClient.SystemLogPage(ctx, &protobuf.SystemLogPageReq{})
	return
}
