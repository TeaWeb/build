package api

import (
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/files"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/rands"
)

// 测试历史
type APITestCase struct {
	Filename     string     `yaml:"filename" json:"filename"`         // 文件名
	Name         string     `yaml:"name" json:"name"`                 // 名称
	Domain       string     `yaml:"domain" json:"domain"`             // 域名
	Method       string     `yaml:"method" json:"method"`             // 方法
	Query        string     `yaml:"query" json:"query"`               // URL附加参数
	Headers      []maps.Map `yaml:"headers" json:"headers"`           // Header
	Params       []maps.Map `yaml:"params" json:"params"`             // 内置参数
	AttachParams []maps.Map `yaml:"attachParams" json:"attachParams"` // 附加参数
	Format       string     `yaml:"format" json:"format"`             // 响应格式
	Username     string     `yaml:"username" json:"username"`         // 用户名
	CreatedAt    int64      `yaml:"createdAt" json:"createdAt"`       // 创建时间
	UpdatedAt    int64      `yaml:"updatedAt" json:"updatedAt"`       // 更新时间
}

// 获取新对象
func NewAPITestCase() *APITestCase {
	return &APITestCase{}
}

// 从文件中加载对象
func NewAPITestCaseFromFile(filename string) *APITestCase {
	reader, err := files.NewReader(Tea.ConfigFile(filename))
	if err != nil {
		logs.Error(err)
		return nil
	}
	defer reader.Close()

	h := NewAPITestCase()
	err = reader.ReadYAML(h)
	if err != nil {
		logs.Error(err)
		return nil
	}

	return h
}

// 保存到文件
func (this *APITestCase) Save() error {
	if len(this.Filename) == 0 {
		this.Filename = "test.case." + rands.HexString(16) + ".conf"
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
func (this *APITestCase) Delete() error {
	if len(this.Filename) > 0 {
		return files.NewFile(Tea.ConfigFile(this.Filename)).Delete()
	}
	return nil
}
