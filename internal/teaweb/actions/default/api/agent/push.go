package agent

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/TeaWeb/build/internal/teaconfigs/agents"
	"github.com/TeaWeb/build/internal/teaconfigs/notices"
	"github.com/TeaWeb/build/internal/teadb"
	"github.com/TeaWeb/build/internal/teadb/shared"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/agents/agentutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/types"
	"github.com/iwind/TeaGo/utils/time"
	"io/ioutil"
	"strings"
	"time"
)

type PushAction actions.Action

// 接收推送的数据
func (this *PushAction) Run(params struct{}) {
	if !teadb.SharedDB().IsAvailable() {
		this.Success()
	}

	agent := this.Context.Get("agent").(*agents.AgentConfig)

	// 是否未启用
	if !agent.On {
		this.Success()
	}

	data, err := ioutil.ReadAll(this.Request.Body)
	if err != nil {
		logs.Error(err)
		this.Fail("read body error")
	}

	m := maps.Map{}
	err = json.Unmarshal(data, &m)
	if err != nil {
		logs.Error(err)
		this.Fail("unmarshal error")
	}

	timestamp := m.GetInt64("timestamp")
	t := time.Unix(timestamp, 0)

	eventDomain := m.GetString("event")

	if eventDomain == "ProcessEvent" { // 进程事件
		event := &agents.ProcessLog{
			Id:         shared.NewObjectId(),
			AgentId:    agent.Id,
			TaskId:     m.GetString("taskId"),
			ProcessId:  m.GetString("uniqueId"),
			ProcessPid: m.GetInt("pid"),
			EventType:  m.GetString("eventType"),
			Data:       m.GetString("data"),
			Timestamp:  timestamp,
			TimeFormat: struct {
				Year   string `bson:"year" json:"year"`
				Month  string `bson:"month" json:"month"`
				Week   string `bson:"week" json:"week"`
				Day    string `bson:"day" json:"day"`
				Hour   string `bson:"hour" json:"hour"`
				Minute string `bson:"minute" json:"minute"`
				Second string `bson:"second" json:"second"`
			}{
				Year:   timeutil.Format("Y", t),
				Month:  timeutil.Format("Ym", t),
				Week:   timeutil.Format("YW", t),
				Day:    timeutil.Format("Ymd", t),
				Hour:   timeutil.Format("YmdH", t),
				Minute: timeutil.Format("YmdHi", t),
				Second: timeutil.Format("YmdHis", t),
			},
		}

		err = teadb.AgentLogDAO().InsertOne(agent.Id, event)
		if err != nil {
			logs.Error(err)
		}
	} else if eventDomain == "ItemEvent" { // 监控项事件
		this.processItemEvent(agent, m, t)
	}

	this.Success()
}

