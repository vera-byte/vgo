package service

import (
	"context"
	"encoding/json"
	"time"

	"github.com/gogf/gf/v2/container/gvar"
	"github.com/vera-byte/vgo/modules/dict/model"
	"github.com/vera-byte/vgo/v"
)

// dictCacheKey 字典全量缓存键（单键重建，类型增删不会残留脏缓存）。
// 缓存内容为 []v.DictFullInfo（类型 + 全部条目完整字段），查询一律走缓存。
const dictCacheKey = "dict:data"

// dictCacheTTL 字典缓存过期时间：写入会主动失效，过期兜底防止数据量大时长期占用内存
const dictCacheTTL = 10 * time.Minute

type dictLookup struct{}

func (dictLookup) GetLabel(ctx context.Context, typeKey string, value string) (string, bool) {
	if t := findDictType(dictData(ctx), typeKey); t != nil {
		for _, e := range t.Entries {
			if e.Value.String() == value && value != "" {
				return e.Name, true
			}
		}
	}
	return "", false
}

func (dictLookup) GetValue(ctx context.Context, typeKey string, label string) (*gvar.Var, bool) {
	e := (dictLookup{}).GetInfo(ctx, typeKey, label)
	if e != nil {
		return e.Value, true
	}
	return nil, false
}

func (dictLookup) GetOptions(ctx context.Context, typeKey string) []v.DictOption {
	options := []v.DictOption{}
	if t := findDictType(dictData(ctx), typeKey); t != nil {
		for _, e := range t.Entries {
			options = append(options, v.DictOption{Label: e.Name, Value: e.Value.String()})
		}
	}
	return options
}

func (dictLookup) GetTypes(ctx context.Context) []v.DictTypeInfo {
	list := dictData(ctx)
	types := make([]v.DictTypeInfo, 0, len(list))
	for _, t := range list {
		types = append(types, v.DictTypeInfo{ID: t.ID, Name: t.Name, Key: t.Key})
	}
	return types
}

// GetInfo 读缓存按显示名返回单条条目完整数据
func (dictLookup) GetInfo(ctx context.Context, typeKey string, label string) *v.DictEntry {
	if t := findDictType(dictData(ctx), typeKey); t != nil && label != "" {
		for i := range t.Entries {
			if t.Entries[i].Name == label {
				return &t.Entries[i]
			}
		}
	}
	return nil
}

// dictData 读取全量字典：缓存优先，未命中查库加载后回填（按类型 id 升序）
func dictData(ctx context.Context) []v.DictFullInfo {
	if cached, err := v.CacheManager.Get(ctx, dictCacheKey); err == nil && !cached.IsEmpty() {
		var list []v.DictFullInfo
		if json.Unmarshal([]byte(cached.String()), &list) == nil {
			return list
		}
	}
	list := loadDictData()
	if b, err := json.Marshal(list); err == nil {
		_ = v.CacheManager.Set(ctx, dictCacheKey, string(b), dictCacheTTL)
	}
	return list
}

func findDictType(list []v.DictFullInfo, typeKey string) *v.DictFullInfo {
	for i := range list {
		if list[i].Key == typeKey {
			return &list[i]
		}
	}
	return nil
}

// loadDictData 全量加载字典（类型 + 条目全字段），类型按 id、条目按 orderNum 升序
func loadDictData() []v.DictFullInfo {
	types, err := v.DBM(model.NewDictType()).Order("id asc").All()
	if err != nil {
		return nil
	}
	entriesByType := map[uint][]v.DictEntry{}
	if infos, err := v.DBM(model.NewDictInfo()).Order("orderNum asc").All(); err == nil {
		for _, row := range infos {
			entry := v.DictEntry{
				ID:       row["id"].Uint(),
				Name:     row["name"].String(),
				Value:    gvar.New(row["value"].Val()),
				OrderNum: row["orderNum"].Int(),
				Remark:   row["remark"].String(),
			}
			if !row["parentId"].IsNil() {
				pid := row["parentId"].Int()
				entry.ParentID = &pid
			}
			typeID := row["typeId"].Uint()
			entriesByType[typeID] = append(entriesByType[typeID], entry)
		}
	}
	list := make([]v.DictFullInfo, 0, len(types))
	for _, row := range types {
		typeID := row["id"].Uint()
		entries := entriesByType[typeID]
		if entries == nil {
			entries = []v.DictEntry{}
		}
		list = append(list, v.DictFullInfo{
			ID:      typeID,
			Name:    row["name"].String(),
			Key:     row["key"].String(),
			Entries: entries,
		})
	}
	return list
}

// invalidateDictCache 字典数据变更后失效缓存，下次查询惰性重建
func invalidateDictCache(ctx context.Context) {
	_, _ = v.CacheManager.Remove(ctx, dictCacheKey)
}

func init() {
	v.RegisterDict(dictLookup{})
}
