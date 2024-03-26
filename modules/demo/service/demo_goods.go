package service

import (
	"github.com/vera-byte/vgo/modules/demo/model"
	"github.com/vera-byte/vgo/v"
)

type DemoGoodsService struct {
	*v.Service
}

func NewDemoGoodsService() *DemoGoodsService {
	return &DemoGoodsService{
		&v.Service{
			Model: model.NewDemoGoods(),
			ListQueryOp: &v.QueryOp{

				Join: []*v.JoinOp{},
			},
		},
	}
}
