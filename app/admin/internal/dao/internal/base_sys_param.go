// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. Created at 2025-03-14 15:21:22
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// BaseSysParamDao is the data access object for the table base_sys_param.
type BaseSysParamDao struct {
	table   string              // table is the underlying table name of the DAO.
	group   string              // group is the database configuration group name of the current DAO.
	columns BaseSysParamColumns // columns contains all the column names of Table for convenient usage.
}

// BaseSysParamColumns defines and stores column names for the table base_sys_param.
type BaseSysParamColumns struct {
	Id         string // ID
	CreateTime string // 创建时间
	UpdateTime string // 更新时间
	TenantId   string // 租户ID
	KeyName    string // 键
	Name       string // 名称
	Data       string // 数据
	DataType   string // 数据类型 0-字符串 1-富文本 2-文件
	Remark     string // 备注
}

// baseSysParamColumns holds the columns for the table base_sys_param.
var baseSysParamColumns = BaseSysParamColumns{
	Id:         "id",
	CreateTime: "createTime",
	UpdateTime: "updateTime",
	TenantId:   "tenantId",
	KeyName:    "keyName",
	Name:       "name",
	Data:       "data",
	DataType:   "dataType",
	Remark:     "remark",
}

// NewBaseSysParamDao creates and returns a new DAO object for table data access.
func NewBaseSysParamDao() *BaseSysParamDao {
	return &BaseSysParamDao{
		group:   "default",
		table:   "base_sys_param",
		columns: baseSysParamColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *BaseSysParamDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *BaseSysParamDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *BaseSysParamDao) Columns() BaseSysParamColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *BaseSysParamDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *BaseSysParamDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rolls back the transaction and returns the error if function f returns a non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note: Do not commit or roll back the transaction in function f,
// as it is automatically handled by this function.
func (dao *BaseSysParamDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
