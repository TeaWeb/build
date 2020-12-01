package agentutils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/TeaWeb/build/internal/teaconfigs/agents"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/files"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/rands"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io"
	"net"
	"os"
	"strings"
	"time"
)

// 安装器
type SSHAuthType = string

const (
	SSHAuthTypePassword = "password"
	SSHAuthTypeKey      = "key"
)

type Installer struct {
	Master       string
	Dir          string
	Host         string
	Port         int
	AuthUsername string
	AuthType     SSHAuthType
	AuthPassword string
	AuthKey      []byte
	Timeout      time.Duration
	GroupId      string

	HostName string
	HostIP   string
	OS       string
	Arch     string

	Logs        []string
	IsInstalled bool
}

// 获取新对象
func NewInstaller() *Installer {
	return &Installer{
		AuthType: SSHAuthTypePassword,
	}
}

// 安装Agent
func (this *Installer) Start() error {
	this.log("start")

	if len(this.Master) == 0 {
		return errors.New("'master' should not be empty")
	}

	if len(this.Dir) == 0 {
		return errors.New("'dir' should not be empty")
	}

	var hostKeyCallback ssh.HostKeyCallback = nil

	// 不使用known_hosts

	if hostKeyCallback == nil {
		hostKeyCallback = func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		}
	}

	methods := []ssh.AuthMethod{}
	if this.AuthType == SSHAuthTypePassword {
		{
			authMethod := ssh.Password(this.AuthPassword)
			methods = append(methods, authMethod)
		}

		{
			authMethod := ssh.KeyboardInteractive(func(user, instruction string, questions []string, echos []bool) (answers []string, err error) {
				if len(questions) == 0 {
					return []string{}, nil
				}
				return []string{this.AuthPassword}, nil
			})
			methods = append(methods, authMethod)
		}
	} else {
		{
			signer, err := ssh.ParsePrivateKey(this.AuthKey)
			if err != nil {
				return err
			}
			authMethod := ssh.PublicKeys(signer)
			methods = append(methods, authMethod)
		}
	}

	config := &ssh.ClientConfig{
		User:            this.AuthUsername,
		Auth:            methods,
		HostKeyCallback: hostKeyCallback,
		Timeout:         this.Timeout,
	}

	this.log("connecting")
	client, err := ssh.Dial("tcp", this.Host+":"+fmt.Sprintf("%d", this.Port), config)
	if err != nil {
		return err
	}
	defer func() {
		_ = client.Close()
	}()

	// hostname
	this.log("get hostname")
	hostName, _, err := this.runCmdOnSSH(client, "hostname")
	if err != nil {
		return err
	}
	this.HostName = string(bytes.TrimSpace(hostName))

	// os
	this.log("get os and arch")
	uname, _, err := this.runCmdOnSSH(client, "uname -a")
	if err != nil {
		return err
	}
	if strings.Index(string(uname), "Darwin") >= 0 {
		this.OS = "darwin"
	} else if strings.Index(string(uname), "Linux") >= 0 {
		this.OS = "linux"
	} else {
		return errors.New("installer only supports darwin and linux")
	}

	if strings.Index(string(uname), "x86_64") > 0 {
		this.Arch = "amd64"
	} else {
		this.Arch = "386"
	}

	// upload installer
	this.log("finding installer file")
	filename := "agentinstaller_" + this.OS + "_" + this.Arch
	installerFile := files.NewFile(Tea.Root + "/web/installers/" + filename)
	if !installerFile.Exists() {
		return errors.New("installer file '" + filename + "' not found")
	}

	this.log("sftp connecting")
	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		return err
	}
	defer func() {
		_ = sftpClient.Close()
	}()

	this.log("create installer file on /tmp")
	writer, err := sftpClient.Create("/tmp/agentinstaller")
	if err != nil {
		return err
	}
	isInstallerWriterClosed := false
	defer func() {
		if !isInstallerWriterClosed {
			_ = writer.Close()
		}

		// 删除
		_, _, _ = this.runCmdOnSSH(client, "unlink /tmp/agentinstaller")
	}()

	this.log("open installer file")
	reader, err := os.OpenFile(installerFile.Path(), os.O_RDONLY, 0666)
	if err != nil {
		return err
	}
	defer func() {
		_ = reader.Close()
	}()

	this.log("copy installer file to host")
	n, err := io.Copy(writer, reader)
	if err != nil {
		return err
	}

	if n == 0 {
		return errors.New("copy installer failed")
	}
	isInstallerWriterClosed = true
	_ = writer.Close() // 明确close一次，以便于下面的chmod和运行

	// chmod
	this.log("chmod")
	_, _, err = this.runCmdOnSSH(client, "chmod 777 /tmp/agentinstaller")
	if err != nil {
		return err
	}

	// run
	this.log("installing")

	agentList, err := agents.SharedAgentList()
	if err != nil {
		return err
	}

	// 创建主机信息
	agent := agents.NewAgentConfig()
	this.log("create new agent " + agent.Id)
	agent.Name = string(bytes.TrimSpace(hostName))
	agent.AutoUpdates = true
	agent.CheckDisconnections = true
	agent.AllowAll = true
	agent.On = true
	agent.Key = rands.HexString(32)
	agent.AddDefaultApps()
	if len(this.GroupId) > 0 {
		agent.AddGroup(this.GroupId)
	}
	err = agent.Save()
	if err != nil {
		return err
	}

	newAgentCreated := false
	defer func() {
		if !newAgentCreated {
			this.log("delete agent " + agent.Id)
			err = agent.Delete()
			if err != nil {
				logs.Error(err)
			}
		} else {
			// 保存到列表
			this.log("add agent to list")
			agentList.AddAgent(agent.Filename())
			err = agentList.Save()
			if err != nil {
				logs.Error(err)
			}
		}
	}()

	output, stderr, err := this.runCmdOnSSH(client, "/tmp/agentinstaller -dir=\""+this.Dir+"\" -master=\""+this.Master+"\" -id=\""+agent.Id+"\" -key=\""+agent.Key+"\"")
	if err != nil {
		return errors.New(err.Error() + ":" + string(stderr))
	}

	outputString := strings.TrimSpace(string(output))
	if len(outputString) == 0 {
		return errors.New("start failed: no response:" + string(stderr))
	}

	m := maps.Map{}
	err = json.Unmarshal([]byte(outputString), &m)
	if err != nil {
		return err
	}

	errString := m.GetString("err")
	isInstalled := m.GetBool("isInstalled")
	ip := m.GetString("ip")
	this.HostIP = ip
	if isInstalled {
		this.IsInstalled = true
		if len(errString) > 0 {
			return errors.New(errString)
		} else {
			newAgentCreated = true

			// 保存IP
			agent.Host = ip
			err := agent.Save()
			if err != nil {
				logs.Error(err)
			}

			this.log("finished")

			// 试着安装启动脚本，不提示错误信息
			if this.AuthUsername == "root" {
				this.installService(sftpClient, client)
			}

			return nil
		}
	}

	return errors.New("error response:" + errString)
}

