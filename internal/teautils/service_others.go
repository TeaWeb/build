// +build !linux,!windows

package teautils

// 安装服务
func (this *ServiceManager) Install(exePath string, args []string) error {
	return nil
}

// 启动服务
func (this *ServiceManager) Start() error {
	return nil
}

// 删除服务
func (this *ServiceManager) Uninstall() error {
	return nil
}
