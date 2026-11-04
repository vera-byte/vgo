package model

import "github.com/vera-byte/vgo/v"

const TableNameDemoSample = "demo_sample"

// DemoSample 示例表 demo_sample
type DemoSample struct {
	*v.Model
	Name string `json:"name"` // 名称
}

// TableName DemoSample 的表名
func (*DemoSample) TableName() string {
	return TableNameDemoSample
}

// NewDemoSample 创建一个 DemoSample 模型实例
func NewDemoSample() *DemoSample {
	return &DemoSample{Model: v.NewModel()}
}
