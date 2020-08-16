package item

import (
	"github.com/TeaWeb/build/internal/teaconfigs/agents"
	"github.com/TeaWeb/build/internal/teadb"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/api/apiutils"
	"github.com/iwind/TeaGo/actions"
)

type LatestAction actions.Action

// 获取最后一次数据
func (this *LatestAction) RunGet(params struct {
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

	value, err := teadb.AgentValueDAO().FindLatestItemValue(params.AgentId, params.AppId, item.Id)
	if err != nil {
		apiutils.Fail(this, "no value yet")
		return
	}

	if value == nil {
		apiutils.Success(this, nil)
		return
	}
	apiutils.Success(this, value.Value)
}
