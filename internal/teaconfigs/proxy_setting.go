package teaconfigs

import (
	"github.com/TeaWeb/build/internal/teaconfigs/shared"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/logs"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

const proxySettingFile = "proxy_setting.conf"

var sharedProxySetting *ProxySetting = nil

// 代理全局设置
type ProxySetting struct {
	MatchDomainStrictly bool `yaml:"matchDomainStrictly" json:"matchDomainStrictly"` // 是否严格匹配域名

}

// 取得共享的代理设置
func SharedProxySetting() *ProxySetting {
	shared.Locker.Lock()
	defer shared.Locker.ReadUnlock()
	if sharedProxySetting == nil {
		sharedProxySetting = LoadProxySetting()
	}
	return sharedProxySetting
}

// 从配置文件中加载配置
func LoadProxySetting() *ProxySetting {
	setting := &ProxySetting{
		MatchDomainStrictly: false,
	}

	data, err := ioutil.ReadFile(Tea.ConfigFile(proxySettingFile))
	if err != nil {
		return setting
	}

	err = yaml.Unmarshal(data, setting)
	if err != nil {
		logs.Error(err)
	}

	return setting
}

// 保存
func (this *ProxySetting) Save() error {
	data, err := yaml.Marshal(this)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(Tea.ConfigFile(proxySettingFile), data, 0666)
	if err == nil {
		shared.Locker.Lock()
		sharedProxySetting = this
		shared.Locker.ReadUnlock()
	}
	return err
}
