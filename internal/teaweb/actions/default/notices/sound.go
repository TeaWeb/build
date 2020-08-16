package notices

import (
	"github.com/TeaWeb/build/internal/teaconfigs/notices"
	"github.com/iwind/TeaGo/actions"
)

type SoundAction actions.Action

// 是否开启通知
func (this *SoundAction) Run(params struct {
	On bool
}) {
	setting := notices.SharedNoticeSetting()
	setting.SoundOn = params.On
	err := setting.Save()
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	this.Success()
}
