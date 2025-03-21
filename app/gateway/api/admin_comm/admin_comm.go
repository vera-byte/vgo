// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package admin_comm

import (
	"context"

	"vgo/app/gateway/api/admin_comm/v1"
)

type IAdminCommV1 interface {
	Person(ctx context.Context, req *v1.PersonReq) (res *v1.PersonRes, err error)
	Permmenu(ctx context.Context, req *v1.PermmenuReq) (res *v1.PermmenuRes, err error)
}
