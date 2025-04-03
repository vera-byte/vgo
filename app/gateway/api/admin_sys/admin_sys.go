// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package admin_sys

import (
	"context"

	"vgo/app/gateway/api/admin_sys/v1"
)

type IAdminSysV1 interface {
	LogPage(ctx context.Context, req *v1.LogPageReq) (res *v1.LogPageRes, err error)
	DepartmentList(ctx context.Context, req *v1.DepartmentListReq) (res *v1.DepartmentListRes, err error)
	UserPage(ctx context.Context, req *v1.UserPageReq) (res *v1.UserPageRes, err error)
}
