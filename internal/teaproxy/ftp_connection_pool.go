package teaproxy

import (
	"errors"
	"github.com/jlaffaye/ftp"
	"runtime"
	"strings"
	"sync/atomic"
	"time"
)

var ErrFTPTooManyConnections = errors.New("ftp: too many connections")

// FTP连接池
type FTPConnectionPool struct {
	addr     string
	username string
	password string
	timeout  time.Duration
	dir      string

	currentConnections int64
	maxConnections     int64

	c chan *ftp.ServerConn
}

// 获取新的连接
func (this *FTPConnectionPool) Get() (*ftp.ServerConn, error) {
	if this.timeout <= 0 {
		this.timeout = 5 * time.Second
	}
	if this.maxConnections <= 0 {
		this.maxConnections = int64(runtime.NumCPU())
	}

	select {
	case client := <-this.c:
		return client, nil
	default:
		if atomic.LoadInt64(&this.currentConnections) >= this.maxConnections {
			return nil, ErrFTPTooManyConnections
		}

		this.Increase()

		// create a new connection
		client, err := ftp.DialTimeout(this.addr, this.timeout)
		if err != nil {
			this.Decrease()
			return nil, err
		}
		if len(this.username) > 0 {
			err = client.Login(this.username, this.password)
			if err != nil {
				this.Decrease()
				_ = client.Quit()
				return nil, err
			}
		}
		if len(this.dir) > 0 {
			err = client.ChangeDir(strings.TrimLeft(this.dir, "/"))
			if err != nil {
				this.Decrease()
				_ = client.Quit()
				return nil, err
			}
		}

		return client, err
	}
}

// 复用连接
func (this *FTPConnectionPool) Put(client *ftp.ServerConn) {
	select {
	case this.c <- client:
	default:
		this.Decrease()
		_ = client.Quit()
	}
}

// 关闭连接
func (this *FTPConnectionPool) Close(client *ftp.ServerConn) error {
	this.Decrease()
	err := client.Quit()
	return err
}

// 关闭所有连接
func (this *FTPConnectionPool) CloseAll() error {
FOR:
	for {
		select {
		case conn := <-this.c:
			this.Decrease()
			_ = conn.Quit()
		default:
			break FOR
		}
	}
	return nil
}

// 增加连接数量
func (this *FTPConnectionPool) Increase() {
	atomic.AddInt64(&this.currentConnections, 1)
}

// 减少连接数量
func (this *FTPConnectionPool) Decrease() {
	atomic.AddInt64(&this.currentConnections, -1)
}
