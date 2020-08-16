package agentutils

import (
	"github.com/TeaWeb/build/internal/teaconfigs/notices"
	"github.com/iwind/TeaGo/maps"
)

// 将Receiver转换为Map
func ConvertReceiversToMaps(receivers []*notices.NoticeReceiver) (result []maps.Map) {
	result = []maps.Map{}
	for _, receiver := range receivers {
		m := maps.Map{
			"name":      receiver.Name,
			"id":        receiver.Id,
			"user":      receiver.User,
			"mediaType": "",
		}

		// 媒介
		media := notices.SharedNoticeSetting().FindMedia(receiver.MediaId)
		if media != nil {
			m["mediaType"] = media.Name
		}
		result = append(result, m)
	}

	return result
}
