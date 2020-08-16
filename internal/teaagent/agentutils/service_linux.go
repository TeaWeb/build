// +build linux

package agentutils

import (
	"errors"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/files"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
)

var systemdServiceFile = "/etc/systemd/system/teaweb-agent.service"
var initServiceFile = "/etc/init.d/teaweb-agent"

// 安装服务
func (this *ServiceManager) Install(exePath string, args []string) error {
	if os.Getgid() != 0 {
		return errors.New("only root users can install the service")
	}

	systemd, err := exec.LookPath("systemctl")
	if err != nil {
		return this.installInitService(exePath, args)
	}

	return this.installSystemdService(systemd, exePath, args)
}

// 启动服务
func (this *ServiceManager) Start() error {
	if os.Getgid() != 0 {
		return errors.New("only root users can start the service")
	}

	if files.NewFile(systemdServiceFile).Exists() {
		systemd, err := exec.LookPath("systemctl")
		if err != nil {
			return err
		}

		return exec.Command(systemd, "start", "teaweb-agent.service").Start()
	}
	return exec.Command("service", "teaweb-agent", "start").Start()
}

// 删除服务
func (this *ServiceManager) Uninstall() error {
	if os.Getgid() != 0 {
		return errors.New("only root users can uninstall the service")
	}

	if files.NewFile(systemdServiceFile).Exists() {
		systemd, err := exec.LookPath("systemctl")
		if err != nil {
			return err
		}

		// disable service
		exec.Command(systemd, "disable", "teaweb-agent.service").Start()

		// reload
		exec.Command(systemd, "daemon-reload")

		return files.NewFile(systemdServiceFile).Delete()
	}

	f := files.NewFile(initServiceFile)
	if f.Exists() {
		return f.Delete()
	}
	return nil
}

// install init service
func (this *ServiceManager) installInitService(exePath string, args []string) error {
	scriptFile := Tea.Root + "/scripts/teaweb-agent"
	if !files.NewFile(scriptFile).Exists() {
		return errors.New("'scripts/teaweb-agent' file not exists")
	}

	data, err := ioutil.ReadFile(scriptFile)
	if err != nil {
		return err
	}

	data = regexp.MustCompile("INSTALL_DIR=.+").ReplaceAll(data, []byte("INSTALL_DIR="+Tea.Root))
	err = ioutil.WriteFile(initServiceFile, data, 0777)
	if err != nil {
		return err
	}

	chkCmd, err := exec.LookPath("chkconfig")
	if err != nil {
		return err
	}

	err = exec.Command(chkCmd, "--add", "teaweb-agent").Start()
	if err != nil {
		return err
	}

	return nil
}

// install systemd service
func (this *ServiceManager) installSystemdService(systemd, exePath string, args []string) error {
	desc := `### BEGIN INIT INFO
# Provides:          teaweb-agent
# Required-Start:    $all
# Required-Stop:
# Default-Start:     2 3 4 5
# Default-Stop:
# Short-Description: TeaWeb Agent Service
### END INIT INFO

[Unit]
Description=TeaWeb Agent Service
Before=shutdown.target

[Service]
Type=forking
ExecStart=` + exePath + ` start
ExecStop=` + exePath + ` stop
ExecReload=` + exePath + ` reload

[Install]
WantedBy=multi-user.target`

	// write file
	err := ioutil.WriteFile(systemdServiceFile, []byte(desc), 0777)
	if err != nil {
		return err
	}

	// stop current systemd service if running
	exec.Command(systemd, "stop", "teaweb-agent.service")

	// reload
	exec.Command(systemd, "daemon-reload")

	// enable
	cmd := exec.Command(systemd, "enable", "teaweb-agent.service")
	return cmd.Run()
}

// 运行
func (this *ServiceManager) Run() {

}
