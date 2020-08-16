package teacluster

type ActionInterface interface {
	Name() string
	Execute() error
	OnSuccess(success *SuccessAction) error
	OnFail(fail *FailAction) error
	TypeId() int8
	BaseAction() *Action
}
