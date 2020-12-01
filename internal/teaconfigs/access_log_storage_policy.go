package teaconfigs

import (
	"github.com/TeaWeb/build/internal/teaconfigs/shared"
	"github.com/TeaWeb/build/internal/teautils"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/rands"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
)

// 日志存储策略
// 存储在configs/accesslog.storage.$id.conf
type AccessLogStoragePolicy struct {
	Id      string                 `yaml:"id" json:"id"`
	Name    string                 `yaml:"name" json:"name"`
	On      bool                   `yaml:"bool" json:"on"`
	Type    string                 `yaml:"type" json:"type"`       // 存储类型
	Options map[string]interface{} `yaml:"options" json:"options"` // 存储选项
	Cond    []*shared.RequestCond  `yaml:"cond" json:"cond"`       // 请求条件
}

// 创建新策略
func NewAccessLogStoragePolicy() *AccessLogStoragePolicy {
	return &AccessLogStoragePolicy{
		Id: rands.HexString(16),
		On: true,
	}
}

// 从文件中加载策略
func NewAccessLogStoragePolicyFromId(id string) *AccessLogStoragePolicy {
	filename := "accesslog.storage." + id + ".conf"
	data, err := ioutil.ReadFile(Tea.ConfigFile(filename))
	if err != nil {
		logs.Error(err)
		return nil
	}
	policy := NewAccessLogStoragePolicy()
	err = yaml.Unmarshal(data, policy)
	if err != nil {
		logs.Error(err)
		return nil
	}
	return policy
}

// 校验
func (this *AccessLogStoragePolicy) Validate() error {
	// cond
	if len(this.Cond) > 0 {
		for _, cond := range this.Cond {
			err := cond.Validate()
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// 保存
func (this *AccessLogStoragePolicy) Save() error {
	shared.Locker.Lock()
	defer shared.Locker.WriteUnlockNotify()

	data, err := yaml.Marshal(this)
	if err != nil {
		return err
	}

	filename := "accesslog.storage." + this.Id + ".conf"
	return ioutil.WriteFile(Tea.ConfigFile(filename), data, 0666)
}

// 删除
func (this *AccessLogStoragePolicy) Delete() error {
	filename := "accesslog.storage." + this.Id + ".conf"
	return os.Remove(Tea.ConfigFile(filename))
}

// 匹配关键词
func (this *AccessLogStoragePolicy) MatchKeyword(keyword string) (matched bool, name string, tags []string) {
	if teautils.MatchKeyword(this.Name, keyword) || teautils.MatchKeyword(this.Type, keyword) {
		matched = true
		name = this.Name
		if len(this.Type) > 0 {
			tags = []string{"类型：" + this.Type}
		}
	}
	return
}

// 匹配条件
func (this *AccessLogStoragePolicy) MatchConds(formatter func(string) string) bool {
	if len(this.Cond) == 0 {
		return true
	}
	for _, cond := range this.Cond {
		if !cond.Match(formatter) {
			return false
		}
	}
	return true
}
