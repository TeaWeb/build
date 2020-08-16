package teaconfigs

import (
	"github.com/TeaWeb/build/internal/teaconfigs/shared"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/files"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/logs"
	"sync"
)

// Locker
var serverListLocker = sync.Mutex{}

// Server列表
type ServerList struct {
	Files []string `yaml:"files" json:"files"`
}

// 取得Server列表
func SharedServerList() (*ServerList, error) {
	serverListLocker.Lock()
	defer serverListLocker.Unlock()

	file := files.NewFile(Tea.ConfigFile("serverlist.conf"))
	if !file.Exists() {
		// 初始化
		serverList := &ServerList{
			Files: []string{},
		}
		servers := LoadServerConfigsFromDir(Tea.ConfigDir())
		for _, s := range servers {
			serverList.Files = append(serverList.Files, s.Filename)
		}
		err := serverList.Save()
		if err != nil {
			logs.Error(err)
		}

		return serverList, nil
	}
	reader, err := file.Reader()
	if err != nil {
		return &ServerList{
			Files: []string{},
		}, err
	}
	defer reader.Close()
	serverList := &ServerList{}
	err = reader.ReadYAML(serverList)
	if err != nil {
		return nil, err
	}
	return serverList, nil
}

// 添加Server
func (this *ServerList) AddServer(serverFile string) {
	if !lists.ContainsString(this.Files, serverFile) {
		this.Files = append(this.Files, serverFile)
	}
}

// 删除Server
func (this *ServerList) RemoveServer(serverFile string) {
	result := []string{}
	for _, f := range this.Files {
		if f == serverFile {
			continue
		}
		result = append(result, f)
	}
	this.Files = result
}

// 查找所有Servers
func (this *ServerList) FindAllServers() []*ServerConfig {
	result := []*ServerConfig{}
	for _, f := range this.Files {
		server, err := NewServerConfigFromFile(f)
		if err != nil {
			logs.Error(err)
			continue
		}
		if server == nil {
			continue
		}
		result = append(result, server)
	}
	return result
}

// 移动位置
func (this *ServerList) MoveServer(fromIndex int, toIndex int) {
	if fromIndex < 0 || fromIndex >= len(this.Files) {
		return
	}
	if toIndex < 0 || toIndex >= len(this.Files) {
		return
	}
	if fromIndex == toIndex {
		return
	}

	file := this.Files[fromIndex]
	newList := []string{}
	for i := 0; i < len(this.Files); i ++ {
		if i == fromIndex {
			continue
		}
		if fromIndex > toIndex && i == toIndex {
			newList = append(newList, file)
		}
		newList = append(newList, this.Files[i])
		if fromIndex < toIndex && i == toIndex {
			newList = append(newList, file)
		}
	}

	this.Files = newList
}

// 保存
func (this *ServerList) Save() error {
	shared.Locker.Lock()
	defer shared.Locker.WriteUnlockNotify()

	writer, err := files.NewWriter(Tea.ConfigFile("serverlist.conf"))
	if err != nil {
		return err
	}
	defer writer.Close()
	_, err = writer.WriteYAML(this)
	return err
}
