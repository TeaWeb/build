package teacluster

import (
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/iwind/TeaGo/logs"
)

type RegisterAction struct {
	Action

	ClusterId     string
	ClusterSecret string
	NodeId        string
	NodeName      string
	NodeRole      string
}

func (this *RegisterAction) Name() string {
	return "register"
}

func (this *RegisterAction) Execute() error {
	return nil
}

func (this *RegisterAction) OnSuccess(success *SuccessAction) error {
	if this.NodeRole == teaconfigs.NodeRoleMaster {
		logs.Println("[cluster]register master ok")
		SharedManager.Write(&SumAction{})
	} else {
		logs.Println("[cluster]register node ok")
		SharedManager.Write(&SumAction{})
	}
	return nil
}

func (this *RegisterAction) OnFail(fail *FailAction) error {
	logs.Println("[cluster]fail to register node:", fail.Message)
	return nil
}

func (this *RegisterAction) TypeId() int8 {
	return ActionCodeRegister
}
