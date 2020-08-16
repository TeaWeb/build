package notices

import (
	"github.com/TeaWeb/build/internal/teaconfigs/notices"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type UpdateReceiverAction actions.Action

// 修改接收者
func (this *UpdateReceiverAction) Run(params struct {
	ReceiverId string
}) {
	setting := notices.SharedNoticeSetting()
	level, receiver := setting.FindReceiver(params.ReceiverId)
	if receiver == nil {
		this.Fail("找不到要修改的Receiver")
	}

	this.Data["receiver"] = receiver
	this.Data["level"] = notices.FindNoticeLevel(level)

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
func (this *UpdateReceiverAction) RunPost(params struct {
	ReceiverId string
	On         bool
	Name       string
	MediaId    string
	User       string
	Must       actions.Must
}) {
	params.Must.
		Field("name", params.Name).
		Require("请输入接收人名称").
		Field("mediaId", params.MediaId).
		Require("请选择使用的媒介")

	setting := notices.SharedNoticeSetting()
	_, receiver := setting.FindReceiver(params.ReceiverId)
	if receiver == nil {
		this.Fail("找不到要修改的Receiver")
	}

	// 是否校验接收人
	media := setting.FindMedia(params.MediaId)
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

	receiver.On = params.On
	receiver.Name = params.Name
	receiver.MediaId = params.MediaId
	receiver.User = params.User

	err = setting.Save()
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	this.Success()
}
