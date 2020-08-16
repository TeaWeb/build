package agentconfigs

import (
	"github.com/go-yaml/yaml"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/files"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/types"
	"io/ioutil"
	"os"
)

type AgentConfig struct {
	Master   string `yaml:"master" json:"master"`
	Id       string `yaml:"id" json:"id"`
	Key      string `yaml:"key" json:"key"`
	GroupKey string `yaml:"group" json:"group"`
}

func SharedAgentConfig() (*AgentConfig, error) {
	// 是否是本地
	isLocal := files.NewFile(Tea.ConfigFile("server.conf")).Exists()
	if isLocal {
		agent := &AgentConfig{
			Id: "local",
		}

		serverReader, err := files.NewReader(Tea.ConfigFile("server.conf"))
		if err != nil {
			return nil, err
		}
		defer func() {
			_ = serverReader.Close()
		}()

		m := maps.Map{}
		err = serverReader.ReadYAML(&m)
		if err != nil {
			return nil, err
		}
		httpConfig := m.GetMap("http")
		if httpConfig.GetBool("on") {
			listenConfig := httpConfig.GetSlice("listen")
			if len(listenConfig) != 0 {
				agent.Master = "http://" + types.String(listenConfig[0])
			}
		}

		if len(agent.Master) == 0 {
			httpsConfig := m.GetMap("https")
			if httpsConfig.GetBool("on") {
				listenConfig := httpsConfig.GetSlice("listen")
				if len(listenConfig) > 0 {
					agent.Master = "https://" + types.String(listenConfig[0])
				}
			}
		}

		return agent, nil
	}

	reader, err := files.NewReader(Tea.ConfigFile("agent.conf"))
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = reader.Close()
	}()

	a := &AgentConfig{}
	err = reader.ReadYAML(a)
	if err != nil {
		return nil, err
	}
	return a, nil
}

// 保存到文件
func (this *AgentConfig) Save() error {
	data, err := yaml.Marshal(this)
	if err != nil {
		return err
	}

	_, err = os.Stat(Tea.ConfigDir())
	if err != nil {
		_ = os.Mkdir(Tea.ConfigDir(), 0777)
	}

	return ioutil.WriteFile(Tea.ConfigFile("agent.conf"), data, 0777)
}
