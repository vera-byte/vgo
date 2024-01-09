package service

import (
	"github.com/vera-byte/vgo/cool"
)

type DemoTestService struct {
	*cool.Service
}

func NewDemoTestService() *DemoTestService {
	return &DemoTestService{
		&cool.Service{},
	}
}

func (s *DemoTestService) GetDemoTestList() (interface{}, error) {
	// gsvc.SetRegistry(etcd.New(`127.0.0.1:2379`))

	return nil, nil
}
