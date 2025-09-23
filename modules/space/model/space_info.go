package model

import (
	"github.com/vera-byte/vgo/v"
)

const TableNameSpaceInfo = "space_info"

// SpaceInfo mapped from table <space_info>
type SpaceInfo struct {
	*v.Model
	URL        string `json:"url"`        // 地址
	Type       string `json:"type"`       // 类型
	ClassifyID *int64 `json:"classifyId"` // 分类ID
}

// TableName SpaceInfo's table name
func (*SpaceInfo) TableName() string {
	return TableNameSpaceInfo
}

// GroupName SpaceInfo's table group
func (*SpaceInfo) GroupName() string {
	return "default"
}

// NewSpaceInfo create a new SpaceInfo
func NewSpaceInfo() *SpaceInfo {
	return &SpaceInfo{
		Model: v.NewModel(),
	}
}
