package groups

import (
	"github.com/TeaWeb/build/internal/teaconfigs/agents"
	"github.com/iwind/TeaGo/actions"
)

type GenerateKeyAction actions.Action

// 生成密钥
func (this *GenerateKeyAction) RunPost(params struct {
	GroupId string
}) {
	config := agents.SharedGroupList()
	group := config.FindGroup(params.GroupId)
	if group == nil {
		this.Fail("找不到Group")
	}
	group.Key = group.GenerateKey()
	err := config.Save()
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	this.Success()
}
