package agents

import (
	"github.com/TeaWeb/build/internal/teaconfigs/notices"
	"github.com/TeaWeb/build/internal/teaconfigs/shared"
	"github.com/TeaWeb/build/internal/teautils"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/rands"
	timeutil "github.com/iwind/TeaGo/utils/time"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

// Agent分组
type Group struct {
	Id            string                                            `yaml:"id" json:"id"`
	IsDefault     bool                                              `yaml:"isDefault" json:"isDefault"`
	On            bool                                              `yaml:"on" json:"on"`
	Name          string                                            `yaml:"name" json:"name"`
	Index         int                                               `yaml:"index" json:"index"`
	NoticeSetting map[notices.NoticeLevel][]*notices.NoticeReceiver `yaml:"noticeSetting" json:"noticeSetting"`
	Key           string                                            `yaml:"key" json:"key"` // 密钥

	DayFrom     string `yaml:"dayFrom" json:"dayFrom"`         // 有效开始日期
	DayTo       string `yaml:"dayTo" json:"dayTo"`             // 有效结束日期
	MaxAgents   int    `yaml:"maxAgents" json:"maxAgents"`     // 可以容纳的Agents最大数量
	CountAgents int    `yaml:"countAgents" json:"countAgents"` //

	Keys []*GroupKey `yaml:"keys" json:"keys"` // 临时的Key
}

// 获取新分组
func NewGroup(name string) *Group {
	return &Group{
		Id:   rands.HexString(16),
		On:   true,
		Name: name,
	}
}

// 默认的分组
// deprecated in v0.1.8
func loadOldDefaultGroup() *Group {
	data, err := ioutil.ReadFile(Tea.ConfigFile("agents/group.default.conf"))
	if err != nil {
		return nil
	}
	group := new(Group)
	err = yaml.Unmarshal(data, group)
	if err != nil {
		return nil
	}

	group.IsDefault = true
	return group
}

// 添加通知接收者
func (this *Group) AddNoticeReceiver(level notices.NoticeLevel, receiver *notices.NoticeReceiver) {
	if this.NoticeSetting == nil {
		this.NoticeSetting = map[notices.NoticeLevel][]*notices.NoticeReceiver{}
	}
	receivers, found := this.NoticeSetting[level]
	if !found {
		receivers = []*notices.NoticeReceiver{}
	}
	receivers = append(receivers, receiver)
	this.NoticeSetting[level] = receivers
}

// 删除通知接收者
func (this *Group) RemoveNoticeReceiver(level notices.NoticeLevel, receiverId string) {
	if this.NoticeSetting == nil {
		return
	}
	receivers, found := this.NoticeSetting[level]
	if !found {
		return
	}

	result := []*notices.NoticeReceiver{}
	for _, r := range receivers {
		if r.Id == receiverId {
			continue
		}
		result = append(result, r)
	}
	this.NoticeSetting[level] = result
}

// 删除媒介
func (this *Group) RemoveMedia(mediaId string) (found bool) {
	for level, receivers := range this.NoticeSetting {
		result := []*notices.NoticeReceiver{}
		for _, receiver := range receivers {
			if receiver.MediaId == mediaId {
				found = true
				continue
			}
			result = append(result, receiver)
		}
		this.NoticeSetting[level] = result
	}
	return
}

// 查找一个或多个级别对应的接收者，并合并相同的接收者
func (this *Group) FindAllNoticeReceivers(level ...notices.NoticeLevel) []*notices.NoticeReceiver {
	if len(level) == 0 {
		return []*notices.NoticeReceiver{}
	}

	m := maps.Map{} // mediaId_user => bool
	result := []*notices.NoticeReceiver{}
	for _, l := range level {
		receivers, ok := this.NoticeSetting[l]
		if !ok {
			continue
		}
		for _, receiver := range receivers {
			if !receiver.On {
				continue
			}
			key := receiver.Key()
			if m.Has(key) {
				continue
			}
			m[key] = true
			result = append(result, receiver)
		}
	}
	return result
}

// 写入文件
func (this *Group) WriteToFile(path string) error {
	shared.Locker.Lock()
	defer shared.Locker.WriteUnlockNotify()
	data, err := yaml.Marshal(this)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, data, 0666)
}

// 生成密钥
func (this *Group) GenerateKey() string {
	return rands.HexString(32)
}

// 匹配关键词
func (this *Group) MatchKeyword(keyword string) (matched bool, name string, tags []string) {
	if teautils.MatchKeyword(this.Name, keyword) {
		return true, this.Name, nil
	}
	return
}

// 添加密钥
func (this *Group) AddKey(key *GroupKey) {
	this.Keys = append(this.Keys, key)
}

// 查找密钥
func (this *Group) FindKey(key string) *GroupKey {
	for _, k := range this.Keys {
		if k.Key == key {
			return k
		}
	}
	return nil
}

// 当前分组日期是否可用
func (this *Group) IsDateAvailable() bool {
	today := timeutil.Format("Y-m-d")
	if len(this.DayFrom) > 0 && this.DayFrom > today {
		return false
	}
	if len(this.DayTo) > 0 && this.DayTo < today {
		return false
	}
	return true
}

// 是否可用
func (this *Group) IsAvailable() bool {
	if !this.On {
		return false
	}
	if this.MaxAgents > 0 && this.CountAgents >= this.MaxAgents {
		return false
	}
	return this.IsDateAvailable()
}
