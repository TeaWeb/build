package teaconfigs

import (
	"github.com/TeaWeb/build/internal/teaconfigs/shared"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/logs"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
)

type ServerGroupList struct {
	Groups []*ServerGroup `yaml:"groups" json:"groups"`
}

func SharedServerGroupList() *ServerGroupList {
	list := &ServerGroupList{
		Groups: []*ServerGroup{},
	}
	filename := "servergrouplist.conf"
	data, err := ioutil.ReadFile(Tea.ConfigFile(filename))
	if err != nil {
		if !os.IsNotExist(err) {
			logs.Error(err)
		}
		return list
	}

	err = yaml.Unmarshal(data, list)
	if err != nil {
		logs.Error(err)
		return list
	}

	return list
}

func (this *ServerGroupList) Add(group *ServerGroup) {
	this.Groups = append(this.Groups, group)
}

func (this *ServerGroupList) Remove(id string) {
	result := []*ServerGroup{}
	for _, group := range this.Groups {
		if group.Id == id {
			continue
		}
		result = append(result, group)
	}
	this.Groups = result
}

// 查找分组
func (this *ServerGroupList) Find(id string) *ServerGroup {
	for _, group := range this.Groups {
		if group.Id == id {
			return group
		}
	}
	return nil
}

// 检查是否包含某个代理服务ID
func (this *ServerGroupList) ContainsServer(serverId string) bool {
	for _, group := range this.Groups {
		for _, id := range group.ServerIds {
			if id == serverId {
				return true
			}
		}
	}
	return false
}

// 保存
func (this *ServerGroupList) Save() error {
	data, err := yaml.Marshal(this)
	if err != nil {
		return err
	}

	shared.Locker.Lock()
	defer shared.Locker.WriteUnlock()

	filename := "servergrouplist.conf"
	return ioutil.WriteFile(Tea.ConfigFile(filename), data, 0666)
}
