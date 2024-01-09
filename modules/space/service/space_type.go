package service

import (
	"github.com/vera-byte/vgo/modules/space/model"
	"github.com/vera-byte/vgo/v"
)

type SpaceTypeService struct {
	*v.Service
}

func NewSpaceTypeService() *SpaceTypeService {
	return &SpaceTypeService{
		&v.Service{
			Model: model.NewSpaceType(),
		},

		// Service: v.NewService(model.NewSpaceType()),
	}
}
