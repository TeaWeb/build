// +build windows

package teautils

import (
	"errors"
	"fmt"
	"github.com/Microsoft/go-winio"
	"github.com/TeaWeb/build/internal/teaconst"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/types"
	"net"
	"os"
	"strconv"
	"syscall"
)

var signalPipe net.Listener

const signalPipePath = `\\.\pipe\` + teaconst.TeaProcessName + `.signal.pipe`

// 监听Signal
func ListenSignal(f func(sig os.Signal), sig ...os.Signal) {
	pipe, err := winio.ListenPipe(signalPipePath+"."+strconv.Itoa(os.Getpid()), nil)
	if err != nil {
		logs.Error(err)
		return
	}
	signalPipe = pipe

	go func() {
		for {
			conn, err := signalPipe.Accept()
			if err != nil {
				logs.Error(err)
				continue
			}

			go func(conn net.Conn) {
				buf := make([]byte, 16)
				for {
					n, err := conn.Read(buf)
					if n > 0 {
						data := buf[:n]
						r := syscall.Signal(types.Int(string(data)))

						// 是否存在
						found := false
						for _, s := range sig {
							if r == s {
								found = true
								break
							}
						}
						if !found {
							logs.Println("[ERROR]undefined signal '" + r.String() + "'")
							continue
						}

						f(r)
					}
					if err != nil {
						break
					}
				}
			}(conn)
		}
	}()
}

// 通知Signal
func NotifySignal(proc *os.Process, sig os.Signal) error {
	conn, err := winio.DialPipe(signalPipePath+"."+strconv.Itoa(proc.Pid), nil)
	if err != nil {
		return errors.New("can not connect to signal pipe: " + err.Error())
	}
	_, err = conn.Write([]byte(fmt.Sprintf("%d", sig)))
	if err != nil {
		return errors.New("signal sending failed: " + err.Error())
	}

	_ = conn.Close()

	return nil
}
