package settings

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
}) {
	this.Data["selectedTab"] = "noticeSetting"

	agent := agents.NewAgentConfigFromId(params.AgentId)
	if agent == nil {
		this.Fail("找不到Agent")
	}
	this.Data["agent"] = agent

	group := agent.FirstGroup()
	this.Data["groupId"] = group.Id
	this.Data["levels"] = lists.Map(notices.AllNoticeLevels(), func(k int, v interface{}) interface{} {
		level := v.(maps.Map)
		code := level["code"].(notices.NoticeLevel)
		receivers, found := agent.NoticeSetting[code]

		// 当前Agent的设置
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

		return level
	})

	this.Show()
}
