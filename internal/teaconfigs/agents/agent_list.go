package agents

import (
	"github.com/TeaWeb/build/internal/teaconfigs/shared"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/files"
	"github.com/iwind/TeaGo/logs"
	"sync"
)

// Agent列表是否有变化
var agentListChanged = false
var agentList = []*AgentConfig{}
var agentListLocker = sync.Mutex{}

// Agent列表
type AgentList struct {
	Files       []string `yaml:"files" json:"files"`
	filesLocker sync.Mutex
}

// 取得Agent列表
func SharedAgentList() (*AgentList, error) {
	file := files.NewFile(Tea.ConfigFile("agents/agentlist.conf"))
	if !file.Exists() {
		// 创建目录
		dir := files.NewFile(Tea.ConfigFile("agents"))
		if !dir.Exists() {
			err := dir.MkdirAll()
			if err != nil {
				logs.Error(err)
			}
		}

		return &AgentList{}, nil
	}
	reader, err := file.Reader()
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = reader.Close()
	}()
	agentList := &AgentList{}
	err = reader.ReadYAML(agentList)
	if err != nil {
		return nil, err
	}
	return agentList, nil
}

// 取得AgentId列表，不包括Local
func SharedAgents() []*AgentConfig {
	agentListLocker.Lock()
	defer agentListLocker.Unlock()

	if !agentListChanged && len(agentList) > 0 {
		return agentList
	}

	agentList = []*AgentConfig{}
	list, _ := SharedAgentList()
	for _, agent := range list.FindAllAgents() {
		agentList = append(agentList, agent)
	}

	agentListChanged = false
	return agentList
}

// 取得AgentId列表，包括Local
func AllSharedAgents() []*AgentConfig {
	agents := SharedAgents()
	return append(agents, NewAgentConfigFromId("local"))
}

// 通知Agent变化
func NotifyAgentsChange() {
	agentListLocker.Lock()
	defer agentListLocker.Unlock()
	agentListChanged = true
}

// 添加Agent
func (this *AgentList) AddAgent(agentFile string) {
	this.filesLocker.Lock()
	defer this.filesLocker.Unlock()
	this.Files = append(this.Files, agentFile)
}

// 删除Agent
func (this *AgentList) RemoveAgent(agentFile string) {
	result := []string{}
	for _, f := range this.Files {
		if f == agentFile {
			continue
		}
		result = append(result, f)
	}
	this.Files = result
}

// 查找所有Agents
func (this *AgentList) FindAllAgents() []*AgentConfig {
	result := []*AgentConfig{}
	for _, f := range this.Files {
		agent := NewAgentConfigFromFile(f)
		if agent == nil {
			continue
		}
		result = append(result, agent)
	}
	return result
}

// 计算所有分组中的Agent
func (this *AgentList) CountAgentsInGroup(groupId string) int {
	if len(groupId) == 0 {
		groupId = "default"
	}
	count := 0
	for _, agent := range this.FindAllAgents() {
		if agent.BelongsToGroup(groupId) {
			count++
		}
	}
	return count
}

// 保存
func (this *AgentList) Save() error {
	shared.Locker.Lock()
	defer shared.Locker.WriteUnlockNotify()

	defer func() {
		NotifyAgentsChange()
	}()

	writer, err := files.NewWriter(Tea.ConfigFile("agents/agentlist.conf"))
	if err != nil {
		return err
	}
	defer func() {
		_ = writer.Close()
	}()
	_, err = writer.WriteYAML(this)
	return err
}

// 移动位置
func (this *AgentList) MoveAgent(fromId string, toId string) {
	fromIndex := -1
	toIndex := -1

	for index, f := range this.Files {
		if f == "agent."+fromId+".conf" {
			fromIndex = index
		}
		if f == "agent."+toId+".conf" {
			toIndex = index
		}
	}

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
	for i := 0; i < len(this.Files); i++ {
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
