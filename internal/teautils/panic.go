package teautils

import (
	"github.com/iwind/TeaGo/logs"
	"runtime/debug"
)

// 记录panic日志
func Recover() {
	p := recover()
	if p != nil {
		logs.Println("panic:", p)
		logs.Println(string(debug.Stack()))
	}
}
