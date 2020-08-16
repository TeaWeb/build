package teaconfigs

import (
	"github.com/TeaWeb/build/internal/teaconfigs/shared"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/files"
	"github.com/iwind/TeaGo/logs"
	"gopkg.in/yaml.v3"
)

// 节点配置文件名
var nodeConfigFile = "node.conf"
var sharedNodeConfig *NodeConfig = nil

// 节点配置
type NodeConfig struct {
	Id            string   `yaml:"id" json:"id"`                       // ID
	On            bool     `yaml:"on" json:"on"`                       // 是否启用
	Name          string   `yaml:"name" json:"name"`                   // 名称
	ClusterId     string   `yaml:"clusterId" json:"clusterId"`         // 集群ID
	ClusterSecret string   `yaml:"clusterSecret" json:"clusterSecret"` // 集群秘钥
	ClusterAddr   string   `yaml:"clusterAddr" json:"clusterAddr"`     // 集群通讯地址
	Role          NodeRole `yaml:"role" json:"role"`                   // 角色
}

// 取得当前节点配置
// 如果为nil，表示尚未配置集群
func SharedNodeConfig() *NodeConfig {
	shared.Locker.Lock()
	defer shared.Locker.ReadUnlock()

	if sharedNodeConfig != nil {
		return sharedNodeConfig
	}

	configFile := files.NewFile(Tea.ConfigFile(nodeConfigFile))
	if !configFile.Exists() {
		return nil
	}

	data, err := configFile.ReadAll()
	if err != nil {
		logs.Error(err)
		return nil
	}

	config := &NodeConfig{}
	err = yaml.Unmarshal(data, config)
	if err != nil {
		return nil
	}

	sharedNodeConfig = config
	return config
}

// 保存到文件
func (this *NodeConfig) Save() error {
	shared.Locker.Lock()
	defer shared.Locker.WriteUnlock()

	data, err := yaml.Marshal(this)
	if err != nil {
		return err
	}

	configFile := files.NewFile(Tea.ConfigFile(nodeConfigFile))
	err = configFile.Write(data)
	if err != nil {
		return err
	}
	sharedNodeConfig = nil
	return nil
}

// 是否为Master
func (this *NodeConfig) IsMaster() bool {
	return this.Role == NodeRoleMaster
}
