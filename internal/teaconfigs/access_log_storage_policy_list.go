package teaconfigs

import (
	"github.com/TeaWeb/build/internal/teaconfigs/shared"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/files"
	"github.com/iwind/TeaGo/logs"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

var accessLogStoragePolicyListFilename = "accesslog.storage.list.conf"

// 获取共享的存储策略列表
func SharedAccessLogStoragePolicyList() *AccessLogStoragePolicyList {
	path := Tea.ConfigFile(accessLogStoragePolicyListFilename)
	file := files.NewFile(path)
	if !file.Exists() {
		return &AccessLogStoragePolicyList{}
	}
	reader, err := file.Reader()
	if err != nil {
		logs.Error(err)
		return &AccessLogStoragePolicyList{}
	}
	defer reader.Close()
	policyList := &AccessLogStoragePolicyList{}
	err = reader.ReadYAML(policyList)
	if err != nil {
		logs.Error(err)
		return policyList
	}
	return policyList
}

// 存储策略列表
type AccessLogStoragePolicyList struct {
	Ids []string `yaml:"id" json:"id"`
}

// 添加策略ID
func (this *AccessLogStoragePolicyList) AddId(id string) {
	this.Ids = append(this.Ids, id)
}

// 删除策略ID
func (this *AccessLogStoragePolicyList) RemoveId(id string) {
	result := []string{}
	for _, id2 := range this.Ids {
		if id2 == id {
			continue
		}
		result = append(result, id2)
	}
	this.Ids = result
}

// 保存
func (this *AccessLogStoragePolicyList) Save() error {
	data, err := yaml.Marshal(this)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(Tea.ConfigFile(accessLogStoragePolicyListFilename), data, 0666)
}

// 查找所有的策略列表
func (this *AccessLogStoragePolicyList) FindAllPolicies() []*AccessLogStoragePolicy {
	shared.Locker.Lock()
	defer shared.Locker.WriteUnlock()

	result := []*AccessLogStoragePolicy{}
	for _, id := range this.Ids {
		policy := NewAccessLogStoragePolicyFromId(id)
		if policy != nil {
			result = append(result, policy)
		}
	}
	return result
}
