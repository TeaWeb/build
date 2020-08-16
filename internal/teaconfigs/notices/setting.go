package notices

import (
	"github.com/TeaWeb/build/internal/teaconfigs/shared"
	"github.com/TeaWeb/build/internal/teaconst"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/files"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/types"
	"strings"
)

// 通知设置
type NoticeSetting struct {
	Levels map[NoticeLevel]*NoticeLevelConfig `yaml:"levels" json:"levels"`
	Medias []*NoticeMediaConfig               `yaml:"medias" json:"medias"`

	SoundOn bool `yaml:"soundOn" json:"soundOn"` // 提示声音
}

// 取得当前的配置
func SharedNoticeSetting() *NoticeSetting {
	filename := "notice.conf"
	file := files.NewFile(Tea.ConfigFile(filename))
	config := &NoticeSetting{
		Levels: map[NoticeLevel]*NoticeLevelConfig{},
		Medias: []*NoticeMediaConfig{},
	}
	if !file.Exists() {
		return config
	}

	reader, err := file.Reader()
	if err != nil {
		logs.Error(err)
		return config
	}
	defer func() {
		_ = reader.Close()
	}()
	err = reader.ReadYAML(config)
	if err != nil {
		logs.Error(err)
		return config
	}
	return config
}

// 保存配置
func (this *NoticeSetting) Save() error {
	shared.Locker.Lock()
	defer shared.Locker.WriteUnlockNotify()

	filename := "notice.conf"
	writer, err := files.NewWriter(Tea.ConfigFile(filename))
	if err != nil {
		return err
	}
	defer func() {
		_ = writer.Close()
	}()
	_, err = writer.WriteYAML(this)
	return err
}

// 查找级别配置
func (this *NoticeSetting) LevelConfig(level NoticeLevel) *NoticeLevelConfig {
	config, found := this.Levels[level]
	if found {
		return config
	}
	config = &NoticeLevelConfig{
		ShouldNotify: true,
	}
	this.Levels[level] = config
	return config
}

// 添加媒介配置
func (this *NoticeSetting) AddMedia(mediaConfig *NoticeMediaConfig) {
	this.Medias = append(this.Medias, mediaConfig)
}

// 删除媒介
func (this *NoticeSetting) RemoveMedia(mediaId string) {
	medias := []*NoticeMediaConfig{}
	for _, m := range this.Medias {
		if m.Id == mediaId {
			continue
		}
		medias = append(medias, m)
	}
	this.Medias = medias

	// 移除关联的接收人
	for _, l := range this.Levels {
		l.RemoveMediaReceivers(mediaId)
	}
}

// 查找媒介
func (this *NoticeSetting) FindMedia(mediaId string) *NoticeMediaConfig {
	for _, m := range this.Medias {
		if m.Id == mediaId {
			err := m.Validate()
			if err != nil {
				logs.Error(err)
			}
			return m
		}
	}
	return nil
}

// 查找接收人
func (this *NoticeSetting) FindReceiver(receiverId string) (level NoticeLevel, receiver *NoticeReceiver) {
	for levelCode, levelConfig := range this.Levels {
		receiver := levelConfig.FindReceiver(receiverId)
		if receiver != nil {
			return levelCode, receiver
		}
	}
	return 0, nil
}

// 发送通知
func (this *NoticeSetting) Notify(level NoticeLevel, subject string, message string, counter func(receiverId string, minutes int) int) (receiverIds []string) {
	config, found := this.Levels[level]
	if !found {
		return
	}
	this.NotifyReceivers(level, config.Receivers, subject, message, counter)
	return
}

// 发送通知给一组接收者
func (this *NoticeSetting) NotifyReceivers(level NoticeLevel, receivers []*NoticeReceiver, subject string, message string, counter func(receiverId string, minutes int) int) (receiverIds []string) {
	for _, r := range receivers {
		if !r.On {
			continue
		}
		media := this.FindMedia(r.MediaId)
		if media == nil || !media.On {
			continue
		}
		mediaType := FindNoticeMediaType(media.Type)
		if mediaType == nil {
			continue
		}
		raw, err := media.Raw()
		if err != nil {
			logs.Error(err)
			continue
		}
		if !media.ShouldNotify(counter(r.Id, media.RateMinutes)) {
			continue
		}
		receiverIds = append(receiverIds, r.Id)
		go func(raw NoticeMediaInterface, mediaType maps.Map, user string) {
			body := message
			if types.Bool(mediaType["supportsHTML"]) {
				body = strings.Replace(body, "\n", "<br/>", -1)
			}
			subjectContent := subject
			if len(subjectContent) == 0 {
				subjectContent = "[" + teaconst.TeaProductName + "][" + FindNoticeLevelName(level) + "]有新的通知"
			} else if !strings.HasPrefix(subject, "["+teaconst.TeaProductName+"]") {
				subjectContent = "[" + teaconst.TeaProductName + "][" + FindNoticeLevelName(level) + "]" + subject
			}
			_, err := raw.Send(user, subjectContent, body)
			if err != nil {
				logs.Error(err)
			}
		}(raw, mediaType, r.User)
	}

	return
}

// 查找一个或多个级别对应的接收者，并合并相同的接收者
func (this *NoticeSetting) FindAllNoticeReceivers(level ...NoticeLevel) []*NoticeReceiver {
	if len(level) == 0 {
		return []*NoticeReceiver{}
	}

	m := maps.Map{} // mediaId_user => bool
	result := []*NoticeReceiver{}
	for _, l := range level {
		config, ok := this.Levels[l]
		if !ok {
			continue
		}
		for _, receiver := range config.Receivers {
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
