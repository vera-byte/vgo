package base

import (
	"context"
	"vgo/app/admin/internal/dao"
	"vgo/app/admin/internal/service"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
)

func init() {
	service.RegisterBaseSysDepartmentLogic(NewBaseSysDepartmentLogic())
}

type sBaseSysDepartmentLogic struct {
}

func NewBaseSysDepartmentLogic() *sBaseSysDepartmentLogic {
	return &sBaseSysDepartmentLogic{}
}

// GetByRoleIds 获取部门
func (s *sBaseSysDepartmentLogic) GetByRoleIds(ctx context.Context, roleIds []string, isAdmin bool) (res []uint) {
	var (
		result gdb.Result
	)
	// 如果roleIds不为空
	if len(roleIds) > 0 {
		// 如果是超级管理员，则返回所有部门
		if isAdmin {
			result, _ = dao.BaseSysDepartment.Ctx(ctx).Fields("id").All()
			for _, v := range result {
				vmap := v.Map()
				if vmap["id"] != nil {
					res = append(res, gconv.Uint(vmap["id"]))
				}
			}
		} else {
			// 如果不是超级管理员，则返回角色所在部门
			result, _ = dao.BaseSysRoleDepartment.Ctx(ctx).Where("roleId IN (?)", roleIds).Fields("departmentId").All()
			for _, v := range result {
				vmap := v.Map()
				if vmap["departmentId"] != nil {
					res = append(res, gconv.Uint(vmap["departmentId"]))
				}
			}
		}

	}
	return
}

// Order 排序部门
func (s *sBaseSysDepartmentLogic) Order(ctx g.Ctx) (err error) {
	r := g.RequestFromCtx(ctx).GetMap()

	type item struct {
		Id       uint64  `json:"id"`
		ParentId *uint64 `json:"parentId,omitempty"`
		OrderNum int32   `json:"orderNum"`
	}

	var data *item

	for _, item := range r {
		err = gconv.Struct(item, &data)
		if err != nil {
			continue
		}
		dao.BaseSysDepartment.Ctx(ctx).Where("id = ?", data.Id).Data(data).Update()
	}

	return
}
