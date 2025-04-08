package base

import (
	"context"

	"github.com/vera-byte/vgo/app/admin/internal/dao"
	"github.com/vera-byte/vgo/app/admin/internal/model/entity"
	"github.com/vera-byte/vgo/app/admin/internal/service"

	"github.com/gogf/gf/v2/errors/gerror"
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
// GetByRoleIds 根据角色 ID 获取部门 ID 列表
// 如果 isAdmin 为 true，则返回所有部门 ID（超级管理员权限）
// 如果 isAdmin 为 false，则返回该角色所属的部门 ID
func (s *sBaseSysDepartmentLogic) GetByRoleIds(ctx context.Context, roleIds []string, isAdmin bool) (res []uint) {
	// 如果 roleIds 为空，直接返回空结果
	if len(roleIds) == 0 {
		return
	}

	// // 获取数据库表名
	// departmentTable := dao.BaseSysDepartment.Table()
	// roleDepartmentTable := dao.BaseSysRoleDepartment.Table()

	// 获取字段名
	departmentIdField := dao.BaseSysDepartment.Columns().Id
	roleDepartmentRoleIdField := dao.BaseSysRoleDepartment.Columns().RoleId
	roleDepartmentDepartmentIdField := dao.BaseSysRoleDepartment.Columns().DepartmentId

	// 定义查询结果
	var result []map[string]interface{}

	if isAdmin {
		// 如果是超级管理员，获取所有部门 ID
		_ = dao.BaseSysDepartment.Ctx(ctx).
			Fields(departmentIdField).
			Scan(&result)
	} else {
		// 角色对应的部门 ID 查询
		_ = dao.BaseSysRoleDepartment.Ctx(ctx).
			Where(roleDepartmentRoleIdField+" IN(?)", roleIds).
			Fields(roleDepartmentDepartmentIdField).
			Scan(&result)
	}

	// 转换查询结果
	for _, v := range result {
		if id, ok := v[departmentIdField]; ok {
			res = append(res, gconv.Uint(id))
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
		dao.BaseSysDepartment.Ctx(ctx).Where(dao.BaseSysDepartment.Columns().Id, data.Id).Data(data).Update()
	}

	return
}

// 部门列表
func (s *sBaseSysDepartmentLogic) List(ctx context.Context) (res []entity.BaseSysDepartment, err error) {
	err = dao.BaseSysDepartment.Ctx(ctx).Scan(&res)
	if err != nil {
		_err := gerror.New("部门列表失败")
		g.Log().Error(ctx, _err.Error(), err)
		return nil, _err
	}
	return
}
