package service

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	v1 "github.com/vera-byte/vgo/modules/base/api/v1"
	"github.com/vera-byte/vgo/modules/base/model"
	"github.com/vera-byte/vgo/v"
)

type BaseSysDepartmentService struct {
	*v.Service
}

// GetByRoleIds 获取部门
func (s *BaseSysDepartmentService) GetByRoleIds(roleIds []string, isAdmin bool) (res []uint, err error) {
	var (
		result                gdb.Result
		BaseSysRoleDepartment = model.NewBaseSysRoleDepartment()
	)
	// 如果roleIds不为空
	if len(roleIds) > 0 {
		// 如果是超级管理员，则返回所有部门
		if isAdmin {
			result, err = v.DBM(s.Model).Fields("id").All()
			if err != nil {
				return nil, err
			}
			for _, v := range result {
				vmap := v.Map()
				if vmap["id"] != nil {
					res = append(res, gconv.Uint(vmap["id"]))
				}
			}
		} else {
			// 如果不是超级管理员，则返回角色所在部门
			result, err = v.DBM(BaseSysRoleDepartment).WhereIn("roleId", roleIds).Fields("departmentId").All()
			if err != nil {
				return nil, err
			}
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

// Order 批量更新部门排序与层级（整体原子生效）
func (s *BaseSysDepartmentService) Order(ctx context.Context, items []v1.OrderItem) (err error) {
	return g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) (err error) {
		for _, item := range items {
			if item.Id == 0 {
				continue
			}
			if _, err = v.DBM(s.Model).TX(tx).Where("id = ?", item.Id).Data(g.Map{
				"parentId": item.ParentId,
				"orderNum": item.OrderNum,
			}).Update(); err != nil {
				return err
			}
		}
		return
	})
}

// NewBaseSysDepartmentService 创建一个BaseSysDepartmentService实例
func NewBaseSysDepartmentService() *BaseSysDepartmentService {
	return &BaseSysDepartmentService{
		Service: &v.Service{
			Model:       model.NewBaseSysDepartment(),
			ListQueryOp: &v.QueryOp{},
		},
	}
}
