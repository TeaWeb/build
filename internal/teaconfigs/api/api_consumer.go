package api

import (
	"github.com/TeaWeb/build/internal/teaconfigs/shared"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/files"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/rands"
)

// API consumer
type APIConsumer struct {
	Filename string `yaml:"filename" json:"filename"` // 文件名
	On       bool   `yaml:"on" json:"on"`             // 是否开启 TODO
	Name     string `yaml:"name" json:"name"`         // 名称

	// 认证
	Auth struct {
		On      bool                   `yaml:"on" json:"on"`           // 是否开启 TODO
		Type    string                 `yaml:"type" json:"type"`       // 类型
		Options map[string]interface{} `yaml:"options" json:"options"` // 选项
	} `yaml:"auth" json:"auth"`

	// API控制
	API struct {
		On       bool     `yaml:"on" json:"on"`             // 是否开启
		AllowAll bool     `yaml:"allowAll" json:"allowAll"` // 是否允许所有
		DenyAll  bool     `yaml:"denyAll" json:"denyAll"`   // 是否禁止所有
		Allow    []string `yaml:"allow" json:"allow"`       // 允许的API
		Deny     []string `yaml:"deny" json:"deny"`         // 禁止的API
	} `yaml:"api" json:"api"` // API控制

	Policy shared.AccessPolicy `yaml:"policy" json:"policy"` // 控制策略
}

// 获取新对象
func NewAPIConsumer() *APIConsumer {
	return &APIConsumer{}
}

// 从文件中加载对象
func NewAPIConsumerFromFile(filename string) *APIConsumer {
	if len(filename) == 0 {
		return nil
	}
	reader, err := files.NewReader(Tea.ConfigFile(filename))
	if err != nil {
		logs.Error(err)
		return nil
	}
	defer reader.Close()

	c := NewAPIConsumer()
	err = reader.ReadYAML(c)
	if err != nil {
		logs.Error(err)
		return nil
	}
	return c
}

// 校验
func (this *APIConsumer) Validate() error {
	return nil
}

// 保存
func (this *APIConsumer) Save() error {
	if len(this.Filename) == 0 {
		this.Filename = "consumer." + rands.HexString(16) + ".conf"
	}
	writer, err := files.NewWriter(Tea.ConfigFile(this.Filename))
	if err != nil {
		return err
	}
	defer writer.Close()
	_, err = writer.WriteYAML(this)
	return err
}

// 删除
func (this *APIConsumer) Delete() error {
	if len(this.Filename) == 0 {
		return nil
	}
	f := files.NewFile(Tea.ConfigFile(this.Filename))
	return f.Delete()
}

// 消费API
func (this *APIConsumer) AllowAPI(apiPath string) (passed bool) {
	// API控制
	if this.API.On {
		// deny
		if this.API.DenyAll {
			return false
		}
		if len(this.API.Deny) > 0 && lists.Contains(this.API.Deny, apiPath) {
			return false
		}

		// allow
		if !this.API.AllowAll && !lists.Contains(this.API.Allow, apiPath) {
			return false
		}
	}

	return true
}
