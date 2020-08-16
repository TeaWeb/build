package notices

import (
	"crypto/tls"
	"errors"
	"github.com/TeaWeb/build/internal/teaconst"
	"net"
	"net/mail"
	"net/smtp"
	"time"
)

// 邮件媒介
type NoticeEmailMedia struct {
	SMTP     string `yaml:"smtp" json:"smtp"`
	Username string `yaml:"username" json:"username"`
	Password string `yaml:"password" json:"password"`
	From     string `yaml:"from" json:"from"`
}

// 获取新对象
func NewNoticeEmailMedia() *NoticeEmailMedia {
	return &NoticeEmailMedia{}
}

func (this *NoticeEmailMedia) Send(user string, subject string, body string) (resp []byte, err error) {
	if len(this.SMTP) == 0 {
		return nil, errors.New("host address should be specified")
	}

	// 自动加端口

	if _, _, err := net.SplitHostPort(this.SMTP); err != nil {
		this.SMTP += ":587"
	}

	if len(this.From) == 0 {
		this.From = this.Username
	}

	contentType := "Content-Type: text/html; charset=UTF-8"
	msg := []byte("To: " + user + "\r\nFrom: \"" + teaconst.TeaProductName + "\" <" + this.From + ">\r\nSubject: " + subject + "\r\n" + contentType + "\r\n\r\n" + body)

	return nil, this.SendMail(this.From, []string{user}, msg)
}

// 是否需要用户标识
func (this *NoticeEmailMedia) RequireUser() bool {
	return true
}

func (this *NoticeEmailMedia) SendMail(from string, to []string, message []byte) error {
	serverName := this.SMTP
	username := this.Username
	password := this.Password

	_, err := mail.ParseAddress(from)
	if err != nil {
		return err
	}

	if len(to) == 0 {
		return errors.New("recipients should not be empty")
	}

	for _, to1 := range to {
		_, err := mail.ParseAddress(to1)
		if err != nil {
			return err
		}
	}

	host, port, _ := net.SplitHostPort(serverName)

	var client *smtp.Client

	// TLS config
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}

	// 587 port: prefer START_TLS
	if port == "587" {
		conn, err := net.DialTimeout("tcp", serverName, 10*time.Second)
		if err != nil {
			conn, err := tls.Dial("tcp", serverName, tlsConfig)
			if err != nil {
				return err
			}
			client, err = smtp.NewClient(conn, host)
			if err != nil {
				return err
			}
		} else {
			client, err = smtp.NewClient(conn, host)
			if err != nil {
				return err
			}
			client.StartTLS(tlsConfig)
		}
	} else {
		conn, err := tls.Dial("tcp", serverName, tlsConfig)
		if err != nil {
			conn, err := net.DialTimeout("tcp", serverName, 10*time.Second)
			if err != nil {
				return err
			}

			client, err = smtp.NewClient(conn, host)
			if err != nil {
				return err
			}
			client.StartTLS(tlsConfig)
		} else {
			client, err = smtp.NewClient(conn, host)
			if err != nil {
				return err
			}
		}
	}

	// Auth
	auth := smtp.PlainAuth("", username, password, host)
	if err := client.Auth(auth); err != nil {
		return err
	}

	// To && From
	if err := client.Mail(from); err != nil {
		return err
	}

	for _, to1 := range to {
		if err := client.Rcpt(to1); err != nil {
			return err
		}
	}

	// Data
	w, err := client.Data()
	if err != nil {
		return err
	}

	_, err = w.Write(message)
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	client.Quit()
	client.Close()

	return nil
}
