package agent

import (
	"encoding/base64"
	"github.com/TeaWeb/build/internal/teaconfigs/agents"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/rands"
	"net"
)

type Helper struct {
}

func (this *Helper) BeforeAction(action actions.ActionObject) bool {
	agentId := action.Header("Tea-Agent-Id")
	if len(agentId) == 0 {
		action.Fail("Authenticate Failed 001: invalid agent id")
	}

	key := action.Header("Tea-Agent-Key")
	if len(key) == 0 {
		action.Fail("Authenticate Failed 002: invalid agent key")
	}

	agent := agents.NewAgentConfigFromId(agentId)
	if agent == nil {
		// 是否自动生成
		groupKey := action.Header("Tea-Agent-Group-Key")
		if len(groupKey) > 0 && groupKey != "GroupKey" {
			groupList := agents.SharedGroupList()
			group := groupList.FindGroupWithKey(groupKey)
			if group == nil {
				action.Fail("Authenticate Failed 006: invalid group key")
			}

			// 创建新Agent
			agent := agents.NewAgentConfig()

			hostName := action.Header("Tea-Agent-Hostname")
			if len(hostName) > 0 {
				b, err := base64.StdEncoding.DecodeString(hostName)
				if err == nil && len(b) > 0 {
					agent.Name = string(b)
				}
			}
			if len(agent.Name) == 0 {
				agent.Name = action.RequestRemoteIP()
			}
			agent.Key = rands.HexString(32)
			agent.AddGroup(group.Id)
			agent.Host = action.RequestRemoteIP()
			agent.AllowAll = true
			agent.AutoUpdates = true
			agent.CheckDisconnections = true
			agent.GroupKey = groupKey
			err := agent.Save()
			if err != nil {
				logs.Error(err)
				action.Fail("Authenticate Failed 007: server error")
			}

			agentList, err := agents.SharedAgentList()
			agentList.AddAgent("agent." + agent.Id + ".conf")
			err = agentList.Save()
			if err != nil {
				logs.Error(err)

				err = agent.Delete()
				if err != nil {
					logs.Error(err)
				}
				action.Fail("Authenticate Failed 007: server error")
			}

			// 重建索引
			err = groupList.BuildIndexes()
			if err != nil {
				logs.Error(err)
			}

			action.Context.Set("agent", agent)

			return true
		}

		action.Fail("Authenticate Failed 003: agent id not found")
	}
	if agent.Id != agentId || agent.Key != key {
		action.Fail("Authenticate Failed 004: wrong agent key")
	}

	// 检查IP
	addr := action.Request.RemoteAddr
	host, _, err := net.SplitHostPort(addr)
	if err == nil {
		addr = host
	}
	if !agent.IsLocal() && !agent.AllowAll && !lists.ContainsString(agent.Allow, addr) {
		action.Fail("Access Denied 005: address denied")
	}

	action.Context.Set("agent", agent)

	return true
}
