package log

import (
	"fmt"
	"github.com/TeaWeb/build/internal/teadb"
	"github.com/TeaWeb/build/internal/tealogs/accesslogs"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/utils/time"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type ListAction actions.Action

var whitespaceSplitReg = regexp.MustCompile(`\s+`)

// 获取日志
func (this *ListAction) Run(params struct {
	ServerId string
	FromId   string
	Size     int `default:"10"`

	RemoteAddr string  // 终端地址
	Domain     string  // 域名
	OsName     string  // 终端OS
	Browser    string  // 终端浏览器
	Cost       float64 // 耗时
	Keyword    string  // 关键词
	BackendId  string  // 后端服务器
	LocationId string  // 路径规则
	RewriteId  string  // 重写规则
	FastcgiId  string  // FastcgiId

	BodyFetching bool
	LogType      string
}) {
	if params.Size < 1 {
		params.Size = 20
	}

	params.Domain = strings.ToLower(params.Domain)
	params.OsName = strings.ToLower(params.OsName)
	params.Browser = strings.ToLower(params.Browser)
	params.Keyword = strings.ToLower(params.Keyword)

	serverId := params.ServerId

	requestBodyFetching = params.BodyFetching
	requestBodyTime = time.Now()

	accessLogs, err := teadb.AccessLogDAO().ListLatestAccessLogs(timeutil.Format("Ymd"), serverId, params.FromId, params.LogType == "errorLog", params.Size)

	this.Data["lastId"] = ""
	if err != nil {
		if err != teadb.ErrorDBUnavailable {
			logs.Error(err)
		}
		this.Data["logs"] = []interface{}{}
	} else {
		result := []maps.Map{}
		if len(accessLogs) > 0 {
			this.Data["lastId"] = accessLogs[0].Id.Hex()
		}
		for _, accessLog := range accessLogs {
			if accessLog.Extend == nil {
				accessLog.Extend = new(accesslogs.AccessLogExtend)
			}

			// filters
			if len(params.RemoteAddr) > 0 && !this.match(accessLog.RemoteAddr, params.RemoteAddr) {
				continue
			}

			if len(params.Domain) > 0 && !this.match(accessLog.Host, params.Domain) {
				continue
			}

			if len(params.OsName) > 0 && !this.match(accessLog.Extend.Client.OS.Family+" "+accessLog.Extend.Client.OS.Major, params.OsName) {
				continue
			}

			if len(params.Browser) > 0 && !this.match(accessLog.Extend.Client.Browser.Family+" "+accessLog.Extend.Client.Browser.Major, params.Browser) {
				continue
			}

			if params.Cost > 0 && accessLog.RequestTime*1000 < params.Cost {
				continue
			}

			if len(params.Keyword) > 0 &&
				!this.match(accessLog.Request, params.Keyword) &&
				!this.match(accessLog.Host, params.Keyword) &&
				!this.match(accessLog.RemoteAddr, params.Keyword) &&
				!this.match(accessLog.UserAgent, params.Keyword) &&
				!this.match(accessLog.Extend.Client.OS.Family+" "+accessLog.Extend.Client.OS.Major, params.Keyword) &&
				!this.match(accessLog.Extend.Client.Browser.Family+" "+accessLog.Extend.Client.Browser.Major, params.Keyword) &&
				!this.match(fmt.Sprintf("%d", accessLog.Status), params.Keyword) &&
				!this.match(accessLog.StatusMessage, params.Keyword) &&
				!this.match(accessLog.ContentType, params.Keyword) &&
				!this.match(accessLog.TimeLocal, params.Keyword) &&
				!this.match(accessLog.TimeISO8601, params.Keyword) {
				continue
			}

			if len(params.BackendId) > 0 && accessLog.BackendId != params.BackendId {
				continue
			}

			if len(params.LocationId) > 0 && accessLog.LocationId != params.LocationId {
				continue
			}

			if len(params.RewriteId) > 0 && accessLog.RewriteId != params.RewriteId {
				continue
			}

			if len(params.FastcgiId) > 0 && accessLog.FastcgiId != params.FastcgiId {
				continue
			}

			result = append(result, map[string]interface{}{
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
				"backendId":      accessLog.BackendId,
				"locationId":     accessLog.LocationId,
				"rewriteId":      accessLog.RewriteId,
				"fastcgiId":      accessLog.FastcgiId,
				"attrs":          accessLog.Attrs,
			})
		}

		this.Data["logs"] = result
	}

	this.Success()
}

func (this *ListAction) match(s string, keyword string) bool {
	if len(keyword) == 0 {
		return false
	}
	if len(s) == 0 {
		return false
	}

	s = strings.ToLower(s)
	ok := true
	for _, piece := range whitespaceSplitReg.Split(keyword, -1) {
		if strings.Index(s, piece) == -1 {
			ok = false
			break
		}
	}
	return ok
}
