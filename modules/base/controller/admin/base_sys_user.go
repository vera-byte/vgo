package admin

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/vera-byte/vgo/modules/base/service"
	"github.com/vera-byte/vgo/v"
)

type BaseSysUserController struct {
	*v.Controller
}

func init() {
	var base_sys_user_controller = &BaseSysUserController{
		&v.Controller{
			Perfix:  "/admin/base/sys/user",
			Api:     []string{"Add", "Delete", "Update", "Info", "List", "Page", "Move"},
			Service: service.NewBaseSysUserService(),
		},
	}
	// 注册路由
	v.RegisterController(base_sys_user_controller)
}

type UserMoveReq struct {
	g.Meta        `path:"/move" method:"GET"`
	Authorization string `json:"Authorization" in:"header"`
}

func (c *BaseSysUserController) Move(ctx context.Context, req *UserMoveReq) (res *v.BaseRes, err error) {
	err = service.NewBaseSysUserService().Move(ctx)
	res = v.Ok(nil)
	return
}
