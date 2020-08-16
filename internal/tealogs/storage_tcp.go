package tealogs

import (
	"errors"
	"github.com/TeaWeb/build/internal/tealogs/accesslogs"
	"github.com/iwind/TeaGo/logs"
	"net"
	"sync"
)

// TCP存储策略
type TCPStorage struct {
	Storage `yaml:", inline"`

	Network string `yaml:"network" json:"network"` // tcp, unix
	Addr    string `yaml:"addr" json:"addr"`

	writeLocker sync.Mutex

	connLocker sync.Mutex
	conn       net.Conn
}

// 开启
func (this *TCPStorage) Start() error {
	if len(this.Network) == 0 {
		return errors.New("'network' should not be nil")
	}
	if len(this.Addr) == 0 {
		return errors.New("'addr' should not be nil")
	}
	return nil
}

// 写入日志
func (this *TCPStorage) Write(accessLogs []*accesslogs.AccessLog) error {
	if len(accessLogs) == 0 {
		return nil
	}

	err := this.connect()
	if err != nil {
		return err
	}

	conn := this.conn
	if conn == nil {
		return errors.New("connection should not be nil")
	}

	this.writeLocker.Lock()
	defer this.writeLocker.Unlock()

	for _, accessLog := range accessLogs {
		data, err := this.FormatAccessLogBytes(accessLog)
		if err != nil {
			logs.Error(err)
			continue
		}
		_, err = conn.Write(data)
		if err != nil {
			_ = this.Close()
			break
		}
		_, err = conn.Write([]byte("\n"))
		if err != nil {
			_ = this.Close()
			break
		}
	}

	return nil
}

// 关闭
func (this *TCPStorage) Close() error {
	this.connLocker.Lock()
	defer this.connLocker.Unlock()

	if this.conn != nil {
		err := this.conn.Close()
		this.conn = nil
		return err
	}
	return nil
}

func (this *TCPStorage) connect() error {
	this.connLocker.Lock()
	defer this.connLocker.Unlock()

	if this.conn != nil {
		return nil
	}

	conn, err := net.Dial(this.Network, this.Addr)
	if err != nil {
		return err
	}
	this.conn = conn

	return nil
}
