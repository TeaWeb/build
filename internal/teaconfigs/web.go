package teaconfigs

import (
	"github.com/TeaWeb/build/internal/teaconfigs/shared"
	"github.com/iwind/TeaGo"
	"github.com/iwind/TeaGo/Tea"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

const (
	webConfigFile = "server.conf"
)

// Web 配置
type WebConfig struct {
	TeaGo.ServerConfig `yaml:",inline"`

	CertId string `yaml:"certId" json:"certId"` // 代理全局中的证书
}

// 加载配置
func LoadWebConfig() (*WebConfig, error) {
	shared.Locker.Lock()
	defer shared.Locker.ReadUnlock()

	data, err := ioutil.ReadFile(Tea.ConfigFile(webConfigFile))
	if err != nil {
		return nil, err
	}
	config := &WebConfig{}
	err = yaml.Unmarshal(data, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

// 保存
func (this *WebConfig) Save() error {
	shared.Locker.Lock()
	defer shared.Locker.WriteUnlock()

	data, err := yaml.Marshal(this)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(Tea.ConfigFile(webConfigFile), data, 0777)
}
