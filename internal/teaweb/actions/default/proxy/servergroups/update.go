package servergroups

import (
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/iwind/TeaGo/actions"
	"net/http"
)

type UpdateAction actions.Action

func (this *UpdateAction) RunGet(params struct {
	GroupId string
}) {
	group := teaconfigs.SharedServerGroupList().Find(params.GroupId)
	if group == nil {
		this.Error("not found", http.StatusNotFound)
		return
	}

	this.Data["group"] = group

	this.Show()
}

func (this *UpdateAction) RunPost(params struct {
	GroupId string
	Name    string

	Must *actions.Must
}) {
	params.Must.
		Field("name", params.Name).
		Require("请输入分组名称")

	groupList := teaconfigs.SharedServerGroupList()
	group := groupList.Find(params.GroupId)
	if group == nil {
		this.Fail("找不到要修改的分组")
	}
	group.Name = params.Name
	err := groupList.Save()
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	this.Success()
}
