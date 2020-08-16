package teaagents

type EventInterface interface {
	AsJSON() (data []byte, err error)
}
