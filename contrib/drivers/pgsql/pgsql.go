// Package pgsql PostgreSQL数据库驱动包
// 功能: 提供PostgreSQL数据库的GORM连接驱动实现
package pgsql

import (
	"fmt"

	// GoFrame框架相关包
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/text/gstr"

	// PostgreSQL驱动相关包
	_ "github.com/gogf/gf/contrib/drivers/pgsql/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	// 项目内部包
	"github.com/vera-byte/vgo/v/vdb"
)

// DriverPgsql PostgreSQL数据库驱动结构体
// 功能: 实现vdb.Driver接口，提供PostgreSQL数据库连接功能
type DriverPgsql struct{}

// NewDriverPgsql 创建PostgreSQL驱动实例
// 功能: 创建并返回一个新的PostgreSQL驱动实例
// 返回值: *DriverPgsql - PostgreSQL驱动实例
func NewDriverPgsql() *DriverPgsql {
	return &DriverPgsql{}
}

// ExactNamingStrategy 精确命名策略结构体
// 功能: 自定义GORM命名策略，保持表名和字段名不变
type ExactNamingStrategy struct {
	schema.NamingStrategy
}

// TableName 获取表名
// 功能: 返回原始表名，不进行任何转换
// 参数: table string - 原始表名
// 返回值: string - 处理后的表名（保持不变）
func (ns ExactNamingStrategy) TableName(table string) string {
	return table
}

// ColumnName 获取字段名
// 功能: 返回原始字段名，不进行任何转换
// 参数: 
//   - table string - 表名
//   - column string - 原始字段名
// 返回值: string - 处理后的字段名（保持不变）
func (ns ExactNamingStrategy) ColumnName(table, column string) string {
	return column
}
// GetConn 获取数据库连接
// 功能: 根据配置节点创建PostgreSQL数据库连接
// 参数: config *gdb.ConfigNode - 数据库配置节点
// 返回值: 
//   - db *gorm.DB - GORM数据库连接实例
//   - err error - 错误信息
func (d *DriverPgsql) GetConn(config *gdb.ConfigNode) (db *gorm.DB, err error) {
	var source string

	// 处理连接字符串配置
	if config.Link != "" {
		// ============================================================================
		// 从v2.2.0开始已弃用Link配置方式
		// ============================================================================
		source = config.Link
		// 运行时自定义数据库名称
		if config.Name != "" {
			source, _ = gregex.ReplaceString(`dbname=([\w\.\-]+)+`, "dbname="+config.Name, source)
		}
	} else {
		// 构建标准PostgreSQL连接字符串
		if config.Name != "" {
			source = fmt.Sprintf(
				"user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
				config.User, config.Pass, config.Host, config.Port, config.Name,
			)
		} else {
			source = fmt.Sprintf(
				"user=%s password=%s host=%s port=%s sslmode=disable",
				config.User, config.Pass, config.Host, config.Port,
			)
		}

		// 添加时区配置
		if config.Timezone != "" {
			source = fmt.Sprintf("%s timezone=%s", source, config.Timezone)
		}

		// 处理额外配置参数
		if config.Extra != "" {
			var extraMap map[string]interface{}
			if extraMap, err = gstr.Parse(config.Extra); err != nil {
				return nil, fmt.Errorf("解析额外配置参数失败: %w", err)
			}
			for k, v := range extraMap {
				source += fmt.Sprintf(` %s=%s`, k, v)
			}
		}
	}

	// 创建GORM配置
	gormConfig := &gorm.Config{
		NamingStrategy: ExactNamingStrategy{},
	}

	// 创建数据库连接
	db, err = gorm.Open(postgres.Open(source), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("创建PostgreSQL连接失败: %w", err)
	}

	// 获取底层sql.DB实例进行连接池配置
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("获取底层数据库连接失败: %w", err)
	}

	// 设置连接池参数
	if config.MaxIdleConnCount > 0 {
		sqlDB.SetMaxIdleConns(config.MaxIdleConnCount)
	}
	if config.MaxOpenConnCount > 0 {
		sqlDB.SetMaxOpenConns(config.MaxOpenConnCount)
	}
	if config.MaxConnLifeTime > 0 {
		sqlDB.SetConnMaxLifetime(config.MaxConnLifeTime)
	}

	return db, nil
}

// init 初始化函数
// 功能: 注册PostgreSQL数据库驱动到vdb驱动管理器
func init() {
	var (
		err         error
		driverObj   = NewDriverPgsql()
		driverNames = g.SliceStr{"pgsql", "postgres", "postgresql"}
	)
	
	// 注册多个驱动名称别名
	for _, driverName := range driverNames {
		if err = vdb.Register(driverName, driverObj); err != nil {
			panic(fmt.Sprintf("注册PostgreSQL驱动失败 [%s]: %v", driverName, err))
		}
	}
}
