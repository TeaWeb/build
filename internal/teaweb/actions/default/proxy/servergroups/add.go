package servergroups

import (
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/iwind/TeaGo/actions"
)

type AddAction actions.Action

// 添加分组
func (this *AddAction) RunGet(params struct{}) {
	this.Data["selectedMenu"] = "add"

	this.Show()
}

func (this *AddAction) RunPost(params struct {
	Name string

	Must *actions.Must
}) {
	params.Must.
		Field("name", params.Name).
		Require("请输入分组名称")

	group := teaconfigs.NewServerGroup()
	group.Name = params.Name

	groupList := teaconfigs.SharedServerGroupList()
	groupList.Add(group)
	err := groupList.Save()
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	this.Success()
}
