package api

import (
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/files"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/utils/string"
)

// API的mock格式定义
const (
	APIMockFormatJSON = "json"
	APIMockFormatXML  = "xml"
	APIMockFormatText = "text"
	APIMockFormatFile = "file"
	// APIMockFormatJavascript = ""
)

// API Mock定义
type APIMock struct {
	Filename  string     `yaml:"filename" json:"filename"`   // 保存的文件名
	On        bool       `yaml:"on" json:"on"`               // 是否开启
	Headers   []maps.Map `yaml:"headers" json:"headers"`     // 输出的Header
	Format    string     `yaml:"format" json:"format"`       // 格式
	Text      string     `yaml:"text" json:"text"`           // 文本
	File      string     `yaml:"file" json:"file"`           // 文件名，一般是和文本二选一
	Username  string     `yaml:"username" json:"username"`   // 创建的用户名
	CreatedAt int64      `yaml:"createdAt" json:"createdAt"` // 创建时间
}

// 获取新对象
func NewAPIMock() *APIMock {
	return &APIMock{
		On: true,
	}
}

// 从文件中加载对象
func NewAPIMockFromFile(filename string) *APIMock {
	if len(filename) == 0 {
		return nil
	}

	reader, err := files.NewReader(Tea.ConfigFile(filename))
	if err != nil {
		return nil
	}
	defer reader.Close()

	mock := NewAPIMock()
	err = reader.ReadYAML(mock)
	if err != nil {
		logs.Error(err)
		return nil
	}

	return mock
}

// 保存
func (this *APIMock) Save() error {
	if len(this.Filename) == 0 {
		this.Filename = "mock." + stringutil.Rand(16) + ".conf"
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
func (this *APIMock) Delete() error {
	if len(this.Filename) == 0 {
		return nil
	}

	if len(this.File) > 0 {
		err := files.NewFile(Tea.ConfigFile(this.File)).Delete()
		if err != nil {
			return err
		}
	}
	return files.NewFile(Tea.ConfigFile(this.Filename)).Delete()
}
