package ssl

import (
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/actions"
	"io/ioutil"
	"path/filepath"
	"strings"
)

type DownloadFileAction actions.Action

// 下载证书和密钥相关文件
func (this *DownloadFileAction) RunGet(params struct {
	ServerId string
	File     string
	View     bool
}) {
	server := teaconfigs.NewServerConfigFromId(params.ServerId)
	if server == nil {
		this.Fail("找不到Server")
	}

	if server.SSL == nil {
		this.Fail("还没有配置SSL")
	}

	server.SSL.Validate()

	file := params.File
	if len(file) == 0 {
		this.Fail("file path should not be empty")
	}

	// 校验文件名，以确保安全性
	if !server.SSL.ContainsFile(file) {
		this.Fail("the file has been forbidden to download")
	}

	fullPath := file
	if !strings.ContainsAny(file, "/\\") {
		fullPath = Tea.ConfigFile(file)
	}
	data, err := ioutil.ReadFile(fullPath)
	if err != nil {
		this.WriteString(err.Error())
		return
	}

	if params.View { // 在线浏览
		this.Write(data)
	} else { // 下载
		this.AddHeader("Content-Disposition", "attachment; filename=\""+filepath.Base(file)+"\";")
		this.Write(data)
	}
}
