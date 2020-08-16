package agent

import (
	"encoding/base64"
	"github.com/TeaWeb/build/internal/teaconfigs/agents"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/agents/agentutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/types"
	"math"
	"time"
)

type PullAction actions.Action

// 拉取事件
func (this *PullAction) Run(params struct{}) {
	agentId := this.Context.Get("agent").(*agents.AgentConfig).Id
	agentVersion := this.Request.Header.Get("Tea-Agent-Version")
	agentOsName := this.Request.Header.Get("Tea-Agent-OsName")
	nano := this.Request.Header.Get("Tea-Agent-Nano")
	speed := float64(0)
	if len(nano) > 0 {
		speed = math.Ceil(float64(time.Now().UnixNano()-types.Int64(nano))*1000/1000000) / 1000
		if speed < 0 {
			speed = -speed
		}
	}

	osName := ""
	if len(agentOsName) > 0 {
		data, err := base64.StdEncoding.DecodeString(agentOsName)
		if err == nil {
			osName = string(data)
		}
	}

	state := agentutils.FindAgentState(agentId)
	state.IsActive = true
	state.Version = agentVersion
	state.OsName = osName
	state.Speed = speed
	state.IP = this.RequestRemoteIP()

	isDone := false

	go func() {
		<-this.Request.Context().Done()
		if !isDone {
			state.IsActive = false
		}
	}()

	event := agentutils.Wait(agentId)
	state.IsActive = false
	isDone = true

	events := []*agentutils.AgentEvent{}
	if event != nil {
		events = append(events, event)
	}

	this.Data["events"] = events

	this.Success()
}
