package teautils

import (
	"bytes"
	"errors"
	"os/exec"
)

// 命令执行器
type CommandExecutor struct {
	commands []*Command
}

// 获取新对象
func NewCommandExecutor() *CommandExecutor {
	return &CommandExecutor{}
}

// 添加命令
func (this *CommandExecutor) Add(command string, arg ...string) {
	this.commands = append(this.commands, &Command{
		Name: command,
		Args: arg,
	})
}

// 执行命令
func (this *CommandExecutor) Run() (output string, err error) {
	if len(this.commands) == 0 {
		return "", errors.New("no commands no run")
	}
	var lastCmd *exec.Cmd = nil
	var lastData []byte = nil
	for _, command := range this.commands {
		cmd := exec.Command(command.Name, command.Args...)
		stdout := bytes.NewBuffer([]byte{})
		cmd.Stdout = stdout
		if lastCmd != nil {
			cmd.Stdin = bytes.NewBuffer(lastData)
		}
		err = cmd.Start()
		if err != nil {
			return "", err
		}

		err = cmd.Wait()
		if err != nil {
			_, ok := err.(*exec.ExitError)
			if ok {
				return "", nil
			}

			return "", err
		}
		lastData = stdout.Bytes()

		lastCmd = cmd
	}

	return string(bytes.TrimSpace(lastData)), nil
}
