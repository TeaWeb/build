package teacluster

// cluster -> slave node
type NotifyAction struct {
	Action
}

func (this *NotifyAction) Name() string {
	return "notify"
}

func (this *NotifyAction) TypeId() int8 {
	return ActionCodeNotify
}

func (this *NotifyAction) Execute() error {
	SharedManager.PullItems()
	return nil
}
