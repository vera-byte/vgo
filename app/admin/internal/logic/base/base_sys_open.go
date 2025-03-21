package base

func init() {
	// service.RegisterBaseSysOpenLogic(NewBaseSysOpenLogic())
}

type sBaseSysOpenLogic struct {
}

func NewBaseSysOpenLogic() *sBaseSysOpenLogic {
	return &sBaseSysOpenLogic{}
}
