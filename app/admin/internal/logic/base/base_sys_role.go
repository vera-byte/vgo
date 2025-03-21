package base

import (
	"context"
	"vgo/app/admin/internal/dao"
	"vgo/app/admin/internal/service"

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

	res, err := dao.BaseSysUserRole.Ctx(ctx).Where("user_id = ?", userId).Array("role_id")
	for _, v := range res {
		roles = append(roles, gconv.String(v))
	}
	return roles, err
}
