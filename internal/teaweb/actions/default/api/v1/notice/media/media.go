package media

import (
	"github.com/TeaWeb/build/internal/teaconfigs/notices"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/api/apiutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type MediaAction actions.Action

// 媒介信息
func (this *MediaAction) RunGet(params struct {
	MediaId string
}) {
	media := notices.SharedNoticeSetting().FindMedia(params.MediaId)
	if media == nil {
		apiutils.Fail(this, "media not found")
		return
	}

	apiutils.Success(this, maps.Map{
		"config": media,
	})
}
