package teaconfigs

import "github.com/iwind/TeaGo/lists"

// 代理访问日志配置
type AccessLogConfig struct {
	Id string `yaml:"id" json:"id"`
	On bool   `yaml:"on" json:"on"`

	Fields []int `yaml:"fields" json:"fields"` // 记录的字段

	Status1 bool `yaml:"status1" json:"status1"` // 1xx
	Status2 bool `yaml:"status2" json:"status2"` // 2xx
	Status3 bool `yaml:"status3" json:"status3"` // 3xx
	Status4 bool `yaml:"status4" json:"status4"` // 4xx
	Status5 bool `yaml:"status5" json:"status5"` // 5xx

	StorageOnly     bool     `yaml:"storageOnly" json:"storageOnly"`         // 是否只输出到存储策略
	StoragePolicies []string `yaml:"storagePolicies" json:"storagePolicies"` // 存储策略Ids
}

// 获取新对象
func NewAccessLogConfig() *AccessLogConfig {
	return &AccessLogConfig{
		On:      true,
		Fields:  []int{},
		Status1: true,
		Status2: true,
		Status3: true,
		Status4: true,
		Status5: true,
	}
}

// 校验
func (this *AccessLogConfig) Validate() error {
	return nil
}

// 判断是否应该记录
func (this *AccessLogConfig) Match(status int) bool {
	s := status / 100
	switch s {
	case 1:
		if !this.Status1 {
			return false
		}
	case 2:
		if !this.Status2 {
			return false
		}
	case 3:
		if !this.Status3 {
			return false
		}
	case 4:
		if !this.Status4 {
			return false
		}
	case 5:
		if !this.Status5 {
			return false
		}
	}

	return true
}

// 是否包含某个存储策略
func (this *AccessLogConfig) ContainsStoragePolicy(storagePolicyId string) bool {
	return lists.ContainsString(this.StoragePolicies, storagePolicyId)
}
