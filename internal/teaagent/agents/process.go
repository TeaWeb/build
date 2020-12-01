package teaagents

import (
	"bytes"
	"github.com/TeaWeb/build/internal/teaconfigs/shared"
	"github.com/iwind/TeaGo/rands"
	"io"
	"os"
	"os/exec"
)

type Process struct {
	UniqueId string
	Pid      int

	Env  []*shared.Variable
	Cwd  string
	File string

	proc    *os.Process
	onStart func()
	onStop  func()
}

// 获取新进程
func NewProcess() *Process {
	return &Process{
		UniqueId: rands.HexString(16),
	}
}

// on start
func (this *Process) OnStart(f func()) {
	this.onStart = f
}

// on stop
func (this *Process) OnStop(f func()) {
	this.onStop = f
}

// 立即运行
func (this *Process) Run() (stdout string, stderr string, err error) {
	stdoutWriter := bytes.NewBuffer([]byte{})
	stderrWriter := bytes.NewBuffer([]byte{})
	err = this.RunWriter(stdoutWriter, stderrWriter)

	stdout = stdoutWriter.String()
	stderr = stderrWriter.String()
	if this.onStop != nil {
		this.onStop()
	}

	return
}

// 使用自定义stdout, stderr运行
func (this *Process) RunWriter(stdout io.Writer, stderr io.Writer) (err error) {
	// execute
	cmd := exec.Command(this.File)
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	// cwd
	if len(this.Cwd) > 0 {
		cmd.Dir = this.Cwd
	}

	// env
	for _, envVar := range this.Env {
		cmd.Env = append(cmd.Env, envVar.Name+"="+envVar.Value)
	}

	err = cmd.Start()
	if err != nil {
		stderr.Write([]byte(err.Error()))
	}
	this.proc = cmd.Process
	if this.proc != nil {
		this.Pid = this.proc.Pid
	}
	if this.onStart != nil {
		this.onStart()
	}
	err = cmd.Wait()
	defer func() {
		if this.onStop != nil {
			this.onStop()
		}
	}()
	if err != nil {
		return
	}

	return
}

// Kill进程
func (this *Process) Kill() error {
	if this.proc == nil {
		return nil
	}
	return this.proc.Kill()
}
