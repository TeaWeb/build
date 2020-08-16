package teaevents

const (
	EventTypeConfigChanged EventType = "EventTypeConfigChanged" // config changed
)

var configChangedEvent = new(ConfigChangedEvent)

type ConfigChangedEvent struct {
}

func NewConfigChangedEvent() *ConfigChangedEvent {
	return configChangedEvent
}

func (this *ConfigChangedEvent) Type() string {
	return EventTypeConfigChanged
}
