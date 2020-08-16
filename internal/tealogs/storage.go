package tealogs

import (
	"fmt"
	"github.com/TeaWeb/build/internal/tealogs/accesslogs"
	"github.com/TeaWeb/build/internal/teautils"
	"github.com/mailru/easyjson"
	"strconv"
	"time"
)

// 日志存储接口
type StorageInterface interface {
	// 开启
	Start() error

	// 写入日志
	Write(accessLogs []*accesslogs.AccessLog) error

	// 关闭
	Close() error
}

// 基础的Storage
type Storage struct {
	Format   StorageFormat `yaml:"format" json:"format"`     // 日志格式
	Template string        `yaml:"template" json:"template"` // 只有在Format为template时有效
}

// 格式化访问日志成字节
func (this *Storage) FormatAccessLogBytes(accessLog *accesslogs.AccessLog) ([]byte, error) {
	if this.Format == StorageFormatTemplate {
		return []byte(accessLog.Format(this.Template)), nil
	} else if this.Format == StorageFormatJSON {
		return easyjson.Marshal(accessLog)
	}

	return easyjson.Marshal(accessLog)
}

// 格式化访问日志成字符串
func (this *Storage) FormatAccessLogString(accessLog *accesslogs.AccessLog) (string, error) {
	if this.Format == StorageFormatTemplate {
		return accessLog.Format(this.Template), nil
	} else if this.Format == StorageFormatJSON {
		data, err := easyjson.Marshal(accessLog)
		if err != nil {
			return "", err
		}

		return string(data), err
	}

	// 默认
	data, err := easyjson.Marshal(accessLog)
	if err != nil {
		return "", err
	}

	return string(data), err
}

// 格式化字符串中的变量
func (this *Storage) FormatVariables(s string) string {
	now := time.Now()
	return teautils.ParseVariables(s, func(varName string) (value string) {
		switch varName {
		case "year":
			return strconv.Itoa(now.Year())
		case "month":
			return fmt.Sprintf("%02d", now.Month())
		case "week":
			_, week := now.ISOWeek()
			return fmt.Sprintf("%02d", week)
		case "day":
			return fmt.Sprintf("%02d", now.Day())
		case "hour":
			return fmt.Sprintf("%02d", now.Hour())
		case "minute":
			return fmt.Sprintf("%02d", now.Minute())
		case "second":
			return fmt.Sprintf("%02d", now.Second())
		case "date":
			return fmt.Sprintf("%d%02d%02d", now.Year(), now.Month(), now.Day())
		}

		return varName
	})
}
