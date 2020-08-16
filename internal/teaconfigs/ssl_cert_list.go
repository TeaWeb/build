package teaconfigs

import (
	"github.com/TeaWeb/build/internal/teaconfigs/shared"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/logs"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

const (
	sslCertListFilename = "ssl.certs.conf"
)

// 获取证书列表实例
// 一定会返回不为nil的值
func SharedSSLCertList() *SSLCertList {
	data, err := ioutil.ReadFile(Tea.ConfigFile(sslCertListFilename))
	if err != nil {
		return NewSSLCertList()
	}

	list := &SSLCertList{}
	err = yaml.Unmarshal(data, list)
	if err != nil {
		logs.Error(err)
		return NewSSLCertList()
	}

	return list
}

// 公共的SSL证书列表
type SSLCertList struct {
	Certs []*SSLCertConfig `yaml:"certs" json:"certs"` // 证书
	Tasks []*SSLCertTask   `yaml:"tasks" json:"tasks"` // 任务
}

// 获取新对象
func NewSSLCertList() *SSLCertList {
	return &SSLCertList{
		Certs: []*SSLCertConfig{},
		Tasks: []*SSLCertTask{},
	}
}

// 添加证书
func (this *SSLCertList) AddCert(cert *SSLCertConfig) {
	this.Certs = append(this.Certs, cert)
}

// 删除证书
func (this *SSLCertList) RemoveCert(certId string) {
	result := []*SSLCertConfig{}
	for _, cert := range this.Certs {
		if cert.Id == certId {
			continue
		}
		result = append(result, cert)
	}
	this.Certs = result
}

// 查找证书
func (this *SSLCertList) FindCert(certId string) *SSLCertConfig {
	if len(certId) == 0 {
		return nil
	}
	for _, cert := range this.Certs {
		if cert.Id == certId {
			return cert
		}
	}
	return nil
}

// 添加任务
func (this *SSLCertList) AddTask(task *SSLCertTask) {
	this.Tasks = append(this.Tasks, task)
}

// 删除任务
func (this *SSLCertList) RemoveTask(taskId string) {
	result := []*SSLCertTask{}
	for _, task := range this.Tasks {
		if task.Id == taskId {
			continue
		}
		result = append(result, task)
	}
	this.Tasks = result
}

// 查找任务
func (this *SSLCertList) FindTask(taskId string) *SSLCertTask {
	if len(taskId) == 0 {
		return nil
	}
	for _, task := range this.Tasks {
		if task.Id == taskId {
			return task
		}
	}
	return nil
}

// 保存
func (this *SSLCertList) Save() error {
	shared.Locker.Lock()
	defer shared.Locker.WriteUnlockNotify()

	data, err := yaml.Marshal(this)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(Tea.ConfigFile(sslCertListFilename), data, 0777)
}
