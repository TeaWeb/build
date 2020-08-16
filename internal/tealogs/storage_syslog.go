package tealogs

import (
	"errors"
	"github.com/TeaWeb/build/internal/tealogs/accesslogs"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"os/exec"
	"runtime"
	"strconv"
)

type SyslogStorageProtocol = string

const (
	SyslogStorageProtocolTCP    SyslogStorageProtocol = "tcp"
	SyslogStorageProtocolUDP    SyslogStorageProtocol = "udp"
	SyslogStorageProtocolNone   SyslogStorageProtocol = "none"
	SyslogStorageProtocolSocket SyslogStorageProtocol = "socket"
)

type SyslogStoragePriority = int

const (
	SyslogStoragePriorityEmerg SyslogStoragePriority = iota
	SyslogStoragePriorityAlert
	SyslogStoragePriorityCrit
	SyslogStoragePriorityErr
	SyslogStoragePriorityWarning
	SyslogStoragePriorityNotice
	SyslogStoragePriorityInfo
	SyslogStoragePriorityDebug
)

var SyslogStoragePriorities = []maps.Map{
	{
		"name":  "[无]",
		"value": -1,
	},
	{
		"name":  "EMERG",
		"value": SyslogStoragePriorityEmerg,
	},
	{
		"name":  "ALERT",
		"value": SyslogStoragePriorityAlert,
	},
	{
		"name":  "CRIT",
		"value": SyslogStoragePriorityCrit,
	},
	{
		"name":  "ERR",
		"value": SyslogStoragePriorityErr,
	},
	{
		"name":  "WARNING",
		"value": SyslogStoragePriorityWarning,
	},
	{
		"name":  "NOTICE",
		"value": SyslogStoragePriorityNotice,
	},
	{
		"name":  "INFO",
		"value": SyslogStoragePriorityInfo,
	},
	{
		"name":  "DEBUG",
		"value": SyslogStoragePriorityDebug,
	},
}

// syslog存储策略
type SyslogStorage struct {
	Storage `yaml:", inline"`

	Protocol   string                `yaml:"protocol" json:"protocol"` // SysLogStorageProtocol*
	ServerAddr string                `yaml:"serverAddr" json:"serverAddr"`
	ServerPort int                   `yaml:"serverPort" json:"serverPort"`
	Socket     string                `yaml:"socket" json:"socket"` // sock file
	Tag        string                `yaml:"tag" json:"tag"`
	Priority   SyslogStoragePriority `yaml:"priority" json:"priority"`

	exe string
}

// 开启
func (this *SyslogStorage) Start() error {
	if runtime.GOOS != "linux" {
		return errors.New("'syslog' storage only works on linux")
	}

	exe, err := exec.LookPath("logger")
	if err != nil {
		return err
	}

	this.exe = exe

	return nil
}

// 写入日志
func (this *SyslogStorage) Write(accessLogs []*accesslogs.AccessLog) error {
	if len(accessLogs) == 0 {
		return nil
	}

	args := []string{}
	if len(this.Tag) > 0 {
		args = append(args, "-t", this.Tag)
	}

	if this.Priority >= 0 {
		args = append(args, "-p", strconv.Itoa(this.Priority))
	}

	switch this.Protocol {
	case SyslogStorageProtocolTCP:
		args = append(args, "-T")
		if len(this.ServerAddr) > 0 {
			args = append(args, "-n", this.ServerAddr)
		}
		if this.ServerPort > 0 {
			args = append(args, "-P", strconv.Itoa(this.ServerPort))
		}
	case SyslogStorageProtocolUDP:
		args = append(args, "-d")
		if len(this.ServerAddr) > 0 {
			args = append(args, "-n", this.ServerAddr)
		}
		if this.ServerPort > 0 {
			args = append(args, "-P", strconv.Itoa(this.ServerPort))
		}
	case SyslogStorageProtocolSocket:
		args = append(args, "-u")
		args = append(args, this.Socket)
	case SyslogStorageProtocolNone:
		// do nothing
	}

	args = append(args, "-S", "10240")

	cmd := exec.Command(this.exe, args...)
	w, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	err = cmd.Start()
	if err != nil {
		return err
	}

	for _, accessLog := range accessLogs {
		data, err := this.FormatAccessLogBytes(accessLog)
		if err != nil {
			logs.Error(err)
			continue
		}
		_, err = w.Write(data)
		if err != nil {
			logs.Error(err)
		}

		_, err = w.Write([]byte("\n"))
		if err != nil {
			logs.Error(err)
		}
	}

	_ = w.Close()
	err = cmd.Wait()
	if err != nil {
		return err
	}

	return nil
}

// 关闭
func (this *SyslogStorage) Close() error {
	return nil
}
