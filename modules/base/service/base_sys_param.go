package service

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/vera-byte/vgo/modules/base/model"
	"github.com/vera-byte/vgo/v"
)

type BaseSysParamService struct {
	*v.Service
}

func NewBaseSysParamService() *BaseSysParamService {
	return &BaseSysParamService{
		&v.Service{
			Model: model.NewBaseSysParam(),
		},

		// Service: v.NewService(model.NewBaseSysParam()),
	}
}

// HtmlByKey 根据配置参数key获取网页内容(富文本)
func (s *BaseSysParamService) HtmlByKey(key string) string {
	var (
		html = "<html><body>@content</body></html>"
	)
	m := v.DBM(s.Model)
	record, err := m.Where("keyName = ?", key).One()
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
	for _, v := range result {
		key := "param:" + v["keyName"].String()
		value := v["data"].String()
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
	record, err := m.Where("keyName = ?", key).One()
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
