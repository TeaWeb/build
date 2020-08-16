package groups

import (
	"github.com/TeaWeb/build/internal/teaconfigs/agents"
	"github.com/TeaWeb/build/internal/teaconfigs/notices"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type AddNoticeReceiversAction actions.Action

// 添加通知
func (this *AddNoticeReceiversAction) Run(params struct {
	GroupId string
	Level   notices.NoticeLevel
}) {
	group := agents.SharedGroupList().FindGroup(params.GroupId)
	if group == nil {
		this.Fail("Group不存在")
	}

	level := notices.FindNoticeLevel(params.Level)
	if level == nil {
		this.Fail("Level不存在")
	}

	this.Data["group"] = group
	this.Data["level"] = level

	// 媒介
	setting := notices.SharedNoticeSetting()
	mediaMaps := []maps.Map{}
	for _, media := range setting.Medias {
		if !media.On {
			continue
		}
		mediaType := notices.FindNoticeMediaType(media.Type)
		if mediaType == nil {
			continue
		}
		mediaMaps = append(mediaMaps, maps.Map{
			"id":               media.Id,
			"name":             media.Name,
			"typeName":         mediaType["name"],
			"type":             media.Type,
			"mediaDescription": mediaType["description"],
			"userDescription":  mediaType["user"],
		})
	}
	this.Data["medias"] = mediaMaps

	this.Show()
}

// 提交保存
func (this *AddNoticeReceiversAction) RunPost(params struct {
	GroupId string
	Level   notices.NoticeLevel
	On      bool
	Name    string
	MediaId string
	User    string
	Must    actions.Must
}) {
	params.Must.
		Field("name", params.Name).
		Require("请输入接收人名称").
		Field("mediaId", params.MediaId).
		Require("请选择使用的媒介")

	// 是否校验接收人
	media := notices.SharedNoticeSetting().FindMedia(params.MediaId)
	if media == nil {
		this.Fail("找不到媒介")
	}
	rawMedia, err := media.Raw()
	if err != nil {
		this.Fail("找不到媒介：" + err.Error())
	}
	if rawMedia.RequireUser() {
		params.Must.
			Field("name", params.Name).
			Require("请输入接收人名称")
	}

	receiver := notices.NewNoticeReceiver()
	receiver.On = params.On
	receiver.Name = params.Name
	receiver.MediaId = params.MediaId
	receiver.User = params.User

	config := agents.SharedGroupList()
	group := config.FindGroup(params.GroupId)
	if group == nil {
		this.Fail("找不到Group")
	}

	group.AddNoticeReceiver(params.Level, receiver)
	err = config.Save()
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	this.Success()
}
