package teastats

import (
	"github.com/TeaWeb/build/internal/tealogs/accesslogs"
)

// 服务队列配置
// 为每一个代理服务配置一个服务队列
type ServerQueue struct {
	Queue   *Queue
	Filters map[string]FilterInterface // code => instance
}

// 停止队列
func (this *ServerQueue) Stop() {
	for _, f := range this.Filters {
		f.Stop()
	}

	this.Queue.Stop()
	this.Queue = nil
	this.Filters = nil
}

// 启动一个Filter
func (this *ServerQueue) StartFilter(code string) {
	_, found := this.Filters[code]
	if found {
		return
	}

	instance := FindNewFilter(code)
	if instance == nil {
		return
	}

	this.Filters[code] = instance
	instance.Start(this.Queue, code)
}

// 筛选
func (this *ServerQueue) Filter(accessLog *accesslogs.AccessLog) {
	for _, f := range this.Filters {
		f.Filter(accessLog)
	}
}

// 提交数据
func (this *ServerQueue) Commit() {
	for _, f := range this.Filters {
		f.Commit()
	}
}
