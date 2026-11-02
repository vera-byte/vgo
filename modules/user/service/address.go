package service

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	userModel "github.com/vera-byte/vgo/modules/user/model"
	"github.com/vera-byte/vgo/v"
)

// UserAddressService 收货地址服务（参照 cool-admin user 模块 UserAddressService）
type UserAddressService struct {
	*v.Service
}

// NewUserAddressService 创建地址服务（owner=当前登录用户，行级过滤+注入）
func NewUserAddressService() *UserAddressService {
	return &UserAddressService{
		Service: &v.Service{
			Model: userModel.NewUserAddress(),
			InsertParam: func(ctx context.Context) g.MapStrAny {
				return g.MapStrAny{"userId": CurrentUser(ctx).ID}
			},
			PageQueryOp: &v.QueryOp{
				Where: func(ctx context.Context) []g.Array {
					return []g.Array{{`("userId" = ?)`, []interface{}{CurrentUser(ctx).ID}, true}}
				},
				AddOrderby: g.MapStrStr{"isDefault": "DESC"},
			},
			ListQueryOp: &v.QueryOp{
				Where: func(ctx context.Context) []g.Array {
					return []g.Array{{`("userId" = ?)`, []interface{}{CurrentUser(ctx).ID}, true}}
				},
				AddOrderby: g.MapStrStr{"isDefault": "DESC"},
			},
		},
	}
}

// ownAddress 校验地址归属当前用户（防越权），返回记录
func (s *UserAddressService) ownAddress(ctx context.Context, id int) (map[string]any, error) {
	record, err := v.DBM(userModel.NewUserAddress()).
		Where("id", id).Where("userId", CurrentUser(ctx).ID).One()
	if err != nil {
		return nil, err
	}
	if record.IsEmpty() {
		return nil, gerror.New("地址不存在")
	}
	return record.Map(), nil
}

// ServiceInfo 查询地址（限定本人，防越权）
func (s *UserAddressService) ServiceInfo(ctx context.Context, req *v.InfoReq) (data interface{}, err error) {
	record, err := s.ownAddress(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return record, nil
}

// ServiceUpdate 更新地址（限定本人，防越权）
func (s *UserAddressService) ServiceUpdate(ctx context.Context, req *v.UpdateReq) (data interface{}, err error) {
	id := g.RequestFromCtx(ctx).Get("id").Int()
	if _, err = s.ownAddress(ctx, id); err != nil {
		return nil, err
	}
	return s.Service.ServiceUpdate(ctx, req)
}

// ServiceDelete 删除地址（限定本人，防越权）
func (s *UserAddressService) ServiceDelete(ctx context.Context, req *v.DeleteReq) (data interface{}, err error) {
	for _, id := range gconv.Ints(g.RequestFromCtx(ctx).Get("ids").Slice()) {
		if _, err = s.ownAddress(ctx, id); err != nil {
			return nil, err
		}
	}
	return s.Service.ServiceDelete(ctx, req)
}

// ModifyAfter 增改后处理：设为默认地址时取消该用户其他默认
func (s *UserAddressService) ModifyAfter(ctx context.Context, method string, param g.MapStrAny) error {
	if method != "Add" && method != "Update" {
		return nil
	}
	userId := CurrentUser(ctx).ID
	id := gconv.Uint(param["id"])
	isDefault := gconv.Bool(param["isDefault"])
	if !isDefault {
		return nil
	}
	m := v.DBM(userModel.NewUserAddress()).Where("userId", userId).Where("isDefault", true)
	if id > 0 {
		m = m.Where("id <> ?", id)
	}
	_, err := m.Data(g.Map{"isDefault": false}).Update()
	return err
}

// Default 默认地址
func (s *UserAddressService) Default(ctx context.Context) (map[string]any, error) {
	record, err := v.DBM(userModel.NewUserAddress()).
		Where("userId", CurrentUser(ctx).ID).Where("isDefault", true).One()
	if err != nil {
		return nil, err
	}
	if record.IsEmpty() {
		return nil, gerror.New("暂无默认地址")
	}
	return record.Map(), nil
}
