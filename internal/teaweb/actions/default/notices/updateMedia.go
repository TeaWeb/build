package notices

import (
	"fmt"
	"github.com/TeaWeb/build/internal/teaconfigs/notices"
	"github.com/TeaWeb/build/internal/teaconfigs/shared"
	"github.com/TeaWeb/build/internal/teaconst"
	"github.com/TeaWeb/build/internal/teautils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/logs"
	"net/http"
)

type UpdateMediaAction actions.Action

// 修改媒介
func (this *UpdateMediaAction) Run(params struct {
	MediaId string
	From    string
}) {
	setting := notices.SharedNoticeSetting()
	media := setting.FindMedia(params.MediaId)
	if media == nil {
		this.Fail("找不到Media")
	}

	this.Data["from"] = params.From
	this.Data["media"] = media
	this.Data["mediaTypes"] = notices.AllNoticeMediaTypes()
	this.Data["methods"] = []string{http.MethodGet, http.MethodPost}

	this.Show()
}

// 提交修改
func (this *UpdateMediaAction) RunPost(params struct {
	MediaId string

	Name string
	Type string
	On   bool

	EmailSmtp     string
	EmailUsername string
	EmailPassword string
	EmailFrom     string

	WebhookURL          string
	WebhookMethod       string
	WebhookHeaderNames  []string
	WebhookHeaderValues []string
	WebhookContentType  string
	WebhookParamNames   []string
	WebhookParamValues  []string
	WebhookBody         string

	ScriptType      string
	ScriptPath      string
	ScriptLang      string
	ScriptCode      string
	ScriptCwd       string
	ScriptEnvNames  []string
	ScriptEnvValues []string

	DingTalkWebhookURL string

	QyWeixinCorporateId string
	QyWeixinAgentId     string
	QyWeixinAppSecret   string
	QyWeixinTextFormat  string

	QyWeixinRobotWebhookURL string
	QyWeixinRobotTextFormat string

	AliyunSmsSign              string
	AliyunSmsTemplateCode      string
	AliyunSmsTemplateVarNames  []string
	AliyunSmsTemplateVarValues []string
	AliyunSmsAccessKeyId       string
	AliyunSmsAccessKeySecret   string

	TelegramToken string

	TeaSmsAccessId     string
	TeaSmsAccessSecret string

	TimeFromHour   int
	TimeFromMinute int
	TimeFromSecond int
	TimeToHour     int
	TimeToMinute   int
	TimeToSecond   int
	RateCount      int
	RateMinutes    int

	Must *actions.Must
}) {
	if teaconst.DemoEnabled {
		this.Fail("演示版无法修改媒介")
	}

	params.Must.
		Field("name", params.Name).
		Require("请输入媒介名称")

	if notices.FindNoticeMediaType(params.Type) == nil {
		this.Fail("找不到此媒介类型")
	}

	setting := notices.SharedNoticeSetting()
	mediaConfig := setting.FindMedia(params.MediaId)
	if mediaConfig == nil {
		this.Fail("找不到Media")
	}

	mediaConfig.Name = params.Name
	mediaConfig.Type = params.Type
	mediaConfig.On = params.On

	switch params.Type {
	case notices.NoticeMediaTypeEmail:
		params.Must.
			Field("emailSmtp", params.EmailSmtp).
			Require("请输入SMTP地址").
			Field("emailUsername", params.EmailUsername).
			Require("请输入邮箱账号").
			Field("emailPassword", params.EmailPassword).
			Require("请输入密码或授权码")

		media := notices.NewNoticeEmailMedia()
		media.SMTP = params.EmailSmtp
		media.Username = params.EmailUsername
		media.Password = params.EmailPassword
		media.From = params.EmailFrom
		err := teautils.ObjectToMapJSON(media, &mediaConfig.Options)
		if err != nil {
			logs.Error(err)
		}
	case notices.NoticeMediaTypeWebhook:
		params.Must.
			Field("webhookURL", params.WebhookURL).
			Require("请输入URL地址").
			Match("(?i)^(http|https)://", "URL地址必须以http或https开头").
			Field("webhookMethod", params.WebhookMethod).
			Require("请选择请求方法")

		media := notices.NewNoticeWebhookMedia()
		media.URL = params.WebhookURL
		media.Method = params.WebhookMethod

		media.ContentType = params.WebhookContentType
		if len(params.WebhookHeaderNames) > 0 {
			for index, name := range params.WebhookHeaderNames {
				if index < len(params.WebhookHeaderValues) {
					media.AddHeader(name, params.WebhookHeaderValues[index])
				}
			}
		}

		if params.WebhookContentType == "params" {
			for index, name := range params.WebhookParamNames {
				if index < len(params.WebhookParamValues) {
					media.AddParam(name, params.WebhookParamValues[index])
				}
			}
		} else if params.WebhookContentType == "body" {
			media.Body = params.WebhookBody
		}
		err := teautils.ObjectToMapJSON(media, &mediaConfig.Options)
		if err != nil {
			logs.Error(err)
		}
	case notices.NoticeMediaTypeScript:
		if params.ScriptType == "path" {
			params.Must.
				Field("scriptPath", params.ScriptPath).
				Require("请输入脚本路径")
		} else if params.ScriptType == "code" {
			params.Must.
				Field("scriptCode", params.ScriptCode).
				Require("请输入脚本代码")
		} else {
			params.Must.
				Field("scriptPath", params.ScriptPath).
				Require("请输入脚本路径")
		}

		media := notices.NewNoticeScriptMedia()
		media.ScriptType = params.ScriptType
		media.Path = params.ScriptPath
		media.ScriptLang = params.ScriptLang
		media.Script = params.ScriptCode
		media.Cwd = params.ScriptCwd

		for index, envName := range params.ScriptEnvNames {
			if index < len(params.ScriptEnvValues) {
				media.AddEnv(envName, params.ScriptEnvValues[index])
			}
		}

		err := teautils.ObjectToMapJSON(media, &mediaConfig.Options)
		if err != nil {
			logs.Error(err)
		}
	case notices.NoticeMediaTypeDingTalk:
		params.Must.
			Field("dingTalkWebhookURL", params.DingTalkWebhookURL).
			Require("请输入Hook地址").
			Match("^https:", "Hook地址必须以https://开头")

		media := notices.NewNoticeDingTalkMedia()
		media.WebhookURL = params.DingTalkWebhookURL
		err := teautils.ObjectToMapJSON(media, &mediaConfig.Options)
		if err != nil {
			logs.Error(err)
		}
	case notices.NoticeMediaTypeQyWeixin:
		params.Must.
			Field("qyWeixinCorporateId", params.QyWeixinCorporateId).
			Require("请输入企业ID").
			Field("qyWeixinAgentId", params.QyWeixinAgentId).
			Require("请输入应用AgentId").
			Field("qyWeixinSecret", params.QyWeixinAppSecret).
			Require("请输入应用Secret")

		media := notices.NewNoticeQyWeixinMedia()
		media.CorporateId = params.QyWeixinCorporateId
		media.AgentId = params.QyWeixinAgentId
		media.AppSecret = params.QyWeixinAppSecret
		media.TextFormat = params.QyWeixinTextFormat
		err := teautils.ObjectToMapJSON(media, &mediaConfig.Options)
		if err != nil {
			logs.Error(err)
		}
	case notices.NoticeMediaTypeQyWeixinRobot:
		params.Must.
			Field("qyWeixinRobotWebhookURL", params.QyWeixinRobotWebhookURL).
			Require("请输入Webhook地址").
			Match("^https:", "Webhook地址必须以https://开头")

		media := notices.NewNoticeQyWeixinRobotMedia()
		media.WebhookURL = params.QyWeixinRobotWebhookURL
		media.TextFormat = params.QyWeixinRobotTextFormat
		err := teautils.ObjectToMapJSON(media, &mediaConfig.Options)
		if err != nil {
			logs.Error(err)
		}
	case notices.NoticeMediaTypeAliyunSms:
		params.Must.
			Field("aliyunSmsSign", params.AliyunSmsSign).
			Require("请输入签名名称").
			Field("aliyunSmsTemplateCode", params.AliyunSmsTemplateCode).
			Require("请输入模板CODE").
			Field("aliyunSmsAccessKeyId", params.AliyunSmsAccessKeyId).
			Require("请输入AccessKey ID").
			Field("aliyunSmsAccessKeySecret", params.AliyunSmsAccessKeySecret).
			Require("请输入AccessKey Secret")

		media := notices.NewNoticeAliyunSmsMedia()
		media.Sign = params.AliyunSmsSign
		media.TemplateCode = params.AliyunSmsTemplateCode
		media.AccessKeyId = params.AliyunSmsAccessKeyId
		media.AccessKeySecret = params.AliyunSmsAccessKeySecret

		for index, name := range params.AliyunSmsTemplateVarNames {
			if index < len(params.AliyunSmsTemplateVarValues) {
				media.Variables = append(media.Variables, &shared.Variable{
					Name:  name,
					Value: params.AliyunSmsTemplateVarValues[index],
				})
			}
		}

		err := teautils.ObjectToMapJSON(media, &mediaConfig.Options)
		if err != nil {
			logs.Error(err)
		}
	case notices.NoticeMediaTypeTelegram:
		params.Must.
			Field("telegramToken", params.TelegramToken).
			Require("请输入机器人Token")
		media := notices.NewNoticeTelegramMedia()
		media.Token = params.TelegramToken
		err := teautils.ObjectToMapJSON(media, &mediaConfig.Options)
		if err != nil {
			logs.Error(err)
		}
	case notices.NoticeMediaTypeTeaSms:
		params.Must.
			Field("teaSmsAccessId", params.TeaSmsAccessId).
			Require("请输入AccessId").
			Field("teaSmsAccessSecret", params.TeaSmsAccessSecret).
			Require("请输入AccessSecret")
		media := notices.NewNoticeTeaSmsMedia()
		media.AccessId = params.TeaSmsAccessId
		media.AccessSecret = params.TeaSmsAccessSecret
		err := teautils.ObjectToMapJSON(media, &mediaConfig.Options)
		if err != nil {
			logs.Error(err)
		}
	}

	// 时间
	params.Must.
		Field("timeFromHour", params.TimeFromHour).
		Require("请输入正确的小时数").
		Gte(0, "请输入正确的小时数").
		Lte(23, "请输入正确的小时数").
		Field("timeFromMinute", params.TimeFromMinute).
		Require("请输入正确的分钟数").
		Gte(0, "请输入正确的分钟数").
		Lte(59, "请输入正确的分钟数").
		Field("timeFromSecond", params.TimeFromSecond).
		Require("请输入正确的秒数").
		Gte(0, "请输入正确的秒数").
		Lte(59, "请输入正确的秒数").

		Field("timeToHour", params.TimeToHour).
		Require("请输入正确的小时数").
		Gte(0, "请输入正确的小时数").
		Lte(23, "请输入正确的小时数").
		Field("timeToMinute", params.TimeToMinute).
		Require("请输入正确的分钟数").
		Gte(0, "请输入正确的分钟数").
		Lte(59, "请输入正确的分钟数").
		Field("timeToSecond", params.TimeToSecond).
		Require("请输入正确的秒数").
		Gte(0, "请输入正确的秒数").
		Lte(59, "请输入正确的秒数")

	mediaConfig.TimeFrom = fmt.Sprintf("%02d:%02d:%02d", params.TimeFromHour, params.TimeFromMinute, params.TimeFromSecond)
	mediaConfig.TimeTo = fmt.Sprintf("%02d:%02d:%02d", params.TimeToHour, params.TimeToMinute, params.TimeToSecond)
	mediaConfig.RateCount = params.RateCount
	mediaConfig.RateMinutes = params.RateMinutes

	err := setting.Save()
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	this.Success()
}
