package notices

import "github.com/iwind/TeaGo/maps"

// 通知媒介类型
type NoticeMediaType = string

const (
	NoticeMediaTypeEmail         = "email"
	NoticeMediaTypeWebhook       = "webhook"
	NoticeMediaTypeScript        = "script"
	NoticeMediaTypeDingTalk      = "dingTalk"
	NoticeMediaTypeQyWeixin      = "qyWeixin"
	NoticeMediaTypeQyWeixinRobot = "qyWeixinRobot"
	NoticeMediaTypeAliyunSms     = "aliyunSms"
	NoticeMediaTypeTelegram      = "telegram"
	NoticeMediaTypeTeaSms        = "teaSms"
)

// 所有媒介
func AllNoticeMediaTypes() []maps.Map {
	return []maps.Map{
		{
			"name":         "邮件",
			"code":         NoticeMediaTypeEmail,
			"supportsHTML": true,
			"instance":     new(NoticeEmailMedia),
			"description":  "通过邮件发送通知",
			"user":         "接收人邮箱地址",
		},
		{
			"name":         "Webhook",
			"code":         NoticeMediaTypeWebhook,
			"supportsHTML": false,
			"instance":     new(NoticeWebhookMedia),
			"description":  "通过HTTP请求发送通知",
			"user":         "通过${NoticeUser}参数传递到URL上",
		},
		{
			"name":         "脚本",
			"code":         NoticeMediaTypeScript,
			"supportsHTML": false,
			"instance":     new(NoticeScriptMedia),
			"description":  "通过运行脚本发送通知",
			"user":         "可以在脚本中使用${NoticeUser}来获取这个标识",
		},
		{
			"name":         "钉钉群机器人",
			"code":         NoticeMediaTypeDingTalk,
			"supportsHTML": false,
			"instance":     new(NoticeDingTalkMedia),
			"description":  "通过钉钉群机器人发送通知消息，<a href=\"http://teaos.cn/doc/notices/DingTalk.md\" target=\"_blank\">相关文档&raquo;</a>",
			"user":         "要At（@）的群成员的手机号，多个手机号用英文逗号隔开，也可以为空",
		},
		{
			"name":         "企业微信应用",
			"code":         NoticeMediaTypeQyWeixin,
			"supportsHTML": false,
			"instance":     new(NoticeQyWeixinMedia),
			"description":  "通过企业微信应用发送通知消息，<a href=\"http://teaos.cn/doc/notices/QyWeixin.md\" target=\"_blank\">相关文档&raquo;</a>",
			"user":         "接收消息的成员的用户账号，多个成员用竖线（|）分隔，如果所有成员使用@all。留空表示所有成员。<a href=\"http://teaos.cn/doc/notices/QyWeixin.md#%E7%94%A8%E6%88%B7%E8%B4%A6%E5%8F%B7\" target=\"_blank\">如何查看成员的用户账号？</a>",
		},
		{
			"name":         "企业微信群机器人",
			"code":         NoticeMediaTypeQyWeixinRobot,
			"supportsHTML": false,
			"instance":     new(NoticeQyWeixinRobotMedia),
			"description":  "通过微信群机器人发送通知消息",
			"user":         "要At（@）的群成员的手机号，多个手机号用英文逗号隔开，也可以为空",
		},
		{
			"name":         "阿里云短信",
			"code":         NoticeMediaTypeAliyunSms,
			"supportsHTML": false,
			"instance":     new(NoticeAliyunSmsMedia),
			"description":  "通过<a href=\"https://www.aliyun.com/product/sms?spm=5176.11533447.1097531.2.12055cfa6UnIix\" target=\"_blank\">阿里云短信服务</a>发送短信，<a href=\"http://teaos.cn/doc/notices/AliyunSms.md\" target=\"_blank\">相关文档&raquo;</a>",
			"user":         "接收消息的手机号",
		},
		{
			"name":         "Telegram机器人",
			"code":         NoticeMediaTypeTelegram,
			"supportsHTML": false,
			"instance":     new(NoticeTelegramMedia),
			"description":  "通过机器人向群或者某个用户发送消息，需要确保所在网络能够访问Telegram API服务",
			"user":         "群或用户的Chat ID，通常是一个数字，可以通过和 @get_id_bot 建立对话并发送任意消息获得",
		},
		{
			"name":         "TeaOS云短信",
			"code":         NoticeMediaTypeTeaSms,
			"supportsHTML": false,
			"instance":     new(NoticeTeaSmsMedia),
			"description":  "通过<a href=\"http://cloud.teaos.cn\" target=\"_blank\">TeaOS官方</a>提供的云短信接口发送短信，<a href=\"http://teaos.cn/doc/notices/TeaSms.md\" target=\"_blank\">相关文档&raquo;</a>",
			"user":         "接收消息的手机号",
		},
	}
}

// 查找媒介类型
func FindNoticeMediaType(mediaType string) maps.Map {
	for _, m := range AllNoticeMediaTypes() {
		if m["code"] == mediaType {
			return m
		}
	}
	return nil
}

// 查找媒介类型名称
func FindNoticeMediaTypeName(mediaType string) string {
	m := FindNoticeMediaType(mediaType)
	if m == nil {
		return ""
	}
	return m["name"].(string)
}

// 媒介接口
type NoticeMediaInterface interface {
	// 发送
	Send(user string, subject string, body string) (resp []byte, err error)

	// 是否可以需要用户标识
	RequireUser() bool
}
