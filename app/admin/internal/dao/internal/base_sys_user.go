// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. Created at 2025-03-14 15:21:22
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// BaseSysUserDao is the data access object for the table base_sys_user.
type BaseSysUserDao struct {
	table   string             // table is the underlying table name of the DAO.
	group   string             // group is the database configuration group name of the current DAO.
	columns BaseSysUserColumns // columns contains all the column names of Table for convenient usage.
}

// BaseSysUserColumns defines and stores column names for the table base_sys_user.
type BaseSysUserColumns struct {
	Id           string // ID
	CreateTime   string // 创建时间
	UpdateTime   string // 更新时间
	TenantId     string // 租户ID
	DepartmentId string // 部门ID
	Name         string // 姓名
	Username     string // 用户名
	Password     string // 密码
	PasswordV    string // 密码版本, 作用是改完密码，让原来的token失效
	NickName     string // 昵称
	HeadImg      string // 头像
	Phone        string // 手机
	Email        string // 邮箱
	Remark       string // 备注
	Status       string // 状态 0-禁用 1-启用
	SocketId     string // socketId
}

// baseSysUserColumns holds the columns for the table base_sys_user.
var baseSysUserColumns = BaseSysUserColumns{
	Id:           "id",
	CreateTime:   "createTime",
	UpdateTime:   "updateTime",
	TenantId:     "tenantId",
	DepartmentId: "departmentId",
	Name:         "name",
	Username:     "username",
	Password:     "password",
	PasswordV:    "passwordV",
	NickName:     "nickName",
	HeadImg:      "headImg",
	Phone:        "phone",
	Email:        "email",
	Remark:       "remark",
	Status:       "status",
	SocketId:     "socketId",
}

// NewBaseSysUserDao creates and returns a new DAO object for table data access.
func NewBaseSysUserDao() *BaseSysUserDao {
	return &BaseSysUserDao{
		group:   "default",
		table:   "base_sys_user",
		columns: baseSysUserColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *BaseSysUserDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *BaseSysUserDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *BaseSysUserDao) Columns() BaseSysUserColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *BaseSysUserDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *BaseSysUserDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rolls back the transaction and returns the error if function f returns a non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note: Do not commit or roll back the transaction in function f,
// as it is automatically handled by this function.
func (dao *BaseSysUserDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
