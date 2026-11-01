package service

import (
	"context"

	"github.com/gogf/gf/v2/container/garray"
	"github.com/gogf/gf/v2/container/gset"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/vera-byte/vgo/modules/base/model"
	"github.com/vera-byte/vgo/v"
)

type BaseSysUserService struct {
	*v.Service
}

// Person 方法 返回不带密码的用户信息
func (s *BaseSysUserService) Person(userId uint) (res gdb.Record, err error) {
	m := v.DBM(s.Model)
	res, err = m.Where("id = ?", userId).FieldsEx("password").One()
	return
}

func (s *BaseSysUserService) ModifyBefore(ctx context.Context, method string, param g.MapStrAny) (err error) {
	if method == "Delete" {
		// 禁止删除超级管理员
		userIds := garray.NewIntArrayFrom(gconv.Ints(param["ids"]))
		currentId, found := userIds.Get(0)
		superAdminId := 1

		if userIds.Len() == 1 && found && currentId == superAdminId {
			err = gerror.New("超级管理员不能删除")
			return
		}

		// 删除超级管理员
		userIds.RemoveValue(1)
		g.RequestFromCtx(ctx).SetParam("ids", userIds.Slice())
	}
	return
}

func (s *BaseSysUserService) ModifyAfter(ctx context.Context, method string, param g.MapStrAny) (err error) {
	if method == "Delete" {
		userIds := garray.NewIntArrayFrom(gconv.Ints(param["ids"]))
		userIds.RemoveValue(1)
		// 删除用户时删除相关数据
		v.DBM(model.NewBaseSysUserRole()).WhereIn("userId", userIds.Slice()).Delete()
	}
	return
}

// ServiceAdd 方法 添加用户
func (s *BaseSysUserService) ServiceAdd(ctx context.Context, req *v.AddReq) (data interface{}, err error) {
	var (
		m      = v.DBM(s.Model)
		r      = g.RequestFromCtx(ctx)
		reqmap = r.GetMap()
	)
	// 如果reqmap["password"]不为空，则对密码进行bcrypt加密
	if !r.Get("password").IsNil() {
		hashed, e := hashPassword(r.Get("password").String())
		if e != nil {
			return nil, e
		}
		reqmap["password"] = hashed
	}
	lastInsertId, err := m.Data(reqmap).InsertAndGetId()
	if err != nil {
		return
	}
	data = g.Map{"id": lastInsertId}
	return
}

// ServiceInfo 方法 返回服务信息
func (s *BaseSysUserService) ServiceInfo(ctx g.Ctx, req *v.InfoReq) (data interface{}, err error) {
	result, err := s.Service.ServiceInfo(ctx, req)
	if err != nil {
		return nil, err
	}
	if result.(gdb.Record).IsEmpty() {
		return nil, nil
	}
	// g.DumpWithType(result)
	resultMap := result.(gdb.Record).Map()

	// 获取角色
	roleIds, err := v.DBM(model.NewBaseSysUserRole()).Where("userId", resultMap["id"]).Fields("roleId").Array()
	if err != nil {
		return nil, err
	}
	resultMap["roleIdList"] = roleIds
	data = resultMap

	return
}

// ServiceUpdate 方法 更新用户信息
func (s *BaseSysUserService) ServiceUpdate(ctx context.Context, req *v.UpdateReq) (data interface{}, err error) {
	var (
		admin = v.GetAdmin(ctx)
		m     = v.DBM(s.Model)
	)

	r := g.RequestFromCtx(ctx)
	rMap := r.GetMap()

	// 如果不传入ID代表更新当前用户
	userId := r.Get("id", admin.UserId).Uint()
	userInfo, err := m.Where("id = ?", userId).One()

	if err != nil {
		return
	}
	if userInfo.IsEmpty() {
		err = gerror.New("用户不存在")
		return
	}

	// 禁止禁用管理员账号（按角色 label=admin 判定）
	if !r.Get("status").IsNil() && r.Get("status").Int() == 0 {
		if roleIds, e := NewBaseSysRoleService().GetByUser(userId); e == nil && satokenIsAdmin(roleIds) {
			err = gerror.New("禁止禁用管理员账号")
			return
		}
	}

	// 如果请求的password不为空并且密码加密后的值有变动，说明要修改密码
	var rPassword = r.Get("password", "").String()
	passwordChanged := false
	if rPassword != "" && rPassword != userInfo["password"].String() {
		hashed, e := hashPassword(rPassword)
		if e != nil {
			return nil, e
		}
		rMap["password"] = hashed
		rMap["passwordV"] = userInfo["passwordV"].Int() + 1
		passwordChanged = true
	} else {
		delete(rMap, "password")
	}

	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) (err error) {
		roleModel := v.DBM(model.NewBaseSysUserRole()).TX(tx).Where("userId", userId)
		roleIds, err := roleModel.Fields("roleId").Array()
		if err != nil {
			return
		}

		// 如果请求参数中不包含roleIdList说明不修改角色信息
		if !r.Get("roleIdList").IsNil() {
			inRoleIdSet := gset.NewFrom(r.Get("roleIdList").Ints())
			roleIdsSet := gset.NewFrom(gconv.Ints(roleIds))

			// 如果请求的角色信息未发生变化则跳过更新逻辑
			if roleIdsSet.Diff(inRoleIdSet).Size() != 0 || inRoleIdSet.Diff(roleIdsSet).Size() != 0 {
				roleArray := garray.NewArray()
				inRoleIdSet.Iterator(func(v interface{}) bool {
					roleArray.PushRight(g.Map{
						"userId": gconv.Uint(userId),
						"roleId": gconv.Uint(v),
					})
					return true
				})

				_, err = roleModel.Delete()

				if err != nil {
					return err
				}
				_, err = roleModel.Fields("userId,roleId").Insert(roleArray)
				if err != nil {
					return err
				}
			}
		}

		_, err = m.TX(tx).Update(rMap)

		if err != nil {
			return err
		}
		return
	})
	if err != nil {
		return
	}
	// 密码变更后强制下线该用户全部会话
	if passwordChanged {
		logoutSatoken(userId)
	}
	// 角色变更后刷新该用户的权限
	if !r.Get("roleIdList").IsNil() {
		if e := NewBaseSysPermsService().RefreshPerms(ctx, userId); e != nil {
			g.Log().Warningf(ctx, "刷新用户权限失败: %v", e)
		}
	}
	return
}

