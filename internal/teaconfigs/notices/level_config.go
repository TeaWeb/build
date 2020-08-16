package notices

// 级别配置
type NoticeLevelConfig struct {
	ShouldNotify bool              `yaml:"shouldNotify" json:"shouldNotify"`
	Receivers    []*NoticeReceiver `yaml:"receivers" json:"receivers"`
}

// 添加接收人
func (this *NoticeLevelConfig) AddReceiver(receiver *NoticeReceiver) {
	this.Receivers = append(this.Receivers, receiver)
}

// 移除接收人
func (this *NoticeLevelConfig) RemoveReceiver(receiverId string) {
	result := []*NoticeReceiver{}
	for _, r := range this.Receivers {
		if r.Id == receiverId {
			continue
		}
		result = append(result, r)
	}
	this.Receivers = result
}

// 移除某个媒介的所有接收人
func (this *NoticeLevelConfig) RemoveMediaReceivers(mediaId string) {
	result := []*NoticeReceiver{}
	for _, r := range this.Receivers {
		if r.MediaId == mediaId {
			continue
		}
		result = append(result, r)
	}
	this.Receivers = result
}

// 查找单个接收人
func (this *NoticeLevelConfig) FindReceiver(receiverId string) *NoticeReceiver {
	for _, r := range this.Receivers {
		if r.Id == receiverId {
			return r
		}
	}
	return nil
}
