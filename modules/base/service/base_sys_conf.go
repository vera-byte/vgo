package service

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/vera-byte/vgo/modules/base/model"
	"github.com/vera-byte/vgo/v"
)

type BaseSysConfService struct {
	*v.Service
}

func NewBaseSysConfService() *BaseSysConfService {
	return &BaseSysConfService{
		&v.Service{
			Model: model.NewBaseSysConf(),
			UniqueKey: map[string]string{
				"cKey": "配置键不能重复",
			},
		},
	}
}

// UpdateValue 更新配置值
func (s *BaseSysConfService) UpdateValue(cKey, cValue string) error {
	m := v.DBM(s.Model).Where("cKey = ?", cKey)
	record, err := m.One()
	if err != nil {
		return err
	}
	if record == nil {
		_, err = v.DBM(s.Model).Insert(g.Map{
			"cKey":   cKey,
			"cValue": cValue,
		})
	} else {
		_, err = v.DBM(s.Model).Where("cKey = ?", cKey).Data(g.Map{"cValue": cValue}).Update()
	}
	return err
}

// GetValue 获取配置值
func (s *BaseSysConfService) GetValue(cKey string) string {
	m := v.DBM(s.Model).Where("cKey = ?", cKey)
	record, err := m.One()
	if err != nil {
		return ""
	}
	if record == nil {
		return ""
	}
	return record["cValue"].String()
}
