// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. Created at 2025-03-14 15:21:22
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// BaseSysUserRoleDao is the data access object for the table base_sys_user_role.
type BaseSysUserRoleDao struct {
	table   string                 // table is the underlying table name of the DAO.
	group   string                 // group is the database configuration group name of the current DAO.
	columns BaseSysUserRoleColumns // columns contains all the column names of Table for convenient usage.
}

// BaseSysUserRoleColumns defines and stores column names for the table base_sys_user_role.
type BaseSysUserRoleColumns struct {
	Id         string // ID
	CreateTime string // 创建时间
	UpdateTime string // 更新时间
	TenantId   string // 租户ID
	UserId     string // 用户ID
	RoleId     string // 角色ID
}

// baseSysUserRoleColumns holds the columns for the table base_sys_user_role.
var baseSysUserRoleColumns = BaseSysUserRoleColumns{
	Id:         "id",
	CreateTime: "createTime",
	UpdateTime: "updateTime",
	TenantId:   "tenantId",
	UserId:     "userId",
	RoleId:     "roleId",
}

// NewBaseSysUserRoleDao creates and returns a new DAO object for table data access.
func NewBaseSysUserRoleDao() *BaseSysUserRoleDao {
	return &BaseSysUserRoleDao{
		group:   "default",
		table:   "base_sys_user_role",
		columns: baseSysUserRoleColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *BaseSysUserRoleDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *BaseSysUserRoleDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *BaseSysUserRoleDao) Columns() BaseSysUserRoleColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *BaseSysUserRoleDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *BaseSysUserRoleDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rolls back the transaction and returns the error if function f returns a non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note: Do not commit or roll back the transaction in function f,
// as it is automatically handled by this function.
func (dao *BaseSysUserRoleDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
