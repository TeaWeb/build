package waf

import (
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/TeaWeb/build/internal/teadb"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	timeutil "github.com/iwind/TeaGo/utils/time"
	"time"
)

type HistoryAction actions.Action

func (this *HistoryAction) RunGet(params struct {
	WafId string
}) {
	waf := teaconfigs.SharedWAFList().FindWAF(params.WafId)
	if waf == nil {
		this.Fail("找不到WAF")
	}

	this.Data["config"] = maps.Map{
		"id":            waf.Id,
		"name":          waf.Name,
		"countInbound":  waf.CountInboundRuleSets(),
		"countOutbound": waf.CountOutboundRuleSets(),
		"on":            waf.On,
		"actionBlock":   waf.ActionBlock,
		"cond":          waf.Cond,
	}

	// 检查MongoDB连接
	this.Data["mongoError"] = ""
	err := teadb.SharedDB().Test()
	mongoAvailable := true
	if err != nil {
		this.Data["mongoError"] = "此功能需要连接数据库"
		mongoAvailable = false
	}

	// 列出最近30天的日志
	days := []maps.Map{}
	if mongoAvailable {
		for i := 0; i < 60; i++ {
			day := timeutil.Format("Ymd", time.Now().Add(time.Duration(-i*24)*time.Hour))

			b, err := teadb.AccessLogDAO().HasAccessLogWithWAF(day, waf.Id)
			if err != nil {
				logs.Error(err)
			}
			if b {
				days = append(days, maps.Map{
					"day": day,
					"has": true,
				})
			} else {
				days = append(days, maps.Map{
					"day": day,
					"has": false,
				})
			}
		}
	}

	this.Data["days"] = days
	this.Data["today"] = timeutil.Format("Ymd")

	this.Show()
}
