package teaconfigs

import "github.com/TeaWeb/build/internal/teaevents"

const (
	EventBackendDown teaevents.EventType = "EventBackendDown"
	EventBackendUp   teaevents.EventType = "EventBackendUp"
)

// 后端服务器下线事件
type BackendDownEvent struct {
	Server    *ServerConfig
	Backend   *BackendConfig
	Location  *LocationConfig
	Websocket *WebsocketConfig
}

func (this *BackendDownEvent) Type() string {
	return EventBackendDown
}

// 后端服务器上线事件
type BackendUpEvent struct {
	Server    *ServerConfig
	Backend   *BackendConfig
	Location  *LocationConfig
	Websocket *WebsocketConfig
}

func (this *BackendUpEvent) Type() string {
	return EventBackendUp
}
