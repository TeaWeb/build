package notices

import (
	"github.com/TeaWeb/build/internal/teaconfigs/notices"
	"github.com/iwind/TeaGo/actions"
)

type TestMediaAction actions.Action

// 测试媒介
func (this *TestMediaAction) Run(params struct {
	MediaId string
}) {
	setting := notices.SharedNoticeSetting()
	media := setting.FindMedia(params.MediaId)
	if media == nil {
		this.Fail("找不到媒介")
	}

	this.Data["media"] = media
	this.Data["mediaType"] = notices.FindNoticeMediaType(media.Type)

	this.Show()
}

// 提交测试
func (this *TestMediaAction) RunPost(params struct {
	MediaId string
	Subject string
	Body    string
	User    string
	Must    *actions.Must
}) {
	setting := notices.SharedNoticeSetting()
	media := setting.FindMedia(params.MediaId)
	if media == nil {
		this.Fail("找不到媒介")
	}

	rawMedia, err := media.Raw()
	if err != nil {
		this.Fail("发现配置错误：" + err.Error())
	}

	this.Data["response"] = ""
	if rawMedia.RequireUser() {
		if len(params.User) == 0 {
			this.Data["error"] = "请输入用户标识"
			this.Success()
		}
	}

	resp, err := rawMedia.Send(params.User, params.Subject, params.Body)
	this.Data["response"] = string(resp)
	if err != nil {
		this.Data["error"] = err.Error()
	} else {
		this.Data["error"] = ""
	}
	this.Success()
}
