package teaproxy

import (
	"io"
	"sync"
)

// Tunnel响应包装，主要是为了覆盖Close()方法
type TunnelResponseBody struct {
	io.ReadCloser
	locker *sync.Mutex
}

// 关闭
func (this *TunnelResponseBody) Close() error {
	err := this.ReadCloser.Close()
	this.locker.Unlock()
	return err
}
