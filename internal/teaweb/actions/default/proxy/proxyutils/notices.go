package proxyutils

import (
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/TeaWeb/build/internal/teaconfigs/notices"
	"github.com/TeaWeb/build/internal/teadb"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/notices/noticeutils"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"strings"
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

// 发送一个后端下线通知
func NotifyProxyBackendDownMessage(event *teaconfigs.BackendDownEvent) error {
	server := event.Server
	server.SetupNoticeItems()

	noticeItem := server.NoticeItems.BackendDown
	if noticeItem == nil || !noticeItem.On {
		return nil
	}

	positions := []string{}
	cond := notices.ProxyCond{
		ServerId:  server.Id,
		BackendId: event.Backend.Id,
		Level:     noticeItem.Level,
	}
	positions = append(positions, server.Description)
	if event.Location != nil {
		cond.LocationId = event.Location.Id
		positions = append(positions, "\"" + event.Location.Pattern + "\"")
	}
	if event.Websocket != nil {
		cond.Websocket = true
		positions = append(positions, "Websocket")
	}

	params := maps.Map{
		"server.description": server.Description,
		"backend.address":    event.Backend.Address,
		"position":           strings.Join(positions, " / "),
		"cause":              "错误过多",
	}

	// 不阻塞
	go func() {
		err := teadb.NoticeDAO().NotifyProxyMessage(cond, noticeItem.FormatBody(params))
		if err != nil {
			logs.Error(err)
		}
	}()

	NotifyServer(server, noticeItem.Level, noticeItem.FormatSubject(params), noticeItem.FormatBody(params))

	return nil
}

// 发送一个后端上线通知
func NotifyProxyBackendUpMessage(event *teaconfigs.BackendUpEvent) error {
	server := event.Server
	server.SetupNoticeItems()

	noticeItem := server.NoticeItems.BackendUp
	if noticeItem == nil || !noticeItem.On {
		return nil
	}

	positions := []string{}
	cond := notices.ProxyCond{
		ServerId:  server.Id,
		BackendId: event.Backend.Id,
		Level:     noticeItem.Level,
	}
	positions = append(positions, server.Description)
	if event.Location != nil {
		cond.LocationId = event.Location.Id
		positions = append(positions, "\"" + event.Location.Pattern + "\"")
	}
	if event.Websocket != nil {
		cond.Websocket = true
		positions = append(positions, "Websocket")
	}

	params := maps.Map{
		"server.description": server.Description,
		"backend.address":    event.Backend.Address,
		"position":           strings.Join(positions, " / "),
		"cause":              "重新检测成功",
	}

	// 不阻塞
	go func() {
		err := teadb.NoticeDAO().NotifyProxyMessage(cond, noticeItem.FormatBody(params))
		if err != nil {
			logs.Error(err)
		}
	}()

	NotifyServer(server, noticeItem.Level, noticeItem.FormatSubject(params), noticeItem.FormatBody(params))

	return nil
}

// 推送代理服务相关通知
func NotifyServer(server *teaconfigs.ServerConfig, level notices.NoticeLevel, subject string, message string) {
	receivers := server.FindAllNoticeReceivers(level)
	if len(receivers) == 0 {
		setting := notices.SharedNoticeSetting()
		if setting != nil {
			receivers = setting.FindAllNoticeReceivers(level)
		}
		if len(receivers) == 0 {
			return
		}
	}
	noticeutils.AddTask(level, receivers, subject, message)
}
