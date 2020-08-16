package settings

import (
	"github.com/TeaWeb/build/internal/teaconfigs/agents"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/agents/agentutils"
	"github.com/iwind/TeaGo/actions"
	"net/http"
)

type Helper struct {
}

func (this *Helper) BeforeAction(actionWrapper actions.ActionWrapper) {
	action := actionWrapper.Object()
	if action.Request.Method != http.MethodGet {
		return
	}

	action.Data["selectedTab"] = "detail"
	action.Data["countNoticeReceivers"] = 0

	agentId := action.ParamString("agentId")
	if len(agentId) > 0 {
		agent := agents.NewAgentConfigFromId(agentId)
		if agent != nil {
			action.Data["countNoticeReceivers"] = agent.CountNoticeReceivers()
		}
	}

	agentutils.AddTabbar(action)
}
