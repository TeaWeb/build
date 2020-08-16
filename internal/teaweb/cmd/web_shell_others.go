// +build !windows

package cmd

import (
	"io"
)

// 启动服务模式
func (this *WebShell) ExecService(writer io.Writer) bool {
	// do nothing beyond windows
	return true
}
