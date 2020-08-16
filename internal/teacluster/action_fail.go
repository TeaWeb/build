package teacluster

type FailAction struct {
	Action

	Message string
}

func (this *FailAction) Name() string {
	return "fail"
}

func (this *FailAction) TypeId() int8 {
	return ActionCodeFail
}
