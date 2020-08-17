// +build windows

package teaplugins

import (
	"errors"
	"github.com/Microsoft/go-winio"
	"github.com/iwind/TeaGo/logs"
	"net"
	"os"
	"path/filepath"
	"time"
)

// 加载插件
func (this *Loader) Load() error {
	reader, w /** 子进程写入器 **/, err := os.Pipe()
	if err != nil {
		return err
	}

	r /** 子进程读取器 **/, writer, err := os.Pipe()
	if err != nil {
		return err
	}

	rFile := `\\.\pipe\teaweb.reader.` + this.shortFileName() + `.pipe`
	wFile := `\\.\pipe\teaweb.writer.` + this.shortFileName() + `.pipe`

	rListener, err := winio.ListenPipe(rFile, nil)
	if err != nil {
		return errors.New("ERROR1:" + err.Error())
	}
	go func() {
		for {
			conn, err := rListener.Accept()
			if err != nil {
				logs.Error(err)
				break
			}

			go func(conn net.Conn) {
				data := make([]byte, 1024)
				for {
					n, err := conn.Read(data)
					if n > 0 {
						_, _ = w.Write(data[:n])
					}
					if err != nil {
						logs.Println("[plugin]read error:", err)
						break
					}
				}
			}(conn)
		}
	}()

	wListener, err := winio.ListenPipe(wFile, nil)
	if err != nil {
		return errors.New("ERROR2:" + err.Error())
	}
	go func() {
		for {
			conn, err := wListener.Accept()
			if err != nil {
				logs.Error(err)
				break
			}

			go func(conn net.Conn) {
				data := make([]byte, 1024)
				for {
					n, err := r.Read(data)
					if n > 0 {
						_, _ = conn.Write(data[:n])
					}
					if err != nil {
						logs.Println("[plugin]write error:", err)
						break
					}
				}
			}(conn)
		}
	}()

	this.writer = writer

	go this.pipe(reader, writer)

	p, err := this.startProcess(this.path)

	if err != nil {
		logs.Println("[plugin][" + this.shortFileName() + "]start failed:" + err.Error())
		return err
	}

	_, err = p.Wait()

	_ = rListener.Close()
	_ = wListener.Close()
	_ = reader.Close()
	_ = writer.Close()

	// 重新加载
	time.Sleep(1 * time.Second)
	return this.Load()
}

func (this *Loader) startProcess(path string) (*os.Process, error) {
	attrs := &os.ProcAttr{
		Dir:   filepath.Dir(path),
		Env:   os.Environ(),
		Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
	}
	return os.StartProcess(path, []string{}, attrs)
}
