package teautils

import (
	"github.com/iwind/TeaGo/Tea"
	"path/filepath"
)

// 临时文件
func TmpFile(path string) string {
	return filepath.Clean(WebRoot() + Tea.DS + "tmp" + Tea.DS + path)
}
