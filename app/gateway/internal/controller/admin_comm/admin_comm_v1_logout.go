package admin_comm

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"

	protobuf "github.com/vera-byte/vgo/app/admin/api/comm/v1"
	v1 "github.com/vera-byte/vgo/app/gateway/api/admin_comm/v1"
)

func (c *ControllerV1) Logout(ctx context.Context, req *v1.LogoutReq) (res *v1.LogoutRes, err error) {
	_, _err := c.AdminBaseCommClient.LoginOut(ctx, &protobuf.LoginOutRpcInvoke{})
	if _err != nil {
		g.Log().Error(ctx, _err)
		return nil, err
	}
	return nil, nil
}
