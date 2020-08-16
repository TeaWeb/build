package mongo

import (
	"github.com/TeaWeb/build/internal/teaconfigs/db"
	"github.com/iwind/TeaGo/actions"
)

type CleanUpdateAction actions.Action

// 自动清理设置
func (this *CleanUpdateAction) Run(params struct{}) {
	config, _ := db.LoadMongoConfig()
	if config != nil {
		this.Data["accessLog"] = config.AccessLog
	} else {
		this.Data["accessLog"] = nil
	}

	this.Show()
}

// 保存自动清理设置
func (this *CleanUpdateAction) RunPost(params struct {
	AccessLogCleanHour int
	AccessLogKeepDays  int
	Must               *actions.Must
}) {
	params.Must.
		Field("accessLogCleanHour", params.AccessLogCleanHour).
		Gte(0, "请输入一个不小于0的数字").
		Lte(23, "请输入一个在0-23之间的数字").
		Field("accessLogKeepDays", params.AccessLogKeepDays).
		Gte(1, "请输入一个大于0的数字")

	config, _ := db.LoadMongoConfig()
	if config == nil {
		config = db.NewMongoConfig()
	}
	config.AccessLog = &db.MongoAccessLogConfig{
		CleanHour: params.AccessLogCleanHour,
		KeepDays:  params.AccessLogKeepDays,
	}
	err := config.Save()
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	this.Success()
}
