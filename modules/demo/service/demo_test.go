package service

import "github.com/vera-byte/vgo/v"

type DemoTestService struct {
	*v.Service
}

func NewDemoTestService() *DemoTestService {
	return &DemoTestService{
		&v.Service{},
	}
}

func (s *DemoTestService) GetDemoTestList() (interface{}, error) {
	// gsvc.SetRegistry(etcd.New(`127.0.0.1:2379`))

	return nil, nil
}
