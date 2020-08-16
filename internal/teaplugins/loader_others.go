// +build !windows

package teaplugins

import (
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/processes"
	"os"
	"time"
)

func (this *Loader) Load() error {
	reader, w /** 子进程写入器 **/, err := os.Pipe()
	if err != nil {
		return err
	}

	r /** 子进程读取器 **/, writer, err := os.Pipe()
	if err != nil {
		return err
	}

	this.writer = writer

	go this.pipe(reader, writer)

	p := processes.NewProcess(this.path)
	p.AppendFile(r, w)

	err = p.Start()
	if err != nil {
		logs.Println("[plugin][" + this.shortFileName() + "]start failed:" + err.Error())
		return err
	}

	err = p.Wait()
	if err != nil {
		logs.Println("[plugin][" + this.shortFileName() + "]wait failed" + err.Error())

		reader.Close()

		// 重新加载
		time.Sleep(1 * time.Second)
		return this.Load()
	}

	return nil
}
