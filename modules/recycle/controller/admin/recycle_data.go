package admin

import (
	"github.com/gogf/gf/v2/frame/g"
	v1 "github.com/vera-byte/vgo/modules/recycle/api/v1"
	"github.com/vera-byte/vgo/modules/recycle/service"
	"github.com/vera-byte/vgo/v"
)

// RecycleDataController 数据回收控制器
type RecycleDataController struct {
	*v.Controller
}

func init() {
	c := &RecycleDataController{
		&v.Controller{
			Prefix:  "/admin/recycle/data",
			Api:     []string{"Add", "Delete", "Update", "Info", "List", "Page"},
			Service: service.NewRecycleDataService(),
		},
	}
	v.RegisterController(c)
}

// Restore 恢复数据
func (c *RecycleDataController) Restore(ctx g.Ctx, req *v1.RecycleDataRestoreReq) (res *v.BaseRes, err error) {
	err = service.NewRecycleDataService().Restore(ctx, req.Ids)
	if err != nil {
		return
	}
	res = v.Ok(nil)
	return
}

// Clear 清理回收站
func (c *RecycleDataController) Clear(ctx g.Ctx, req *v1.RecycleDataClearReq) (res *v.BaseRes, err error) {
	err = service.NewRecycleDataService().Clear(ctx)
	if err != nil {
		return
	}
	res = v.Ok(nil)
	return
}
