package tealogs

import (
	"bytes"
	"errors"
	"github.com/TeaWeb/build/internal/tealogs/accesslogs"
	"github.com/iwind/TeaGo/logs"
	"os/exec"
	"sync"
)

// 通过命令行存储
type CommandStorage struct {
	Storage `yaml:", inline"`

	Command string   `yaml:"command" json:"command"`
	Args    []string `yaml:"args" json:"args"`
	Dir     string   `yaml:"dir" json:"dir"`

	writeLocker sync.Mutex
}

// 启动
func (this *CommandStorage) Start() error {
	if len(this.Command) == 0 {
		return errors.New("'command' should not be empty")
	}
	return nil
}

// 写入日志
func (this *CommandStorage) Write(accessLogs []*accesslogs.AccessLog) error {
	if len(accessLogs) == 0 {
		return nil
	}

	this.writeLocker.Lock()
	defer this.writeLocker.Unlock()

	cmd := exec.Command(this.Command, this.Args...)
	if len(this.Dir) > 0 {
		cmd.Dir = this.Dir
	}

	stdout := bytes.NewBuffer([]byte{})
	cmd.Stdout = stdout

	w, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	err = cmd.Start()
	if err != nil {
		return err
	}
	for _, accessLog := range accessLogs {
		data, err := this.FormatAccessLogBytes(accessLog)
		if err != nil {
			logs.Error(err)
			continue
		}
		_, err = w.Write(data)
		if err != nil {
			logs.Error(err)
		}

		_, err = w.Write([]byte("\n"))
		if err != nil {
			logs.Error(err)
		}
	}
	_ = w.Close()
	err = cmd.Wait()
	if err != nil {
		logs.Error(err)

		if stdout.Len() > 0 {
			logs.Error(errors.New(string(stdout.Bytes())))
		}
	}

	return nil
}

// 关闭
func (this *CommandStorage) Close() error {
	return nil
}
