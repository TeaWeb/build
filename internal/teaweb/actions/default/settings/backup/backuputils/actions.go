package backuputils

import (
	"archive/zip"
	"fmt"
	"github.com/TeaWeb/build/internal/teaconfigs/agents"
	"github.com/TeaWeb/build/internal/teaproxy"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/files"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	timeutil "github.com/iwind/TeaGo/utils/time"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

// 是否需要重启
var shouldRestart = false

// 列出所有文件
func ActionListFiles() []maps.Map {
	// 已备份
	result := []maps.Map{}

	reg := regexp.MustCompile("^\\d{8}\\.zip$")
	dir := files.NewFile(Tea.Root + "/backups/")
	if dir.Exists() {
		for _, f := range dir.List() {
			if !reg.MatchString(f.Name()) {
				continue
			}
			modifiedTime, _ := f.LastModified()
			size, _ := f.Size()
			result = append(result, maps.Map{
				"name":        f.Name(),
				"time":        timeutil.Format("Y-m-d H:i:s", modifiedTime),
				"size":        fmt.Sprintf("%.2f", float64(size)/1024/1024), // M
				"sizeBytes":   size,
				"isToday":     timeutil.Format("Ymd")+".zip" == f.Name(),
				"isYesterday": timeutil.Format("Ymd", time.Now().Add(-24*time.Hour))+".zip" == f.Name(),
			})
		}
	}

	lists.Sort(result, func(i int, j int) bool {
		return result[i].GetString("name") > result[j].GetString("name")
	})

	return result
}

// 下载文件
func ActionDownloadFile(filename string, responseWriter http.ResponseWriter, onNotFound func()) {
	reg := regexp.MustCompile("^\\d{8}\\.zip$")
	if !reg.MatchString(filename) {
		onNotFound()
		return
	}

	fp, err := os.Open(Tea.Root + "/backups/" + filename)
	if err != nil {
		onNotFound()
		return
	}
	defer fp.Close()

	responseWriter.Header().Set("Content-Disposition", "attachment; filename=\""+filename+"\"")
	_, err = io.Copy(responseWriter, fp)
	if err != nil {
		logs.Error(err)
		return
	}
}

// 删除文件
func ActionDeleteFile(filename string, onError func(err error)) (goNext bool) {
	file := files.NewFile(Tea.Root + "/backups/" + filename)
	if file.Exists() {
		err := file.Delete()
		if err != nil {
			onError(err)
			return
		}
	}

	return true
}

// 还原
func ActionRestoreFile(filename string, onFail func(message string)) (goNext bool) {
	file := files.NewFile(Tea.Root + "/backups/" + filename)
	if !file.Exists() {
		onFail("指定的备份文件不存在")
		return
	}

	// 解压
	reader, err := zip.OpenReader(file.Path())
	if err != nil {
		onFail("无法读取：" + err.Error())
		return
	}
	defer reader.Close()

	// 清除backup configs
	tmpDir := files.NewFile(Tea.Root + "/backups/configs")
	if tmpDir.Exists() {
		err := tmpDir.DeleteAll()
		if err != nil {
			onFail("无法清除backups/configs")
			return
		}
	}

	for _, entry := range reader.File {
		dir := filepath.Dir(entry.Name)
		target := files.NewFile(Tea.Root + "/backups/" + dir)
		if !target.Exists() {
			err := target.MkdirAll()
			if err != nil {
				onFail("创建目录失败：" + dir)
				return
			}
		}
		reader, err := entry.Open()
		if err != nil {
			onFail("文件读取失败：" + err.Error())
			return
		}
		data := []byte{}
		for {
			buf := make([]byte, 1024)
			n, err := reader.Read(buf)
			if n > 0 {
				data = append(data, buf[:n]...)
			}
			if err != nil {
				break
			}
		}
		err = files.NewFile(Tea.Root + "/backups/" + entry.Name).Write(data)
		if err != nil {
			reader.Close()
			onFail("文件写入失败：" + err.Error())
			return
		}
		reader.Close()
	}

	// 修改老的配置文件
	oldDir := Tea.Root + "/old.configs." + timeutil.Format("YmdHis")
	err = os.Rename(Tea.ConfigDir(), oldDir)
	if err != nil {
		onFail("原配置清空失败：" + err.Error())
		return
	}

	// 创建新的目录
	err = os.Rename(Tea.Root+"/backups/configs", Tea.ConfigDir())
	if err != nil {
		// 还原
		os.Rename(oldDir, Tea.ConfigDir())

		onFail("新配置拷贝失败")
		return
	}

	teaproxy.SharedManager.Reload()
	agents.NotifyAgentsChange()

	shouldRestart = true

	return true
}

// 判断是否需要重启
func ShouldRestart() bool {
	return shouldRestart
}
