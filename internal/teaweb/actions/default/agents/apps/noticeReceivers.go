package apps

import (
	"github.com/TeaWeb/build/internal/teaconfigs/agents"
	"github.com/TeaWeb/build/internal/teaconfigs/notices"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/agents/agentutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/maps"
)

type NoticeReceiversAction actions.Action

// 通知接收人设置
func (this *NoticeReceiversAction) Run(params struct {
	AgentId string
	AppId   string
}) {
	app := agentutils.InitAppData(this, params.AgentId, params.AppId, "noticeSetting")
	if app == nil {
		this.Fail("找不到要操作的App")
	}

	agent := agents.NewAgentConfigFromId(params.AgentId)
	if agent == nil {
		this.Fail("找不到Agent")
	}

	group := agent.FirstGroup()
	this.Data["groupId"] = group.Id
	this.Data["levels"] = lists.Map(notices.AllNoticeLevels(), func(k int, v interface{}) interface{} {
		level := v.(maps.Map)
		code := level["code"].(notices.NoticeLevel)

		// App设置
		receivers, found := app.NoticeSetting[code]
		if found && len(receivers) > 0 {
			level["receivers"] = agentutils.ConvertReceiversToMaps(receivers)
		} else {
			level["receivers"] = []interface{}{}
		}

		// 当前所属分组的设置
		if group != nil {
			groupReceivers, ok := group.NoticeSetting[code]
			if ok {
				level["groupReceivers"] = agentutils.ConvertReceiversToMaps(groupReceivers)
			} else {
				level["groupReceivers"] = []interface{}{}
			}
		} else {
			level["groupReceivers"] = []interface{}{}
		}

		// 当前所属Agent的设置
		agentReceivers, ok := agent.NoticeSetting[code]
		if ok {
			level["agentReceivers"] = agentutils.ConvertReceiversToMaps(agentReceivers)
		} else {
			level["agentReceivers"] = []interface{}{}
		}

		return level
	})

	this.Show()
}
