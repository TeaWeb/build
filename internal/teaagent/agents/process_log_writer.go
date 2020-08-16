package teaagents

// 日志写入器
type StdoutLogWriter struct {
	AppId    string
	TaskId   string
	UniqueId string
	Pid      int
}

func (this *StdoutLogWriter) Write(p []byte) (n int, err error) {
	event := NewProcessEvent(ProcessEventStdout, this.AppId, this.TaskId, this.UniqueId, this.Pid, p)
	PushEvent(event)

	n = len(p)
	return
}

type StderrLogWriter struct {
	AppId    string
	TaskId   string
	UniqueId string
	Pid      int
}

func (this *StderrLogWriter) Write(p []byte) (n int, err error) {
	event := NewProcessEvent(ProcessEventStderr, this.AppId, this.TaskId, this.UniqueId, this.Pid, p)
	PushEvent(event)

	n = len(p)
	return
}
