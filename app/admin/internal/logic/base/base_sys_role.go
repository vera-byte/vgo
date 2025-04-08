package base

import (
	"context"

	"github.com/vera-byte/vgo/app/admin/internal/dao"
	"github.com/vera-byte/vgo/app/admin/internal/service"

	"github.com/gogf/gf/v2/util/gconv"
)

func init() {
	service.RegisterBaseSysRoleLogic(NewBaseSysRoleLogic())
}
func NewBaseSysRoleLogic() *sBaseSysRoleLogic {
	return &sBaseSysRoleLogic{}
}

type sBaseSysRoleLogic struct {
}

// 通过用户ID获取角色集合
func (l *sBaseSysRoleLogic) GetByUser(ctx context.Context, userId int64) (roles []string, err error) {
	var (
		daoBaseSysUserRole = dao.BaseSysUserRole
	)
	res, err := daoBaseSysUserRole.Ctx(ctx).Where(daoBaseSysUserRole.Columns().UserId, userId).Array(daoBaseSysUserRole.Columns().RoleId)
	for _, v := range res {
		roles = append(roles, gconv.String(v))
	}
	return roles, err
}
