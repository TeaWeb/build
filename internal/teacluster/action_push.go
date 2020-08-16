package teacluster

import (
	"github.com/TeaWeb/build/internal/teacluster/configs"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/files"
	"github.com/iwind/TeaGo/logs"
)

// master -> cluster
type PushAction struct {
	Action

	Items []*configs.Item
}

func (this *PushAction) Name() string {
	return "push"
}

func (this *PushAction) Execute() error {
	return nil
}

func (this *PushAction) AddItem(item *configs.Item) {
	this.Items = append(this.Items, item)
}

func (this *PushAction) OnSuccess(success *SuccessAction) error {
	sumData, err := files.NewFile(Tea.ConfigFile("node.sum")).ReadAll()
	if err != nil {
		return err
	}
	err = files.NewFile(Tea.ConfigFile("cluster.sum")).Write(sumData)
	return err
}

func (this *PushAction) OnFail(fail *FailAction) error {
	logs.Println("[push]fail:", fail.Message)

	// TODO retry later

	return nil
}

func (this *PushAction) TypeId() int8 {
	return ActionCodePush
}
