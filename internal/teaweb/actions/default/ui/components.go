package ui

import (
	"bytes"
	"github.com/TeaWeb/build/internal/teautils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/files"
	"github.com/iwind/TeaGo/logs"
)

type ComponentsAction actions.Action

func (this *ComponentsAction) RunGet(params struct{}) {
	this.AddHeader("Content-Type", "text/javascript; charset=utf-8")

	// TODO 增加缓存

	webRoot := teautils.WebRoot() + "/public/js/components/"
	f := files.NewFile(webRoot)

	buf := bytes.NewBuffer([]byte{})
	f.Range(func(file *files.File) {
		if !file.IsFile() {
			return
		}
		if file.Ext() != ".js" {
			return
		}
		data, err := file.ReadAll()
		if err != nil {
			logs.Error(err)
			return
		}
		buf.Write(data)
		buf.WriteByte('\n')
		buf.WriteByte('\n')
	})
	this.Write(buf.Bytes())
}
