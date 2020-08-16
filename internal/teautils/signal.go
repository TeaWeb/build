// +build !windows

package teautils

import (
	"os"
	"os/signal"
)

// 监听Signal
func ListenSignal(f func(sig os.Signal), sig ...os.Signal) {
	ch := make(chan os.Signal, 8)
	signal.Notify(ch, sig...)
	go func() {
		for r := range ch {
			f(r)
		}
	}()
}

// 通知Signal
func NotifySignal(proc *os.Process, sig os.Signal) error {
	return proc.Signal(sig)
}
