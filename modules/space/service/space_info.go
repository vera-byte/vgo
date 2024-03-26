package service

import (
	"github.com/vera-byte/vgo/modules/space/model"
	"github.com/vera-byte/vgo/v"
)

type SpaceInfoService struct {
	*v.Service
}

func NewSpaceInfoService() *SpaceInfoService {
	return &SpaceInfoService{
		&v.Service{
			Model: model.NewSpaceInfo(),
		},

		// Service: v.NewService(model.NewSpaceInfo()),
	}
}
