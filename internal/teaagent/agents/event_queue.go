package teaagents

var eventQueue = make(chan EventInterface, 1024)

func PushEvent(event EventInterface) {
	eventQueue <- event
}
