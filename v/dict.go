package v

import (
	"context"

	"github.com/gogf/gf/v2/container/gvar"
)

// DictOption 字典下拉选项
type DictOption struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

// DictTypeInfo 字典类型信息
type DictTypeInfo struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Key  string `json:"key"`
}

// DictEntry 字典条目完整数据
type DictEntry struct {
	ID       uint      `json:"id"`
	Name     string    `json:"name"`     // 显示名
	Value    *gvar.Var `json:"value"`    // 条目值（gvar 便于按需转数值等）
	OrderNum int       `json:"orderNum"` // 排序号
	Remark   string    `json:"remark"`   // 备注
	ParentID *int      `json:"parentId"` // 父ID（树形字典）
}

// DictFullInfo 字典完整数据（类型 + 全部条目）
type DictFullInfo struct {
	ID      uint        `json:"id"`
	Name    string      `json:"name"`
	Key     string      `json:"key"`
	Entries []DictEntry `json:"entries"`
}

// DictProvider 字典数据源，由 dict 模块注册（见 RegisterDict）
type DictProvider interface {
	// GetLabel 按字典类型 key 与条目值取显示名；不存在返回 ok=false
	GetLabel(ctx context.Context, typeKey string, value string) (label string, ok bool)
	// GetValue 按字典类型 key 与显示名取条目值；不存在返回 ok=false
	GetValue(ctx context.Context, typeKey string, label string) (value *gvar.Var, ok bool)
	// GetOptions 取字典类型全部条目（按排序号升序）
	GetOptions(ctx context.Context, typeKey string) []DictOption
	// GetTypes 取全部字典类型（按 id 升序）
	GetTypes(ctx context.Context) []DictTypeInfo
	// GetInfo 按字典类型 key 与显示名取单条条目完整数据；不存在返回 nil
	GetInfo(ctx context.Context, typeKey string, label string) *DictEntry
}

var dictProvider DictProvider

// RegisterDict 注册字典数据源（dict 模块初始化时调用）
func RegisterDict(p DictProvider) {
	dictProvider = p
}

type dictFacade struct{}

// Dict 字典查询入口
var Dict dictFacade

// GetInfo 按显示名取单条条目完整数据（含 id/排序/备注/父级等字段），
// 形如 v.Dict.GetInfo(ctx, "admin_user_status", "正常")，不存在返回 nil
func (dictFacade) GetInfo(ctx context.Context, typeKey string, label string) *DictEntry {
	if dictProvider == nil {
		return nil
	}
	return dictProvider.GetInfo(ctx, typeKey, label)
}

// GetLabel 取显示名，不存在返回空字符串
func (dictFacade) GetLabel(ctx context.Context, typeKey string, value string) string {
	if dictProvider == nil {
		return ""
	}
	label, _ := dictProvider.GetLabel(ctx, typeKey, value)
	return label
}

// GetOptions 取字典类型全部条目列表（按排序号升序），便于下拉框渲染
func (dictFacade) GetOptions(ctx context.Context, typeKey string) []DictOption {
	if dictProvider == nil {
		return []DictOption{}
	}
	return dictProvider.GetOptions(ctx, typeKey)
}

// GetTypes 取全部字典类型列表（按 id 升序），便于管理页面渲染
func (dictFacade) GetTypes(ctx context.Context) []DictTypeInfo {
	if dictProvider == nil {
		return []DictTypeInfo{}
	}
	return dictProvider.GetTypes(ctx)
}

// GetValue 取条目值（gvar 便于按需转数值等），不存在返回 IsNil 的 gvar
func (dictFacade) GetValue(ctx context.Context, typeKey string, label string) *gvar.Var {
	if dictProvider == nil {
		return gvar.New(nil)
	}
	value, ok := dictProvider.GetValue(ctx, typeKey, label)
	if !ok {
		return gvar.New(nil)
	}
	return value
}
