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
}
