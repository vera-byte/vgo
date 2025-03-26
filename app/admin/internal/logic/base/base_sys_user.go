package base

import (
	"context"
	"vgo/app/admin/internal/dao"
	"vgo/app/admin/internal/model/entity"
	"vgo/app/admin/internal/service"
)

func init() {
	service.RegisterBaseSysUserLogic(NewBaseSysUserLogic())

}

type sBaseSysUserLogic struct {
}

func NewBaseSysUserLogic() *sBaseSysUserLogic {
	return &sBaseSysUserLogic{}
}

func (s *sBaseSysUserLogic) Person(ctx context.Context, userId int64) (user *entity.BaseSysUser, err error) {
	err = dao.BaseSysUser.Ctx(ctx).Where("id = ?", userId).Scan(&user)
	if err != nil {
		return nil, err
	}
	return
}
