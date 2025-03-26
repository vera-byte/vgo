package comm

import (
	"context"
	v1 "vgo/app/admin/api/comm/v1"
	"vgo/app/admin/internal/service"

	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
	vck "github.com/vera-byte/vgo/vgo_core_kit"
)

type Controller struct {
	v1.UnimplementedBaseCommServer
}

func Register(s *grpcx.GrpcServer) {
	v1.RegisterBaseCommServer(s.Server, &Controller{})
}

func (*Controller) Permmenu(ctx context.Context, req *v1.PermmenuRpcInvoke) (res *v1.PermmenuRpcRes, err error) {

	var (
		baseSysMenuService = service.BaseSysMenuLogic()
		_res               = &v1.PermmenuRpcRes{}
		admin              = vck.GetAdminAtGrpcService(ctx)
	)
	_perms := baseSysMenuService.GetPerms(ctx, admin.RoleIds)
	_res.Perms = _perms
	menus := baseSysMenuService.GetMenus(ctx, admin.RoleIds, admin.Username == "admin")
	for _, item := range menus {
		_res.Menus = append(_res.Menus, &v1.Menu{
			ViewPath:   item.ViewPath,
			Name:       item.Name,
			Perms:      item.Perms,
			Icon:       item.Icon,
			OrderNum:   int64(item.OrderNum),
			IsShow:     item.IsShow,
			KeepAlive:  item.KeepAlive,
			Router:     item.Router,
			Type:       int32(item.Type),
			Id:         int64(item.Id),
			ParentId:   int64(item.ParentId),
			TenantId:   int64(item.TenantId),
			CreateTime: item.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdateTime: item.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	return _res, nil
}

func (*Controller) Person(ctx context.Context, req *v1.PersonRpcInvoke) (res *v1.PersonRpcRes, err error) {
	var (
		baseSysUserService = service.BaseSysUserLogic()
		admin              = vck.GetAdminAtGrpcService(ctx)
	)
	person, err := baseSysUserService.Person(ctx, admin.UserId)
	if err != nil {
		return nil, err
	}
	return &v1.PersonRpcRes{
		UserName:     person.Username,
		NickName:     person.NickName,
		HeadImg:      person.HeadImg,
		Phone:        person.Phone,
		Remark:       person.Remark,
		Name:         person.Name,
		PasswordV:    int64(person.PasswordV),
		Status:       int32(person.Status),
		TenantId:     int64(person.TenantId),
		Id:           int64(person.Id),
		DepartmentId: int64(person.DepartmentId),
		Email:        person.Email,
		CreateTime:   person.CreatedAt.Unix(),
		UpdateTime:   person.UpdatedAt.Unix(),
	}, nil
}

func (*Controller) LoginOut(ctx context.Context, req *v1.LoginOutRpcInvoke) (res *v1.LoginOutRes, err error) {
	return nil, service.BaseSysLoginLogic().LoginOut(ctx)
}