// 处理监控项事件
func (this *PushAction) processItemEvent(agent *agents.AgentConfig, m maps.Map, t time.Time) {
	appId := m.GetString("appId")
	itemId := m.GetString("itemId")
	app := agent.FindApp(appId)
	if app == nil {
		this.Success()
	}

	item := app.FindItem(itemId)
	if item == nil {
		this.Success()
	}

	v := m.Get("value")
	oldValue, err := this.findLatestAgentValue(agent.Id, appId, itemId)
	if err != nil {
		if err != context.DeadlineExceeded {
			logs.Error(err)
		}
		return
	}
	if oldValue == nil {
		oldValue = v
	}
	threshold, row, level, message, err := item.TestValue(v, oldValue)
	if err != nil {
		logs.Error(errors.New(item.Name + " " + err.Error()))
		if len(m.GetString("error")) == 0 {
			m["error"] = err.Error()
		}
	}

	// 处理消息中的变量
	message = agents.RegexpParamNamedVariable.ReplaceAllStringFunc(message, func(s string) string {
		result, err := agents.EvalParam(s, v, oldValue, maps.Map{
			"AGENT": maps.Map{
				"name": agent.Name,
				"host": agent.Host,
			},
			"APP": maps.Map{
				"name": app.Name,
			},
			"ITEM": maps.Map{
				"name": item.Name,
			},
			"ROW": row,
		}, false)
		if err != nil {
			logs.Error(err)
		}
		return result
	})

	if threshold != nil && len(message) == 0 {
		message = threshold.Expression()
	}

	// 通知消息
	setting := notices.SharedNoticeSetting()

	isNotified := false

	if level != notices.NoticeLevelNone {
		// 是否发送通知
		shouldNotify := true

		// 检查最近N此数值是否都是同类错误
		if threshold != nil && threshold.MaxFails > 1 {
			values, err := teadb.AgentValueDAO().ListItemValues(agent.Id, app.Id, item.Id, 0, "", 0, threshold.MaxFails-1)
			if err != nil {
				logs.Error(err)
			} else {
				if len(values) != threshold.MaxFails-1 { // 未达到连续失败次数
					shouldNotify = false
				} else {
					for _, v := range values {
						if v.ThresholdId != threshold.Id || v.IsNotified {
							shouldNotify = false
							break
						}
					}
				}
			}
		}

		// 发送通知
		if shouldNotify {
			isNotified = true

			// 添加状态
			agentutils.AddItemState(item.Id, &agentutils.ItemState{
				IsFailed: true,
			})

			notice := notices.NewNotice()
			notice.SetTime(t)
			notice.Message = message
			notice.Agent = notices.AgentCond{
				AgentId: agent.Id,
				AppId:   appId,
				ItemId:  itemId,
				Level:   level,
			}
			if threshold != nil {
				notice.Agent.Threshold = threshold.Expression()
			}
			notice.Hash()

			if notices.IsFailureLevel(level) {
				// 同样的消息短时间内只发送一条
				b, err := teadb.NoticeDAO().ExistNoticesWithHash(notice.MessageHash, map[string]interface{}{
					"agent.agentId": agent.Id,
					"agent.appId":   appId,
					"agent.itemId":  itemId,
				}, 1*time.Hour)
				if err != nil {
					logs.Error(err)
				}
				if b {
					shouldNotify = false
				}
			}

			if shouldNotify {
				err := teadb.NoticeDAO().InsertOne(notice)
				if err != nil {
					logs.Error(err)
				}

				// 通过媒介发送通知
				fullMessage := "消息：" + message + "\n时间：" + timeutil.Format("Y-m-d H:i:s", t)
				linkNames := []string{}
				for _, l := range agentutils.FindNoticeLinks(notice) {
					linkNames = append(linkNames, types.String(l["name"]))
				}
				if len(linkNames) > 0 {
					fullMessage += "\n位置：" + strings.Join(linkNames, "/")
				}

				receiverIds := this.notifyMessage(agent, appId, itemId, setting, level, "有新的通知", fullMessage, false)
				if len(receiverIds) > 0 {
					err = teadb.NoticeDAO().UpdateNoticeReceivers(notice.Id.Hex(), receiverIds)
					if err != nil {
						logs.Error(err)
					}
				}
			}
		}
	}

	// 数值记录
	node := teaconfigs.SharedNodeConfig()
	nodeId := ""
	if node != nil {
		nodeId = node.Id
	}
	beginAt := m.GetInt64("beginAt")
	if beginAt == 0 {
		beginAt = time.Now().Unix()
	}
	value := &agents.Value{
		Id:          shared.NewObjectId(),
		NodeId:      nodeId,
		AppId:       appId,
		AgentId:     agent.Id,
		ItemId:      itemId,
		Value:       v,
		Error:       m.GetString("error"),
		NoticeLevel: level,
		CreatedAt:   beginAt,
		CostMs:      m.GetFloat64("costMs"),
		IsNotified:  isNotified,
	}
	if threshold != nil {
		value.ThresholdId = threshold.Id
		value.Threshold = threshold.Expression()
	}
	value.SetTime(t)

	err = teadb.AgentValueDAO().Insert(agent.Id, value)
	if err != nil {
		logs.Error(err)
		return
	}

	// 是否发送恢复通知
	if len(item.Id) > 0 && !notices.IsFailureLevel(level) {
		// 是否已经存在
		itemState, ok := agentutils.FindItemState(item.Id)

		if !ok || !itemState.IsFailed {
			return
		}

		recoverSuccesses := item.RecoverSuccesses
		if recoverSuccesses <= 0 {
			recoverSuccesses = 1
		}

		values, err := teadb.AgentValueDAO().ListItemValues(agent.Id, app.Id, item.Id, 0, "", 0, recoverSuccesses+1)
		if err != nil {
			logs.Error(err)
			return
		}

		lists.Reverse(values)

		countValues := len(values)
		if countValues != recoverSuccesses+1 {
			return
		}

		if !notices.IsFailureLevel(values[0].NoticeLevel) {
			return
		}

		success := true
		for i := 1; i < countValues; i++ {
			if notices.IsFailureLevel(values[i].NoticeLevel) {
				success = false
				break
			}
		}

		if success {
			// 删除错误的状态
			agentutils.RemoveItemState(item.Id)

			// 发送成功级别的通知
			notice := notices.NewNotice()
			notice.SetTime(t)
			notice.Agent = notices.AgentCond{
				AgentId: agent.Id,
				AppId:   appId,
				ItemId:  itemId,
				Level:   notices.NoticeLevelSuccess,
			}
			notice.Message = "监控项经过" + fmt.Sprintf("%d", recoverSuccesses) + "次刷新后，判定已恢复正常"
			linkNames := []string{}
			for _, l := range agentutils.FindNoticeLinks(notice) {
				linkNames = append(linkNames, types.String(l["name"]))
			}
			if len(linkNames) > 0 {
				notice.Message += " \n位置：" + strings.Join(linkNames, "/")
			}
			notice.Hash()
			err := teadb.NoticeDAO().InsertOne(notice)
			if err != nil {
				logs.Error(err)
			}

			this.notifyMessage(agent, appId, itemId, setting, notices.NoticeLevelSuccess, "有新的通知", notice.Message, true)
		}
	}
}

