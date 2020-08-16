package tealogs

import (
	"github.com/TeaWeb/build/internal/tealogs/accesslogs"
	"github.com/TeaWeb/build/internal/teautils/logbuffer"
	"github.com/iwind/TeaGo/Tea"
	"runtime"
	"strconv"
)

var (
	accessLogger *AccessLogger = nil
)

// 访问日志记录器
type AccessLogger struct {
	queue chan *accesslogs.AccessLog
}

// 获取新日志对象
func NewAccessLogger() *AccessLogger {
	logger := &AccessLogger{
		queue: make(chan *accesslogs.AccessLog, 10*10000),
	}

	go logger.wait()
	return logger
}

// 获取共享的对象
func SharedLogger() *AccessLogger {
	return accessLogger
}

// 推送日志
func (this *AccessLogger) Push(log *accesslogs.AccessLog) {
	if this.queue == nil {
		return
	}
	this.queue <- log
}

// 等待日志到来
func (this *AccessLogger) wait() {
	// 启动queue
	for i := 0; i < runtime.NumCPU(); i++ {
		(func(index int) {
			buf := logbuffer.NewBuffer(Tea.LogDir() + "/accesslog." + strconv.Itoa(index))
			queue := NewAccessLogQueue(buf, index)
			go queue.Receive(this.queue)
			go queue.Dump()
		})(i)
	}
}
