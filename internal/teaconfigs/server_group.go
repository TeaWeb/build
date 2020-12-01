package teaconfigs

import (
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/rands"
)

type ServerGroup struct {
	Id        string   `yaml:"id" json:"id"`
	IsOn      bool     `yaml:"isOn" json:"isOn"`
	Name      string   `yaml:"name" json:"name"`
	ServerIds []string `yaml:"serverIds" json:"serverIds"`
}

func NewServerGroup() *ServerGroup {
	return &ServerGroup{
		Id:        rands.HexString(16),
		IsOn:      true,
		ServerIds: []string{},
	}
}

func (this *ServerGroup) Add(serverId ...string) {
	for _, id := range serverId {
		if lists.ContainsString(this.ServerIds, id) {
			continue
		}
		this.ServerIds = append(this.ServerIds, id)
	}
}

func (this *ServerGroup) Remove(serverId string) {
	result := []string{}
	for _, id := range this.ServerIds {
		if id == serverId {
			continue
		}
		result = append(result, id)
	}
	this.ServerIds = result
}

func (this *ServerGroup) Contains(serverId string) bool {
	return lists.ContainsString(this.ServerIds, serverId)
}
