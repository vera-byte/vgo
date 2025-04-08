package base

import (
	"context"

	"github.com/vera-byte/vgo/app/admin/internal/dao"
	"github.com/vera-byte/vgo/app/admin/internal/model/entity"
	"github.com/vera-byte/vgo/app/admin/internal/service"

	"github.com/gogf/gf/v2/crypto/gmd5"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	vck_request "github.com/vera-byte/vgo/vgo_core_kit/request"
)

func init() {
	service.RegisterBaseSysUserLogic(NewBaseSysUserLogic())

}

type sBaseSysUserLogic struct {
}

func NewBaseSysUserLogic() *sBaseSysUserLogic {
	return &sBaseSysUserLogic{}
}

// 用户信息
func (s *sBaseSysUserLogic) Person(ctx context.Context, userId int64) (user *entity.BaseSysUser, err error) {
	var (
		daoBaseSysUser = dao.BaseSysUser
	)
	err = daoBaseSysUser.Ctx(ctx).Where(daoBaseSysUser.Columns().Id, userId).Scan(&user)
	if err != nil {
		return nil, err
	}
	return
}

// 更新用户信息
func (s *sBaseSysUserLogic) PersonUpdate(ctx context.Context, userId int64, k string, v string, o string) error {
	var (
		daoBaseSysUser = dao.BaseSysUser
		data           = g.Map{k: v}
		user           *entity.BaseSysUser
	)
	if k == "password" {
		if v == o {
			return gerror.New("新密码不能和旧密码相同")
		}
		v, _ = gmd5.Encrypt(v)
		oldPassword, _ := gmd5.Encrypt(o)
		daoBaseSysUser.Ctx(ctx).Where(daoBaseSysUser.Columns().Id, userId).Where(daoBaseSysUser.Columns().Password, oldPassword).Scan(&user)
		if user == nil {
			return gerror.New("原密码不正确")
		}
		if oldPassword == v {
			return gerror.New("不能使用近期密码")
		}
		// 重新赋值加密数据
		data = g.Map{k: v}
		// 校验原密码是否正确
		data[daoBaseSysUser.Columns().PasswordV] = user.PasswordV + 1

	}
	_, err := daoBaseSysUser.Ctx(ctx).Data(data).Where(daoBaseSysUser.Columns().Id, userId).Update()
	return err
}

// page
func (s *sBaseSysUserLogic) Page(ctx context.Context, pageReq *vck_request.PageReq, departmentIds []int64) (res []*entity.BaseSysUser, pagination *vck_request.Pagination, err error) {
	var (
		m = dao.BaseSysUser.Ctx(ctx).WhereIn(dao.BaseSysUser.Columns().DepartmentId, departmentIds).Order(pageReq.Order, pageReq.Sort).Limit(pageReq.Size).Offset((pageReq.Page - 1) * pageReq.Size)
	)
	pagination = &vck_request.Pagination{
		Page: pageReq.Page,
		Size: pageReq.Size,
	}
	pagination.Total, _ = m.Clone().Count()

	err = m.Scan(&res)
	return
}
