package model

import "github.com/vera-byte/vgo/v"

const TableNameRecycleData = "recycle_data"

// RecycleData 数据回收站：被删除数据的快照（参照 cool-admin-midway recycle 模块）
type RecycleData struct {
	*v.Model
	EntityInfo string `json:"entityInfo"` // 来源信息：{"group":"default","table":"base_sys_user"}
	UserID     *int64 `json:"userId"`     // 操作人
	Data       string `json:"data"`       // 被删除的数据（JSON 数组快照）
	Url        string `json:"url"`        // 请求的接口
	Params     string `json:"params"`     // 请求参数（JSON）
	Count      int    `json:"count"`      // 删除数据条数
}

// TableName RecycleData 的表名
func (*RecycleData) TableName() string {
	return TableNameRecycleData
}

// NewRecycleData 创建一个 RecycleData 模型实例
func NewRecycleData() *RecycleData {
	return &RecycleData{
		Model: v.NewModel(),
	}
}
