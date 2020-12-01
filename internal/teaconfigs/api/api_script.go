package api

import (
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/files"
	"github.com/iwind/TeaGo/rands"
)

// 脚本定义
type APIScript struct {
	Filename string `yaml:"filename" json:"filename"` // 脚本路径
	Code     string `yaml:"code" json:"code"`         // 代码
}

// 获取新脚本
func NewAPIScript() *APIScript {
	return &APIScript{}
}

// 保存
func (this *APIScript) Save() error {
	if len(this.Filename) == 0 {
		this.Filename = "script." + rands.HexString(16) + ".conf"
	}
	writer, err := files.NewFile(Tea.ConfigFile(this.Filename)).Writer()
	if err != nil {
		return err
	}
	defer writer.Close()
	_, err = writer.WriteYAML(this)
	return err
}

// 删除
func (this *APIScript) Delete() error {
	if len(this.Filename) == 0 {
		return nil
	}

	return files.NewFile(Tea.ConfigFile(this.Filename)).Delete()
}
