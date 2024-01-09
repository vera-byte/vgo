package service

import (
	"github.com/vera-byte/vgo/cool"
	"github.com/vera-byte/vgo/modules/space/model"
)

type SpaceTypeService struct {
	*cool.Service
}

func NewSpaceTypeService() *SpaceTypeService {
	return &SpaceTypeService{
		&cool.Service{
			Model: model.NewSpaceType(),
		},

		// Service: cool.NewService(model.NewSpaceType()),
	}
}
