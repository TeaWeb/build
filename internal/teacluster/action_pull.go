package teacluster

import (
	"github.com/TeaWeb/build/internal/teacluster/configs"
)

// node <- cluster
type PullAction struct {
	Action

	LocalItems []*configs.Item // items without data
}

func (this *PullAction) Name() string {
	return "pull"
}

func (this *PullAction) Execute() error {
	return nil
}

func (this *PullAction) TypeId() int8 {
	return ActionCodePull
}
