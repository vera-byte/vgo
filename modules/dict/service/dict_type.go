package service

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/vera-byte/vgo/modules/dict/model"
	"github.com/vera-byte/vgo/v"
)

type DictTypeService struct {
	*v.Service
}

// ModifyAfter 字典数据变更后失效查询缓存
func (s *DictTypeService) ModifyAfter(ctx context.Context, method string, param g.MapStrAny) (err error) {
	invalidateDictCache(ctx)
	return
}

func NewDictTypeService() *DictTypeService {
	return &DictTypeService{
		Service: &v.Service{
			Model: model.NewDictType(),
			ListQueryOp: &v.QueryOp{
				KeyWordField: []string{"name"},
			},
		},
	}
}
