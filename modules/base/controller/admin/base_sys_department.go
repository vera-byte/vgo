package admin

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/vera-byte/vgo/modules/base/service"
	"github.com/vera-byte/vgo/v"
)

type BaseSysDepartmentController struct {
	*v.Controller
}

func init() {
	var base_sys_department_controller = &BaseSysDepartmentController{
		&v.Controller{
			Perfix:  "/admin/base/sys/department",
			Api:     []string{"Add", "Delete", "Update", "Info", "List", "Page"},
			Service: service.NewBaseSysDepartmentService(),
		},
	}
	// 注册路由
	v.RegisterController(base_sys_department_controller)
}

// OrderReq 接口请求参数
type OrderReq struct {
	g.Meta        `path:"/order" method:"GET"`
	Authorization string `json:"Authorization" in:"header"`
}

// Order 排序部门
func (c *BaseSysDepartmentController) Order(ctx context.Context, req *OrderReq) (res *v.BaseRes, err error) {
	err = service.NewBaseSysDepartmentService().Order(ctx)
	res = v.Ok(nil)
	return
}
