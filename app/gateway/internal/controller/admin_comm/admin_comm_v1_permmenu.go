package admin_comm

import (
	"context"

	protobuf "vgo/app/admin/api/comm/v1"
	v1 "vgo/app/gateway/api/admin_comm/v1"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/util/gconv"
)

func (c *ControllerV1) Permmenu(ctx context.Context, req *v1.PermmenuReq) (res *v1.PermmenuRes, err error) {

	_res, _err := c.AdminBaseCommClient.Permmenu(ctx, &protobuf.PermmenuRpcInvoke{})
	if _err != nil {
		return nil, gerror.New("获取权限菜单失败")
	}
	_menus := []v1.Menu{}
	gconv.Scan(&_res.Menus, &_menus)
	return &v1.PermmenuRes{
		Menus: _menus,
		Perms: _res.Perms,
	}, nil
}
