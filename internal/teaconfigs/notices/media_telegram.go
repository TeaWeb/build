package notices

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/iwind/TeaGo/types"
)

// Telegram媒介
type NoticeTelegramMedia struct {
	Token string `yaml:"token" json:"token"`
}

// 获取新对象
func NewNoticeTelegramMedia() *NoticeTelegramMedia {
	return &NoticeTelegramMedia{}
}

// 发送消息
func (this *NoticeTelegramMedia) Send(user string, subject string, body string) (respBytes []byte, err error) {
	bot, err := tgbotapi.NewBotAPI(this.Token)
	if err != nil {
		return nil, err
	}

	msg := tgbotapi.NewMessage(types.Int64(user), subject+"\n"+body)
	_, err = bot.Send(msg)
	return nil, err
}

// 是否需要用户标识
func (this *NoticeTelegramMedia) RequireUser() bool {
	return true
}
