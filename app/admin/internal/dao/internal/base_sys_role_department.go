// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. Created at 2025-03-19 16:43:14
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// BaseSysRoleDepartmentDao is the data access object for the table base_sys_role_department.
type BaseSysRoleDepartmentDao struct {
	table    string                       // table is the underlying table name of the DAO.
	group    string                       // group is the database configuration group name of the current DAO.
	columns  BaseSysRoleDepartmentColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler           // handlers for customized model modification.
}

// BaseSysRoleDepartmentColumns defines and stores column names for the table base_sys_role_department.
type BaseSysRoleDepartmentColumns struct {
	Id           string // ID
	CreatedAt    string // 创建时间
	UpdatedAt    string // 更新时间
	TenantId     string // 租户ID
	RoleId       string // 角色ID
	DepartmentId string // 部门ID
	DeletedAt    string // 软删除时间
}

// baseSysRoleDepartmentColumns holds the columns for the table base_sys_role_department.
var baseSysRoleDepartmentColumns = BaseSysRoleDepartmentColumns{
	Id:           "id",
	CreatedAt:    "created_at",
	UpdatedAt:    "updated_at",
	TenantId:     "tenant_id",
	RoleId:       "role_id",
	DepartmentId: "department_id",
	DeletedAt:    "deleted_at",
}

// NewBaseSysRoleDepartmentDao creates and returns a new DAO object for table data access.
func NewBaseSysRoleDepartmentDao(handlers ...gdb.ModelHandler) *BaseSysRoleDepartmentDao {
	return &BaseSysRoleDepartmentDao{
		group:    "default",
		table:    "base_sys_role_department",
		columns:  baseSysRoleDepartmentColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *BaseSysRoleDepartmentDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *BaseSysRoleDepartmentDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *BaseSysRoleDepartmentDao) Columns() BaseSysRoleDepartmentColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *BaseSysRoleDepartmentDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *BaseSysRoleDepartmentDao) Ctx(ctx context.Context) *gdb.Model {
	model := dao.DB().Model(dao.table)
	for _, handler := range dao.handlers {
		model = handler(model)
	}
	return model.Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rolls back the transaction and returns the error if function f returns a non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note: Do not commit or roll back the transaction in function f,
// as it is automatically handled by this function.
func (dao *BaseSysRoleDepartmentDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
