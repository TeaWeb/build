package waf

import (
	"fmt"
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/TeaWeb/build/internal/teadb"
	"github.com/TeaWeb/build/internal/tealogs/accesslogs"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/utils/time"
	"math"
	"net/http"
	"regexp"
	"time"
)

type DayAction actions.Action

// 某天的日志
func (this *DayAction) Run(params struct {
	WafId    string
	Day      string
	LogType  string
	FromId   string
	Page     int
	Size     int
	SearchIP string
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

	if params.Page < 1 {
		params.Page = 1
	}
	if params.Size < 1 {
		params.Size = 20
	}

	this.Data["searchIP"] = params.SearchIP

	// 检查数据库连接
	this.Data["mongoError"] = ""
	err := teadb.SharedDB().Test()
	mongoAvailable := true
	if err != nil {
		this.Data["mongoError"] = "此功能需要连接数据库"
		mongoAvailable = false
	}

	this.Data["day"] = params.Day
	this.Data["isHistory"] = regexp.MustCompile("^\\d+$").MatchString(params.Day)
	this.Data["logType"] = params.LogType
	this.Data["logs"] = []interface{}{}
	this.Data["fromId"] = params.FromId
	this.Data["hasNext"] = false
	this.Data["page"] = params.Page

	// 日志列表
	if mongoAvailable {
		realDay := ""
		if regexp.MustCompile("^\\d+$").MatchString(params.Day) {
			realDay = params.Day
		} else if params.Day == "today" {
			realDay = timeutil.Format("Ymd")
		} else if params.Day == "yesterday" {
			realDay = timeutil.Format("Ymd", time.Now().Add(-24*time.Hour))
		} else {
			realDay = timeutil.Format("Ymd")
		}

		accessLogList, err := teadb.AccessLogDAO().ListAccessLogsWithWAF(realDay, waf.Id, params.FromId, params.LogType == "errorLog", params.SearchIP, params.Size*(params.Page-1), params.Size)

		if err != nil {
			this.Data["mongoError"] = "数据库查询错误：" + err.Error()
		} else {
			result := lists.Map(accessLogList, func(k int, v interface{}) interface{} {
				accessLog := v.(*accesslogs.AccessLog)
				return map[string]interface{}{
					"id":             accessLog.Id.Hex(),
					"requestTime":    accessLog.RequestTime,
					"request":        accessLog.Request,
					"requestURI":     accessLog.RequestURI,
					"requestMethod":  accessLog.RequestMethod,
					"remoteAddr":     accessLog.RemoteAddr,
					"remotePort":     accessLog.RemotePort,
					"userAgent":      accessLog.UserAgent,
					"host":           accessLog.Host,
					"status":         accessLog.Status,
					"statusMessage":  fmt.Sprintf("%d", accessLog.Status) + " " + http.StatusText(accessLog.Status),
					"timeISO8601":    accessLog.TimeISO8601,
					"timeLocal":      accessLog.TimeLocal,
					"requestScheme":  accessLog.Scheme,
					"proto":          accessLog.Proto,
					"contentType":    accessLog.SentContentType(),
					"bytesSent":      accessLog.BytesSent,
					"backendAddress": accessLog.BackendAddress,
					"fastcgiAddress": accessLog.FastcgiAddress,
					"extend":         accessLog.Extend,
					"referer":        accessLog.Referer,
					"upgrade":        accessLog.GetHeader("Upgrade"),
					"day":            timeutil.Format("Ymd", accessLog.Time()),
					"errors":         accessLog.Errors,
					"attrs":          accessLog.Attrs,
				}
			})

			this.Data["logs"] = result

			if len(result) > 0 {
				if len(params.FromId) == 0 {
					fromId := accessLogList[0].Id.Hex()
					this.Data["fromId"] = fromId
				}

				{
					nextId := accessLogList[len(accessLogList)-1].Id.Hex()
					b, err := teadb.AccessLogDAO().HasNextAccessLogWithWAF(realDay, waf.Id, nextId, params.LogType == "errorLog", params.SearchIP)
					if err != nil {
						logs.Error(err)
					} else {
						this.Data["hasNext"] = b
					}
				}
			}
		}

		// 统计
		stat, err := teadb.AccessLogDAO().GroupWAFRuleGroups(realDay, waf.Id)
		if err != nil {
			logs.Error(err)
			this.Data["stat"] = []maps.Map{}
		} else {
			// 计算百分比
			total := 0
			for _, m := range stat {
				total += m.GetInt("count")
			}
			if total > 0 {
				for _, m := range stat {
					percent := math.Ceil(float64(m.GetInt("count"))*10000/float64(total)) / 100
					m["name"] = m.GetString("name") + " " + fmt.Sprintf("%.2f", percent) + "%"
				}
			}

			this.Data["stat"] = stat
		}
	}

	this.Show()
}
