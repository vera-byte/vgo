package service

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/vera-byte/vgo/modules/base/model"
	"github.com/vera-byte/vgo/v"
	"gorm.io/gorm"
)

type BaseSysParamService struct {
	*v.Service
}

func NewBaseSysParamService() *BaseSysParamService {
	return &BaseSysParamService{
		&v.Service{
			Model: model.NewBaseSysParam(),
		},
	}
}

// HtmlByKey 根据配置参数key获取网页内容(富文本)
func (s *BaseSysParamService) HtmlByKey(key string) string {
	var (
		html = "<html><body>@content</body></html>"
	)
	m := v.DBM(s.Model)
	record, err := m.Where("keyName", key).One()
	if err != nil {
		html = gstr.Replace(html, "@content", err.Error())
		return html
	}
	if record.IsEmpty() {
		html = gstr.Replace(html, "@content", "keyName notfound")
		return html
	}
	html = gstr.Replace(html, "@content", record["data"].String())

	return html
}

// ModifyAfter 修改后
func (s *BaseSysParamService) ModifyAfter(ctx context.Context, method string, param g.MapStrAny) (err error) {
	var (
		m = v.DBM(s.Model)
	)
	result, err := m.All()
	if err != nil {
		return
	}
	for _, item := range result {
		key := "param:" + item["keyName"].String()
		value := item["data"].String()
		err = v.CacheManager.Set(ctx, key, value, 0)
		if err != nil {
			return
		}
	}
	return
}

// DataByKey 根据配置参数key获取数据
func (s *BaseSysParamService) DataByKey(ctx context.Context, key string) (data string, err error) {
	var (
		m = v.DBM(s.Model)
	)
	rKey := "param:" + key
	dataCache, err := v.CacheManager.Get(ctx, rKey)
	if err != nil {
		return
	}
	if !dataCache.IsEmpty() {
		data = dataCache.String()
		return
	}
	record, err := m.Where("keyName", key).One()
	if err != nil {
		return
	}
	if record.IsEmpty() {
		return
	}
	data = record["data"].String()
	err = v.CacheManager.Set(ctx, rKey, data, 0)
	return
}

// ModifyBefore 新增|修改前校验 keyName 唯一（与 TS 在 modifyBefore 中定制校验一致）
func (s *BaseSysParamService) ModifyBefore(ctx context.Context, method string, param g.MapStrAny) (err error) {
	if method != "Add" && method != "Update" {
		return
	}
	var entity *model.BaseSysParam
	gconv.Scan(param, &entity)
	if entity == nil {
		return
	}
	query := v.DBM(s.Model).Where("keyName", entity.KeyName)

	var check *model.BaseSysParam
	err = query.Scan(&check)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if check != nil && method == "Add" {
		return errors.New("存在相同的keyName")
	}
	return
}

func (s *BaseSysParamService) ServiceInfo(ctx context.Context, req *v.InfoReq) (data interface{}, err error) {
	result, err := s.Service.ServiceInfo(ctx, req)
	if err != nil {
		return nil, err
	}
	if result.(gdb.Record).IsEmpty() {
		return nil, nil
	}
	resultMap := result.(gdb.Record).Map()

	str, ok := resultMap["data"].(string)
	if !ok {
		return resultMap, nil
	}

	// 先尝试直接解析（原本就是合法 JSON 的情况）
	var parsed interface{}
	if json.Unmarshal([]byte(str), &parsed) == nil {
		resultMap["data"] = parsed
		return resultMap, nil
	}

	// 再尝试替换后解析（兼容自定义格式）
	replaced := strings.ReplaceAll(str, "{", "[")
	replaced = strings.ReplaceAll(replaced, "}", "]")
	if json.Unmarshal([]byte(replaced), &parsed) == nil {
		resultMap["data"] = parsed
		return resultMap, nil
	}
	// 非 JSON 的纯文本/HTML 原样返回
	return resultMap, nil
}
