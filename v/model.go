package v

import (
	"time"
)

type IModel interface {
	TableName() string
	GroupName() string
}
type Model struct {
	ID         uint      `json:"id"`
	CreateTime time.Time `json:"createTime"` // 创建时间
	UpdateTime time.Time `json:"updateTime"` // 更新时间
	DeletedAt  time.Time `json:"deletedAt"`
}

// 返回表名
func (m *Model) TableName() string {
	return "this_table_should_not_exist"
}

// 返回分组名
func (m *Model) GroupName() string {
	return "default"
}

func NewModel() *Model {
	return &Model{
		ID:         0,
		CreateTime: time.Now(), // 避免结构体插入空时间违反 NOT NULL 约束
		UpdateTime: time.Now(),
		DeletedAt:  time.Time{},
	}
}
