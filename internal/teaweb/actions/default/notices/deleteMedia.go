package notices

import (
	"github.com/TeaWeb/build/internal/teaconfigs/agents"
	"github.com/TeaWeb/build/internal/teaconfigs/notices"
	"github.com/iwind/TeaGo/actions"
)

type DeleteMediaAction actions.Action

// 删除媒介
func (this *DeleteMediaAction) Run(params struct {
	MediaId string
}) {
	// 删除agent group中的相关接收人
	groupConfig := agents.SharedGroupList()
	isChanged := false // 有变化才会保存
	for _, group := range groupConfig.FindAllGroups() {
		found := group.RemoveMedia(params.MediaId)
		if found {
			isChanged = true
		}
	}
	if isChanged {
		err := groupConfig.Save()
		if err != nil {
			this.Fail("删除失败：" + err.Error())
		}
	}

	// 删除agent中的相关接收人
	for _, agent := range agents.AllSharedAgents() {
		isChanged := false
		found := agent.RemoveMedia(params.MediaId)
		if found {
			isChanged = true
		}

		// app
		for _, app := range agent.Apps {
			found := app.RemoveMedia(params.MediaId)
			if found {
				isChanged = true
			}
		}

		if isChanged {
			err := agent.Save()
			if err != nil {
				this.Fail("删除失败：" + err.Error())
			}
		}
	}

	setting := notices.SharedNoticeSetting()
	setting.RemoveMedia(params.MediaId)
	err := setting.Save()
	if err != nil {
		this.Fail("删除失败：" + err.Error())
	}

	this.Success()
}
