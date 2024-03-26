package service

import (
	"github.com/vera-byte/vgo/modules/dict/model"
	"github.com/vera-byte/vgo/v"
)

type DictTypeService struct {
	*v.Service
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
