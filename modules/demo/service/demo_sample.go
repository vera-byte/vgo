package service

import (
	"vgo-simple/modules/demo/model"

	"github.com/vera-byte/vgo/v"
)

// DemoSampleService 示例服务：通用 CRUD 由 v.Service 提供
type DemoSampleService struct {
	*v.Service
}

// NewDemoSampleService 创建示例服务
func NewDemoSampleService() *DemoSampleService {
	return &DemoSampleService{
		&v.Service{
			Model: model.NewDemoSample(),
		},
	}
}
