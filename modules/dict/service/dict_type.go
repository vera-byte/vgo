package service

import (
	"github.com/vera-byte/vgo/cool"
	"github.com/vera-byte/vgo/modules/dict/model"
)

type DictTypeService struct {
	*cool.Service
}

func NewDictTypeService() *DictTypeService {
	return &DictTypeService{
		Service: &cool.Service{
			Model: model.NewDictType(),
			ListQueryOp: &cool.QueryOp{
				KeyWordField: []string{"name"},
			},
		},
	}
}
