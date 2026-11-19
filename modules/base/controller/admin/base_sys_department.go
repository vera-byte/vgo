package admin

import (
	"context"

	v1 "github.com/vera-byte/vgo/modules/base/api/v1"
	"github.com/vera-byte/vgo/modules/base/service"
	"github.com/vera-byte/vgo/v"
)

type BaseSysDepartmentController struct {
	*v.Controller
}

func init() {
	var base_sys_department_controller = &BaseSysDepartmentController{
		&v.Controller{
			Prefix:  "/admin/base/sys/department",
			Api:     []string{"Add", "Delete", "Update", "Info", "List", "Page"},
			Service: service.NewBaseSysDepartmentService(),
		},
	}
	// 注册路由
	v.RegisterController(base_sys_department_controller)
}

// Order 排序部门
func (c *BaseSysDepartmentController) Order(ctx context.Context, req *v1.OrderReq) (res *v.BaseRes, err error) {
	err = service.NewBaseSysDepartmentService().Order(ctx, req.Items)
	res = v.Ok(nil)
	return
}
