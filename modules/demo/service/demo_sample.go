package service

import (
	"github.com/vera-byte/vgo/modules/demo/model"
	"github.com/vera-byte/vgo/v"
)

type DemoSampleService struct {
	*v.Service
}

func NewDemoSampleService() *DemoSampleService {
	return &DemoSampleService{
		&v.Service{
			Model: model.NewDemoSample(),
		},
	}
}
