package teaevents

import (
	"github.com/iwind/TeaGo"
	"github.com/iwind/TeaGo/logs"
)

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		handleEvents()
	})
}

func handleEvents() {
	go func() {
		for event := range eventQueue {
			if event == nil {
				continue
			}
			eventType := event.Type()
			logs.Println("[event]post '" + eventType + "'")

			funcList, ok := eventFunctions[eventType]
			if !ok {
				return
			}
			for _, handler := range funcList {
				handler(event)
			}
		}
	}()
}