// 通过SSH运行一个命令
func (this *Installer) runCmdOnSSH(client *ssh.Client, cmd string) (stdoutBytes []byte, stderrBytes []byte, err error) {
	session, err := client.NewSession()
	if err != nil {
		return nil, nil, err
	}
	defer func() {
		_ = session.Close()
	}()

	stdout := bytes.NewBuffer([]byte{})
	stderr := bytes.NewBuffer([]byte{})
	session.Stdout = stdout
	session.Stderr = stderr
	err = session.Run(cmd)
	if err != nil {
		return stdout.Bytes(), stderr.Bytes(), err
	}
	return stdout.Bytes(), stderr.Bytes(), nil
}

// 记录日志
func (this *Installer) log(message string) {
	this.Logs = append(this.Logs, message)
}

// 安装服务脚本
func (this *Installer) installService(sftpClient *sftp.Client, client *ssh.Client) {
	file, err := sftpClient.Create("/etc/init.d/teaweb-agent")
	if err != nil {
		return
	}
	defer func() {
		_ = file.Close()
	}()
	_, err = file.Write([]byte(`#! /bin/bash
#
# teaweb       TeaWeb agent management
#
# chkconfig: 2345 40 90
# description: TeaWeb agent management

# teaweb agent install dir
INSTALL_DIR="` + this.Dir + `/agent"

export TEAROOT=${INSTALL_DIR}

case "$1" in
start)
    ${INSTALL_DIR}/bin/teaweb-agent start
    ;;
stop)
    ${INSTALL_DIR}/bin/teaweb-agent stop
    ;;
restart)
    ${INSTALL_DIR}/bin/teaweb-agent restart
    ;;
*)
    echo $"Usage: $0 {start|stop|restart}"
    exit 2
esac`))
	if err != nil {
		return
	}

	_, _, _ = this.runCmdOnSSH(client, "chmod u+x /etc/init.d/teaweb-agent")
	_, _, _ = this.runCmdOnSSH(client, "chkconfig --add teaweb-agent")
	_, _, _ = this.runCmdOnSSH(client, "systemctl daemon-reload")
}
