package teaproxy

import "io"

// FTP响应内容
type FTPResponseBody struct {
	io.ReadCloser
	callback func() // 关闭时回调
}

// 关闭
func (this *FTPResponseBody) Close() error {
	err := this.ReadCloser.Close()
	if this.callback != nil {
		this.callback()
	}
	return err
}
