// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. Created at 2025-03-14 15:21:22
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// BaseSysRoleDao is the data access object for the table base_sys_role.
type BaseSysRoleDao struct {
	table   string             // table is the underlying table name of the DAO.
	group   string             // group is the database configuration group name of the current DAO.
	columns BaseSysRoleColumns // columns contains all the column names of Table for convenient usage.
}

// BaseSysRoleColumns defines and stores column names for the table base_sys_role.
type BaseSysRoleColumns struct {
	Id               string // ID
	CreateTime       string // 创建时间
	UpdateTime       string // 更新时间
	TenantId         string // 租户ID
	UserId           string // 用户ID
	Name             string // 名称
	Label            string // 角色标签
	Remark           string // 备注
	Relevance        string // 数据权限是否关联上下级
	MenuIdList       string // 菜单权限
	DepartmentIdList string // 部门权限
}

// baseSysRoleColumns holds the columns for the table base_sys_role.
var baseSysRoleColumns = BaseSysRoleColumns{
	Id:               "id",
	CreateTime:       "createTime",
	UpdateTime:       "updateTime",
	TenantId:         "tenantId",
	UserId:           "userId",
	Name:             "name",
	Label:            "label",
	Remark:           "remark",
	Relevance:        "relevance",
	MenuIdList:       "menuIdList",
	DepartmentIdList: "departmentIdList",
}

// NewBaseSysRoleDao creates and returns a new DAO object for table data access.
func NewBaseSysRoleDao() *BaseSysRoleDao {
	return &BaseSysRoleDao{
		group:   "default",
		table:   "base_sys_role",
		columns: baseSysRoleColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *BaseSysRoleDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *BaseSysRoleDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *BaseSysRoleDao) Columns() BaseSysRoleColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *BaseSysRoleDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *BaseSysRoleDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rolls back the transaction and returns the error if function f returns a non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note: Do not commit or roll back the transaction in function f,
// as it is automatically handled by this function.
func (dao *BaseSysRoleDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
