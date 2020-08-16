package servergroups

import (
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/iwind/TeaGo/actions"
)

type DeleteServerAction actions.Action

func (this *DeleteServerAction) RunPost(params struct {
	GroupId  string
	ServerId string
}) {
	groupList := teaconfigs.SharedServerGroupList()
	group := groupList.Find(params.GroupId)
	if group == nil {
		this.Fail("找不到要操作的组")
	}

	group.Remove(params.ServerId)
	err := groupList.Save()
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	this.Success()
}
