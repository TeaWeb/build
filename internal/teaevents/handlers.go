package teaevents

var eventFunctions = map[EventType][]Handler{}
var eventQueue = make(chan EventInterface, 1024)

type Handler func(event EventInterface)

// add handlers
func On(event EventType, f Handler) {
	locker.Lock()
	defer locker.Unlock()

	funcList, ok := eventFunctions[event]
	if ok {
		funcList = append(funcList, f)
	} else {
		funcList = []Handler{f}
	}
	eventFunctions[event] = funcList
}

// call handlers
func Post(event EventInterface) {
	eventQueue <- event
}
