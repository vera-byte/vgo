package admin

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/vera-byte/vgo/cool"
	"github.com/vera-byte/vgo/modules/base/service"
)

type BaseSysUserController struct {
	*cool.Controller
}

func init() {
	var base_sys_user_controller = &BaseSysUserController{
		&cool.Controller{
			Perfix:  "/admin/base/sys/user",
			Api:     []string{"Add", "Delete", "Update", "Info", "List", "Page", "Move"},
			Service: service.NewBaseSysUserService(),
		},
	}
	// 注册路由
	cool.RegisterController(base_sys_user_controller)
}

type UserMoveReq struct {
	g.Meta        `path:"/move" method:"GET"`
	Authorization string `json:"Authorization" in:"header"`
}

func (c *BaseSysUserController) Move(ctx context.Context, req *UserMoveReq) (res *cool.BaseRes, err error) {
	err = service.NewBaseSysUserService().Move(ctx)
	res = cool.Ok(nil)
	return
}
