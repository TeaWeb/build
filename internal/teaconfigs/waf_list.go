package teaconfigs

import (
	"github.com/TeaWeb/build/internal/teaconfigs/shared"
	"github.com/TeaWeb/build/internal/teawaf"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/files"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/logs"
)

// 获取共享的WAF列表
func SharedWAFList() *WAFList {
	path := Tea.ConfigFile("waflist.conf")
	file := files.NewFile(path)
	if !file.Exists() {
		return &WAFList{}
	}
	reader, err := file.Reader()
	if err != nil {
		logs.Error(err)
		return &WAFList{}
	}
	defer reader.Close()
	wafList := &WAFList{}
	err = reader.ReadYAML(wafList)
	if err != nil {
		logs.Error(err)
		return wafList
	}
	return wafList
}

// WAF列表
type WAFList struct {
	Files []string `yaml:"files" json:"files"`
}

// 添加文件名
func (this *WAFList) AddFile(filename string) {
	if lists.ContainsString(this.Files, filename) {
		return
	}
	this.Files = append(this.Files, filename)
}

// 删除文件
func (this *WAFList) RemoveFile(filename string) {
	result := []string{}
	for _, file := range this.Files {
		if file == filename {
			continue
		}
		result = append(result, file)
	}
	this.Files = result
}

// 保存文件
func (this *WAFList) Save() error {
	shared.Locker.Lock()
	defer shared.Locker.WriteUnlockNotify()

	path := Tea.ConfigFile("waflist.conf")
	file := files.NewFile(path)
	writer, err := file.Writer()
	if err != nil {
		return err
	}
	defer writer.Close()
	_, err = writer.WriteYAML(this)
	return err
}

// 查找所有WAF配置
func (this *WAFList) FindAllConfigs() []*teawaf.WAF {
	result := []*teawaf.WAF{}
	for _, filename := range this.Files {
		path := Tea.ConfigFile(filename)
		waf, err := teawaf.NewWAFFromFile(path)
		if err != nil {
			logs.Error(err)
			continue
		}
		result = append(result, waf)
	}
	return result
}

// 查找单个WAF配置
func (this *WAFList) FindWAF(wafId string) *teawaf.WAF {
	if len(wafId) == 0 {
		return nil
	}
	filename := "waf." + wafId + ".conf"
	path := Tea.ConfigFile(filename)
	waf, err := teawaf.NewWAFFromFile(path)
	if err != nil {
		logs.Error(err)
		return nil
	}
	return waf
}

// 保存单个WAF配置
func (this *WAFList) SaveWAF(waf *teawaf.WAF) error {
	filename := "waf." + waf.Id + ".conf"
	return waf.Save(Tea.ConfigFile(filename))
}
