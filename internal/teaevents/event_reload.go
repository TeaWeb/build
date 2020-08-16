package teaevents

const (
	EventTypeReload EventType = "EventTypeReload" // reload system
)

var reloadEvent = new(ReloadEvent)

type ReloadEvent struct {
}

func NewReloadEvent() *ReloadEvent {
	return reloadEvent
}

func (this *ReloadEvent) Type() string {
	return EventTypeReload
}
