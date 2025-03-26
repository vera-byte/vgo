package admin_comm

import (
	"context"

	protobuf "vgo/app/admin/api/comm/v1"
	v1 "vgo/app/gateway/api/admin_comm/v1"
)

func (c *ControllerV1) Person(ctx context.Context, req *v1.PersonReq) (res *v1.PersonRes, err error) {
	personInfo, err := c.AdminBaseCommClient.Person(ctx, &protobuf.PersonRpcInvoke{})
	if err != nil {
		return nil, err
	}
	return &v1.PersonRes{
		HeadImg:      personInfo.HeadImg,
		Name:         personInfo.Name,
		DepartmentId: personInfo.DepartmentId,
		NickName:     personInfo.NickName,
		Id:           personInfo.Id,
		PasswordV:    personInfo.PasswordV,
		Phone:        personInfo.Phone,
		Email:        personInfo.Email,
		CreateTime:   personInfo.CreateTime,
		UpdateTime:   personInfo.UpdateTime,
		Username:     personInfo.UserName,
		Status:       personInfo.Status,
		TenantId:     personInfo.TenantId,
		Remark:       personInfo.Remark,
	}, nil
}
