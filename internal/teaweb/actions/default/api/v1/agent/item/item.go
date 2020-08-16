package item

import (
	"github.com/TeaWeb/build/internal/teaconfigs/agents"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/api/apiutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type ItemAction actions.Action

// 监控项信息
func (this *ItemAction) RunGet(params struct {
	AgentId string
	AppId   string
	ItemId  string
}) {
	agent := agents.NewAgentConfigFromId(params.AgentId)
	if agent == nil {
		apiutils.Fail(this, "agent not found")
		return
	}

	app := agent.FindApp(params.AppId)
	if app == nil {
		apiutils.Fail(this, "app not found")
		return
	}

	item := app.FindItem(params.ItemId)
	if item == nil {
		apiutils.Fail(this, "item not found")
		return
	}

	apiutils.Success(this, maps.Map{
		"config": item,
	})
}
