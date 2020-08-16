package teaconfigs

import (
	"github.com/TeaWeb/build/internal/teaconfigs/shared"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/files"
	"github.com/iwind/TeaGo/lists"
)

// 缓存管理
type CacheConfig struct {
	Filename    string   `yaml:"filename" json:"filename"`       // 文件名
	PolicyFiles []string `yaml:"policyFiles" json:"policyFiles"` // 策略文件
}

// 获取新对象
func NewCacheConfig() *CacheConfig {
	return &CacheConfig{}
}

// 加载对象，无论如何都会返回一个对象
func SharedCacheConfig() (*CacheConfig, error) {
	reader, err := files.NewReader(Tea.ConfigFile("cache.conf"))
	if err != nil {
		return NewCacheConfig(), err
	}
	defer reader.Close()
	config := NewCacheConfig()
	err = reader.ReadYAML(config)
	if err != nil {
		return NewCacheConfig(), err
	}
	return config, nil
}

// 添加缓存策略
func (this *CacheConfig) AddPolicy(file string) {
	this.PolicyFiles = append(this.PolicyFiles, file)
}

// 删除缓存策略
func (this *CacheConfig) DeletePolicy(file string) {
	this.PolicyFiles = lists.Delete(this.PolicyFiles, file).([]string)
}

// 查找所有的缓存策略
func (this *CacheConfig) FindAllPolicies() []*shared.CachePolicy {
	result := []*shared.CachePolicy{}
	for _, file := range this.PolicyFiles {
		policy := shared.NewCachePolicyFromFile(file)
		if policy == nil {
			continue
		}
		policy.Validate()
		result = append(result, policy)
	}
	return result
}

// 保存
func (this *CacheConfig) Save() error {
	shared.Locker.Lock()
	defer shared.Locker.WriteUnlockNotify()

	if len(this.Filename) == 0 {
		this.Filename = "cache.conf"
	}
	writer, err := files.NewWriter(Tea.ConfigFile(this.Filename))
	if err != nil {
		return err
	}
	defer writer.Close()
	_, err = writer.WriteYAML(this)
	return err
}
