package teautils

import (
	"github.com/iwind/TeaGo/Tea"
	"path/filepath"
)

// Web根目录
func WebRoot() string {
	if Tea.IsTesting() {
		return filepath.Dir(Tea.Root) + Tea.DS + "web"
	} else {
		return Tea.Root + Tea.DS + "web"
	}
}
