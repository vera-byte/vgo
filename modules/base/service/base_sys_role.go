package service

import (
	"context"

	"github.com/gogf/gf/v2/container/garray"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/vera-byte/vgo/modules/base/model"
	"github.com/vera-byte/vgo/v"
)

var baseSysUserRole = model.NewBaseSysUserRole()

type BaseSysRoleService struct {
	*v.Service
}

// ModifyAfter modify after
func (s *BaseSysRoleService) ModifyAfter(ctx context.Context, method string, param g.MapStrAny) (err error) {
	if param["id"] != nil {
		err = s.updatePerms(ctx, gconv.Uint(param["id"]), gconv.SliceUint(param["menuIdList"]), gconv.SliceUint(param["departmentIdList"]))
	}
	return
}

// updatePerms(roleId, menuIdList?, departmentIds = [])
func (s *BaseSysRoleService) updatePerms(ctx context.Context, roleId uint, menuIdList, departmentIds []uint) (err error) {
	// 更新菜单权限
	v.DBM(model.NewBaseSysRoleMenu()).Where("roleId", roleId).Delete()
	if len(menuIdList) > 0 {
		roleMenuList := make([]g.MapStrAny, len(menuIdList))
		for i, menuId := range menuIdList {
			roleMenuList[i] = g.MapStrAny{
				"roleId": roleId,
				"menuId": menuId,
			}
		}
		v.DBM(model.NewBaseSysRoleMenu()).Data(roleMenuList).Insert()
	}
	// 更新部门权限
	v.DBM(model.NewBaseSysRoleDepartment()).Where("roleId", roleId).Delete()
	if len(departmentIds) > 0 {
		roleDepartmentList := make([]g.MapStrAny, len(departmentIds))
		for i, departmentId := range departmentIds {
			roleDepartmentList[i] = g.MapStrAny{
				"roleId":       roleId,
				"departmentId": departmentId,
			}
		}
		v.DBM(model.NewBaseSysRoleDepartment()).Data(roleDepartmentList).Insert()
	}
	// 刷新权限
	userRoles, err := v.DBM(model.NewBaseSysUserRole()).Where("roleId", roleId).All()
	if err != nil {
		return
	}
	baseSysPermsService := NewBaseSysPermsService()
	for _, v := range userRoles {
		vmap := v.Map()
		if vmap["userId"] != nil {
			if e := baseSysPermsService.RefreshPerms(ctx, gconv.Uint(vmap["userId"])); e != nil {
				err = e
				return
			}
		}
	}

	return

}

// GetByUser get array  roleId by userId
func (s *BaseSysRoleService) GetByUser(userId uint) (roles []string, err error) {
	res, err := v.DBM(baseSysUserRole).Where("userId", userId).Array("roleId")
	if err != nil {
		return nil, err
	}
	for _, v := range res {
		roles = append(roles, gconv.String(v))
	}
	return roles, nil
}

// BaseSysRoleService Info 方法重构
func (s *BaseSysRoleService) ServiceInfo(ctx context.Context, req *v.InfoReq) (data interface{}, err error) {
	info, err := v.DBM(s.Model).Where("id = ?", req.Id).One()
	if err != nil {
		return nil, err
	}
	if !info.IsEmpty() {
		var menus gdb.Result
		// 管理员角色（label=admin）展示全部菜单
		if satokenIsAdmin([]string{gconv.String(req.Id)}) {
			menus, err = v.DBM(model.NewBaseSysMenu()).Fields("id AS menuId").All()
			if err != nil {
				return nil, err
			}
		} else {
			menus, err = v.DBM(model.NewBaseSysRoleMenu()).Where("roleId", req.Id).All()
			if err != nil {
				return nil, err
			}
		}
		menuIdList := garray.NewIntArray()
		for _, v := range menus {
			menuIdList.Append(gconv.Int(v["menuId"]))
		}
		var departments gdb.Result
		if satokenIsAdmin([]string{gconv.String(req.Id)}) {
			departments, err = v.DBM(model.NewBaseSysRoleDepartment()).All()
			if err != nil {
				return nil, err
			}
		} else {
			departments, err = v.DBM(model.NewBaseSysRoleDepartment()).Where("roleId", req.Id).All()
			if err != nil {
				return nil, err
			}
		}
		departmentIdList := garray.NewIntArray()
		for _, v := range departments {
			departmentIdList.Append(gconv.Int(v["departmentId"]))
		}
		result := gconv.Map(info)
		result["menuIdList"] = menuIdList.Slice()
		result["departmentIdList"] = departmentIdList.Slice()
		data = result
		return
	}
	data = g.Map{}
	return
}

// NewBaseSysRoleService create a new BaseSysRoleService
func NewBaseSysRoleService() *BaseSysRoleService {
	return &BaseSysRoleService{
		Service: &v.Service{
			Model: model.NewBaseSysRole(),
			ListQueryOp: &v.QueryOp{
				Where: func(ctx context.Context) [][]interface{} {
					var (
						admin  = v.GetAdmin(ctx)
						userId = admin.UserId
					)
					return [][]interface{}{
						{"label != ?", g.Slice{"admin"}, true},
						{"(" + v.QuoteIdent(model.NewBaseSysRole().GroupName(), "userId") + "=? or id in (?))", g.Slice{userId, admin.RoleIds}, !admin.IsAdmin},
					}
				},
			},
			PageQueryOp: &v.QueryOp{
				KeyWordField: []string{"name", "label"},
				AddOrderby:   map[string]string{},
				Where: func(ctx context.Context) [][]interface{} {
					var (
						admin  = v.GetAdmin(ctx)
						userId = admin.UserId
					)
					return [][]interface{}{
						{"label != ?", g.Slice{"admin"}, true},
						{"(" + v.QuoteIdent(model.NewBaseSysRole().GroupName(), "userId") + "=? or id in (?))", g.Slice{userId, admin.RoleIds}, !admin.IsAdmin},
					}
				},
			},
			InsertParam: func(ctx context.Context) map[string]interface{} {
				return g.Map{"userId": v.GetAdmin(ctx).UserId}
			},
			UniqueKey: map[string]string{
				"name":  "角色名称不能重复",
				"label": "角色标识不能重复",
			},
		},
	}
}
