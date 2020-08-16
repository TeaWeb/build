package media

import (
	"github.com/TeaWeb/build/internal/teaconfigs/notices"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/api/apiutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type MediasAction actions.Action

// 媒介列表
func (this *MediasAction) RunGet(params struct{}) {
	setting := notices.SharedNoticeSetting()
	result := []maps.Map{}
	for _, media := range setting.Medias {
		result = append(result, maps.Map{
			"config": media,
		})
	}
	apiutils.Success(this, result)
}
