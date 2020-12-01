package groups

import (
	"github.com/TeaWeb/build/internal/teaconfigs/agents"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/rands"
)

type UpdateAction actions.Action

// 分组ID
func (this *UpdateAction) Run(params struct {
	GroupId string
}) {
	group := agents.SharedGroupList().FindGroup(params.GroupId)
	if group == nil {
		this.Fail("找不到Group")
	}
	this.Data["group"] = group

	this.Show()
}

// 提交保存
func (this *UpdateAction) RunPost(params struct {
	GroupId   string
	Name      string
	MaxAgents int
	DayFrom   string
	DayTo     string

	KeysKey       []string
	KeysDayFrom   []string
	KeysDayTo     []string
	KeysMaxAgents []int
	KeysOn        []int
	KeysName      []string

	Must *actions.Must
}) {
	params.Must.
		Field("name", params.Name).
		Require("请输入分组名称")

	groupList := agents.SharedGroupList()
	group := groupList.FindGroup(params.GroupId)
	if group == nil {
		this.Fail("找不到Group")
	}
	group.Name = params.Name
	group.MaxAgents = params.MaxAgents
	group.DayFrom = params.DayFrom
	group.DayTo = params.DayTo

	newKeys := []*agents.GroupKey{}
	for index, key := range params.KeysKey {
		if len(key) == 0 {
			key = rands.HexString(32)
		}
		groupKey := group.FindKey(key)
		if groupKey == nil {
			groupKey = agents.NewGroupKey()
		}
		groupKey.Key = key
		if index < len(params.KeysDayFrom) {
			groupKey.DayFrom = params.KeysDayFrom[index]
		}
		if index < len(params.KeysDayTo) {
			groupKey.DayTo = params.KeysDayTo[index]
		}
		if index < len(params.KeysMaxAgents) {
			groupKey.MaxAgents = params.KeysMaxAgents[index]
		}
		if index < len(params.KeysOn) {
			groupKey.On = params.KeysOn[index] > 0
		}
		if index < len(params.KeysName) {
			groupKey.Name = params.KeysName[index]
		}
		newKeys = append(newKeys, groupKey)
	}

	group.Keys = newKeys
	err := groupList.Save()
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	this.Success()
}
