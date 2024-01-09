package service

import (
	"github.com/vera-byte/vgo/cool"
	"github.com/vera-byte/vgo/modules/demo/model"
)

type DemoGoodsService struct {
	*cool.Service
}

func NewDemoGoodsService() *DemoGoodsService {
	return &DemoGoodsService{
		&cool.Service{
			Model: model.NewDemoGoods(),
			ListQueryOp: &cool.QueryOp{

				Join: []*cool.JoinOp{},
			},
		},
	}
}
