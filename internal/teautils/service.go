package teautils

import (
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/files"
	"github.com/iwind/TeaGo/logs"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

// 服务管理器
type ServiceManager struct {
	Name        string
	Description string

	fp         *os.File
	logger     *log.Logger
	onceLocker sync.Once
}

// 获取对象
func NewServiceManager(name, description string) *ServiceManager {
	manager := &ServiceManager{
		Name:        name,
		Description: description,
	}

	// root
	manager.resetRoot()

	return manager
}

// 设置服务
func (this *ServiceManager) setup() {
	this.onceLocker.Do(func() {
		logFile := files.NewFile(Tea.Root + "/logs/service.log")
		if logFile.Exists() {
			logFile.Delete()
		}

		//logger
		fp, err := os.OpenFile(Tea.Root+"/logs/service.log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
		if err != nil {
			logs.Error(err)
			return
		}
		this.fp = fp
		this.logger = log.New(fp, "", log.LstdFlags)
	})
}

//记录普通日志
func (this *ServiceManager) Log(msg string) {
	this.setup()
	if this.logger == nil {
		return
	}
	this.logger.Println("[info]" + msg)
}

// 记录错误日志
func (this *ServiceManager) LogError(msg string) {
	this.setup()
	if this.logger == nil {
		return
	}
	this.logger.Println("[error]" + msg)
}

// 关闭
func (this *ServiceManager) Close() error {
	if this.fp != nil {
		return this.fp.Close()
	}
	return nil
}

// 重置Root
func (this *ServiceManager) resetRoot() {
	if !Tea.IsTesting() {
		exePath, err := os.Executable()
		if err != nil {
			exePath = os.Args[0]
		}
		link, err := filepath.EvalSymlinks(exePath)
		if err == nil {
			exePath = link
		}
		fullPath, err := filepath.Abs(exePath)
		if err == nil {
			Tea.UpdateRoot(filepath.Dir(filepath.Dir(fullPath)))
		}
	}
	Tea.SetPublicDir(Tea.Root + Tea.DS + "web" + Tea.DS + "public")
	Tea.SetViewsDir(Tea.Root + Tea.DS + "web" + Tea.DS + "views")
	Tea.SetTmpDir(Tea.Root + Tea.DS + "web" + Tea.DS + "tmp")
}

// 保持命令行窗口是打开的
func (this *ServiceManager) PauseWindow() {
	if runtime.GOOS != "windows" {
		return
	}

	b := make([]byte, 1)
	os.Stdin.Read(b)
}