// UpdateProfile 更新本人基础资料与密码（仅白名单字段，强制本人，防越权）
func (s *BaseSysUserService) UpdateProfile(ctx context.Context, userId uint) (err error) {
	r := g.RequestFromCtx(ctx)
	row, err := v.DBM(s.Model).Where("id", userId).One()
	if err != nil {
		return
	}
	if row.IsEmpty() {
		return gerror.New("用户不存在")
	}
	rMap := g.MapStrAny{}
	for _, field := range []string{"nickName", "headImg", "phone", "email", "remark"} {
		if val := r.Get(field); !val.IsNil() {
			rMap[field] = val.Val()
		}
	}
	passwordChanged := false
	if password := r.Get("password").String(); password != "" {
		hashed, e := hashPassword(password)
		if e != nil {
			return e
		}
		rMap["password"] = hashed
		rMap["passwordV"] = row["passwordV"].Int() + 1
		passwordChanged = true
	}
	if len(rMap) == 0 {
		return nil
	}
	if _, err = v.DBM(s.Model).Where("id", userId).Data(rMap).Update(); err != nil {
		return
	}
	// 密码变更后强制重新登录
	if passwordChanged {
		logoutSatoken(userId)
	}
	return
}

// Move 移动用户部门
func (s *BaseSysUserService) Move(ctx g.Ctx) (err error) {
	request := g.RequestFromCtx(ctx)
	departmentId := request.Get("departmentId").Int()
	userIds := request.Get("userIds").Slice()

	_, err = v.DBM(s.Model).WhereIn("id", userIds).Data(g.Map{"departmentId": departmentId}).Update()

	return
}

// NewBaseSysUserService 创建一个新的BaseSysUserService实例
func NewBaseSysUserService() *BaseSysUserService {
	return &BaseSysUserService{
		Service: &v.Service{
			Model:              model.NewBaseSysUser(),
			InfoIgnoreProperty: "password",
			UniqueKey: map[string]string{
				"username": "用户名不能重复",
			},
			PageQueryOp: &v.QueryOp{
				Select: `base_sys_user.*,"dept"."name" as "departmentName",string_agg("role"."name", ',') AS "roleName"`,
				Join: []*v.JoinOp{
					{
						Model:     model.NewBaseSysDepartment(),
						Alias:     "dept",
						Type:      "LeftJoin",
						Condition: `"base_sys_user"."departmentId" = "dept"."id"`,
					},
					{
						Model:     model.NewBaseSysUserRole(),
						Alias:     "user_role",
						Type:      "LeftJoin",
						Condition: `"base_sys_user"."id" = "user_role"."userId"`,
					},
					{
						Model:     model.NewBaseSysRole(),
						Alias:     `"role"`,
						Type:      "LeftJoin",
						Condition: `"role"."id" = "user_role"."roleId"`,
					},
				},
				Where: func(ctx context.Context) []g.Array {
					r := g.RequestFromCtx(ctx).GetMap()
					ids := gconv.SliceStr(r["departmentIds"])
					return []g.Array{
						{`("departmentId" IN (?))`, ids, len(ids) > 0},
					}
				},
				Extend: func(ctx g.Ctx, m *gdb.Model) *gdb.Model {
					return m.Group(`"base_sys_user"."id", "dept"."name"`)
				},
				KeyWordField: []string{"name", "username", "nickName"},
			},
		},
	}
}
