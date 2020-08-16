package teacluster

import "github.com/iwind/TeaGo/maps"

type SuccessAction struct {
	Action

	Message string
	Data    maps.Map
}

func (this *SuccessAction) Name() string {
	return "success"
}

func (this *SuccessAction) TypeId() int8 {
	return ActionCodeSuccess
}
