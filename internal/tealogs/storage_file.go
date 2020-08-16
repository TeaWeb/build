package tealogs

import (
	"errors"
	"github.com/TeaWeb/build/internal/tealogs/accesslogs"
	"github.com/iwind/TeaGo/logs"
	"os"
	"path/filepath"
	"sync"
)

// 文件存储策略
type FileStorage struct {
	Storage `yaml:", inline"`

	Path       string `yaml:"path" json:"path"`             // 文件路径，支持变量：${year|month|week|day|hour|minute|second}
	AutoCreate bool   `yaml:"autoCreate" json:"autoCreate"` // 是否自动创建目录

	writeLocker sync.Mutex

	files       map[string]*os.File // path => *File
	filesLocker sync.Mutex
}

// 开启
func (this *FileStorage) Start() error {
	if len(this.Path) == 0 {
		return errors.New("'path' should not be empty")
	}

	this.files = map[string]*os.File{}

	return nil
}

// 写入日志
func (this *FileStorage) Write(accessLogs []*accesslogs.AccessLog) error {
	if len(accessLogs) == 0 {
		return nil
	}

	fp := this.fp()
	if fp == nil {
		return errors.New("file pointer should not be nil")
	}
	this.writeLocker.Lock()
	defer this.writeLocker.Unlock()

	for _, accessLog := range accessLogs {
		data, err := this.FormatAccessLogBytes(accessLog)
		if err != nil {
			logs.Error(err)
			continue
		}
		_, err = fp.Write(data)
		if err != nil {
			_ = this.Close()
			break
		}
		_, _ = fp.WriteString("\n")
	}
	return nil
}

// 关闭
func (this *FileStorage) Close() error {
	this.filesLocker.Lock()
	defer this.filesLocker.Unlock()

	var resultErr error
	for _, f := range this.files {
		err := f.Close()
		if err != nil {
			resultErr = err
		}
	}
	return resultErr
}

func (this *FileStorage) fp() *os.File {
	path := this.FormatVariables(this.Path)

	this.filesLocker.Lock()
	defer this.filesLocker.Unlock()
	fp, ok := this.files[path]
	if ok {
		return fp
	}

	// 关闭其他的文件
	for _, f := range this.files {
		_ = f.Close()
	}

	// 是否创建文件目录
	if this.AutoCreate {
		dir := filepath.Dir(path)
		_, err := os.Stat(dir)
		if os.IsNotExist(err) {
			err = os.MkdirAll(dir, 0777)
			if err != nil {
				logs.Error(err)
				return nil
			}
		}
	}

	// 打开新文件
	fp, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logs.Error(err)
		return nil
	}
	this.files[path] = fp

	return fp
}