// 发送通知消息
func (this *PushAction) notifyMessage(agent *agents.AgentConfig, appId string, itemId string, setting *notices.NoticeSetting, level notices.NoticeLevel, subject string, message string, isSuccess bool) (receiverIds []string) {
	receiverLevels := []notices.NoticeLevel{level}
	receivers := this.findNoticeReceivers(agent, appId, setting, receiverLevels)
	if len(receivers) == 0 && isSuccess {
		receiverLevels = append(receiverLevels, notices.NoticeLevelError, notices.NoticeLevelWarning)
		receivers = this.findNoticeReceivers(agent, appId, setting, receiverLevels)
	}
	if len(receivers) == 0 {
		return []string{}
	}

	receiverIds = setting.NotifyReceivers(level, receivers, "["+agent.GroupName()+"]["+agent.Name+"]"+subject, message, func(receiverId string, minutes int) int {
		count, err := teadb.NoticeDAO().CountReceivedNotices(receiverId, map[string]interface{}{
			"agent.agentId": agent.Id,
			"agent.appId":   appId,
			"agent.itemId":  itemId,
		}, minutes)
		if err != nil {
			logs.Error(err)
		}
		return count
	})

	return receiverIds
}

// 查找最近的一次数值记录
func (this *PushAction) findLatestAgentValue(agentId string, appId string, itemId string) (interface{}, error) {
	v, err := teadb.AgentValueDAO().FindLatestItemValueNoError(agentId, appId, itemId)
	if err != nil {
		return nil, err
	}
	if v == nil {
		return nil, nil
	}
	return v.Value, nil
}

// 查找接收者
func (this *PushAction) findNoticeReceivers(agent *agents.AgentConfig, appId string, setting *notices.NoticeSetting, receiverLevels []notices.NoticeLevel) (receivers []*notices.NoticeReceiver) {
	if agent == nil {
		return
	}

	// 查找App的通知设置
	app := agent.FindApp(appId)
	if app != nil {
		receivers = app.FindAllNoticeReceivers(receiverLevels...)
		if len(receivers) > 0 {
			return
		}
	}

	// 查找Agent的通知设置
	receivers = agent.FindAllNoticeReceivers(receiverLevels...)
	if len(receivers) > 0 {
		return receivers
	}

	// 查找分组的通知设置
	groupId := "default"
	if len(agent.GroupIds) > 0 {
		groupId = agent.GroupIds[0]
	}
	group := agents.SharedGroupList().FindGroup(groupId)
	if group != nil {
		receivers = group.FindAllNoticeReceivers(receiverLevels...)
		if len(receivers) > 0 {
			return
		}
	}

	// 全局通知
	if setting != nil {
		receivers = setting.FindAllNoticeReceivers(receiverLevels...)
		if len(receivers) > 0 {
			return
		}
	}

	return
}
