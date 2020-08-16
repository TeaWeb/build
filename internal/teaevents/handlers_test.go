package teaevents

import (
	"testing"
	"time"
)

func TestHandlers(t *testing.T) {
	handleEvents()

	On(EventTypeReload, func(event EventInterface) {
		t.Logf("reload1: %#v", event)
	})
	On(EventTypeReload, func(event EventInterface) {
		t.Logf("reload2: %#v", event)
	})
	Post(NewReloadEvent())

	On(EventTypeConfigChanged, func(event EventInterface) {
		t.Logf("config changed: %#v", event)
	})
	Post(NewConfigChangedEvent())

	time.Sleep(1 * time.Second)
}
