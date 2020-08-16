package teacluster

import (
	"time"
)

type Action struct {
	Id          uint64
	RequestTime time.Time
	RequestId   uint64
}

func (this *Action) Execute() error {
	return nil
}

func (this *Action) OnSuccess(success *SuccessAction) error {
	return nil
}

func (this *Action) OnFail(fail *FailAction) error {
	return nil
}

func (this *Action) BaseAction() *Action {
	return this
}
