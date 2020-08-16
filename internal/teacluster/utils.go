package teacluster

import (
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/files"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/logs"
	"path/filepath"
	"strings"
)

func RangeFiles(f func(file *files.File, relativePath string)) {
	configAbs, err := filepath.Abs(Tea.ConfigDir())
	if err != nil {
		logs.Error(err)
		return
	}
	files.NewFile(configAbs).Range(func(file *files.File) {
		// *.conf & ssl.*
		if strings.HasSuffix(file.Name(), ".sum") ||
			strings.HasSuffix(file.Name(), ".script") ||
			strings.HasSuffix(file.Name(), ".bat") ||
			strings.HasPrefix(file.Name(), ".") {
			return
		}
		if lists.ContainsString([]string{
			"node.conf",
			"server.conf",
			"agent.local.conf",
			"board.local.conf",
		}, file.Name()) {
			return
		}
		absPath, _ := file.AbsPath()
		relativePath := strings.TrimLeft(strings.TrimPrefix(absPath, configAbs), "/\\")
		f(file, relativePath)
	})
}
