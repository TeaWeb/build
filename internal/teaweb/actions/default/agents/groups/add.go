package groups

import (
	"github.com/TeaWeb/build/internal/teaconfigs/agents"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/rands"
)

type AddAction actions.Action

// 添加分组
func (this *AddAction) Run(params struct {
	From string
}) {
	this.Data["from"] = params.From

	this.Show()
}

// 提交保存
func (this *AddAction) RunPost(params struct {
	Name string

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

	group := agents.NewGroup(params.Name)
	group.MaxAgents = params.MaxAgents
	group.DayFrom = params.DayFrom
	group.DayTo = params.DayTo

	group.Keys = []*agents.GroupKey{}
	for index, key := range params.KeysKey {
		if len(key) == 0 {
			key = rands.HexString(32)
		}
		groupKey := agents.NewGroupKey()
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
		group.AddKey(groupKey)
	}

	config := agents.SharedGroupList()
	config.AddGroup(group)
	err := config.Save()
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	this.Success()
}
