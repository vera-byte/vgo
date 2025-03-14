// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. Created at 2025-03-14 15:21:22
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// BaseSysDepartmentDao is the data access object for the table base_sys_department.
type BaseSysDepartmentDao struct {
	table   string                   // table is the underlying table name of the DAO.
	group   string                   // group is the database configuration group name of the current DAO.
	columns BaseSysDepartmentColumns // columns contains all the column names of Table for convenient usage.
}

// BaseSysDepartmentColumns defines and stores column names for the table base_sys_department.
type BaseSysDepartmentColumns struct {
	Id         string // ID
	CreateTime string // 创建时间
	UpdateTime string // 更新时间
	TenantId   string // 租户ID
	Name       string // 部门名称
	ParentId   string // 上级部门ID
	OrderNum   string // 排序
}

// baseSysDepartmentColumns holds the columns for the table base_sys_department.
var baseSysDepartmentColumns = BaseSysDepartmentColumns{
	Id:         "id",
	CreateTime: "createTime",
	UpdateTime: "updateTime",
	TenantId:   "tenantId",
	Name:       "name",
	ParentId:   "parentId",
	OrderNum:   "orderNum",
}

// NewBaseSysDepartmentDao creates and returns a new DAO object for table data access.
func NewBaseSysDepartmentDao() *BaseSysDepartmentDao {
	return &BaseSysDepartmentDao{
		group:   "default",
		table:   "base_sys_department",
		columns: baseSysDepartmentColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *BaseSysDepartmentDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *BaseSysDepartmentDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *BaseSysDepartmentDao) Columns() BaseSysDepartmentColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *BaseSysDepartmentDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *BaseSysDepartmentDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rolls back the transaction and returns the error if function f returns a non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note: Do not commit or roll back the transaction in function f,
// as it is automatically handled by this function.
func (dao *BaseSysDepartmentDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
