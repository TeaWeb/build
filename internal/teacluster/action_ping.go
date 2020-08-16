package teacluster

// node -> cluster
type PingAction struct {
	Action

	Version int64
}

func (this *PingAction) Name() string {
	return "ping"
}

func (this *PingAction) TypeId() int8 {
	return ActionCodePing
}
