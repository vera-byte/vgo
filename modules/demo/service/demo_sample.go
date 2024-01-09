package service

import (
	"github.com/vera-byte/vgo/cool"
	"github.com/vera-byte/vgo/modules/demo/model"
)

type DemoSampleService struct {
	*cool.Service
}

func NewDemoSampleService() *DemoSampleService {
	return &DemoSampleService{
		&cool.Service{
			Model: model.NewDemoSample(),
		},
	}
}
