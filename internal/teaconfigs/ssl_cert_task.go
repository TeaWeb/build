package teaconfigs

import "github.com/iwind/TeaGo/utils/string"

// 证书生成任务
type SSLCertTask struct {
	Id       string `yaml:"id" json:"id"`             // ID
	On       bool   `yaml:"on" json:"on"`             // 是否启用
	RunAt    int64  `yaml:"runAt" json:"runAt"`       // 运行时间
	RunError string `yaml:"runError" json:"runError"` // 运行错误

	Request *ACMERequest `yaml:"request" json:"request"` // ACME信息
}

// 获取新对象
func NewSSLCertTask() *SSLCertTask {
	return &SSLCertTask{
		Id: stringutil.Rand(16),
		On: true,
	}
}
