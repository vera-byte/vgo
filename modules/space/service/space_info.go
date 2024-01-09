package service

import (
	"github.com/vera-byte/vgo/cool"
	"github.com/vera-byte/vgo/modules/space/model"
)

type SpaceInfoService struct {
	*cool.Service
}

func NewSpaceInfoService() *SpaceInfoService {
	return &SpaceInfoService{
		&cool.Service{
			Model: model.NewSpaceInfo(),
		},

		// Service: cool.NewService(model.NewSpaceInfo()),
	}
}
