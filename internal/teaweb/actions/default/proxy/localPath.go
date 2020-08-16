package proxy

import (
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/files"
	"runtime"
	"strings"
)

type LocalPathAction actions.Action

// 本地路径查询
func (this *LocalPathAction) Run(params struct {
	Prefix string
}) {
	prefix := params.Prefix
	this.Data["paths"] = []string{}

	if len(prefix) == 0 {
		this.Success()
	}

	dirPrefix := prefix
	sub := ""
	index := 0
	if runtime.GOOS == "windows" {
		prefix = strings.Replace(prefix, "/", "\\", -1)
		index = strings.LastIndex(prefix, "\\")
	} else {
		prefix = strings.Replace(prefix, "\\", "/", -1)
		index = strings.LastIndex(prefix, "/")
	}
	if index > -1 {
		dirPrefix = prefix[:index+1]
		sub = prefix[index+1:]
	}

	this.Data["dir"] = dirPrefix
	this.Data["sub"] = sub

	// 查找路径
	dir := files.NewFile(dirPrefix)
	if !dir.Exists() {
		this.Success()
	}

	paths := []string{}
	for _, f := range dir.List() {
		absPath, _ := f.AbsPath()
		if len(sub) == 0 {
			paths = append(paths, absPath)
		} else if strings.Index(strings.ToLower(f.Name()), strings.ToLower(sub)) > -1 {
			paths = append(paths, absPath)
		}

		if len(paths) >= 10 {
			break
		}
	}
	this.Data["paths"] = paths

	this.Success()
}
