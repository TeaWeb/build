package servergroups

import (
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/iwind/TeaGo/actions"
)

type DeleteAction actions.Action

// 删除分组
func (this *DeleteAction) RunPost(params struct {
	GroupId string
}) {
	groupList := teaconfigs.SharedServerGroupList()
	groupList.Remove(params.GroupId)
	err := groupList.Save()
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	this.Success()
}
