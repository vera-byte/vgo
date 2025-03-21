// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. Created at 2025-03-19 16:43:14
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// BaseSysLogDao is the data access object for the table base_sys_log.
type BaseSysLogDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  BaseSysLogColumns  // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// BaseSysLogColumns defines and stores column names for the table base_sys_log.
type BaseSysLogColumns struct {
	Id        string // ID
	CreatedAt string // 创建时间
	UpdatedAt string // 更新时间
	TenantId  string // 租户ID
	UserId    string // 用户ID
	Action    string // 行为
	Ip        string // ip
	Params    string // 参数
	DeletedAt string // 软删除时间
}

// baseSysLogColumns holds the columns for the table base_sys_log.
var baseSysLogColumns = BaseSysLogColumns{
	Id:        "id",
	CreatedAt: "created_at",
	UpdatedAt: "updated_at",
	TenantId:  "tenant_id",
	UserId:    "user_id",
	Action:    "action",
	Ip:        "ip",
	Params:    "params",
	DeletedAt: "deleted_at",
}

// NewBaseSysLogDao creates and returns a new DAO object for table data access.
func NewBaseSysLogDao(handlers ...gdb.ModelHandler) *BaseSysLogDao {
	return &BaseSysLogDao{
		group:    "default",
		table:    "base_sys_log",
		columns:  baseSysLogColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *BaseSysLogDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *BaseSysLogDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *BaseSysLogDao) Columns() BaseSysLogColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *BaseSysLogDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *BaseSysLogDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *BaseSysLogDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
