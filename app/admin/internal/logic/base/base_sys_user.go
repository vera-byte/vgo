package base

import "context"

func init() {
	// service.RegisterBaseSysMenuLogic(NewBaseSysMenuLogic())

}

type sBaseSysUserLogic struct {
}

func NewBaseSysUserLogic() *sBaseSysUserLogic {
	return &sBaseSysUserLogic{}
}

func (s *sBaseSysUserLogic) Person(ctx context.Context, userId int64) (user interface{}, err error) {
	return nil, nil
}
