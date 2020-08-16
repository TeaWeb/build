package log

import (
	"github.com/TeaWeb/build/internal/teaconfigs/audits"
	"github.com/TeaWeb/build/internal/teadb"
	"github.com/TeaWeb/build/internal/teageo"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/utils/time"
	"math"
	"net"
	"strings"
	"time"
)

type AuditAction actions.Action

// 审计日志
func (this *AuditAction) Run(params struct {
	Read int
	Page int
}) {
	// 分页
	if params.Page < 1 {
		params.Page = 1
	}
	pageSize := 10
	this.Data["page"] = params.Page

	count, err := teadb.AuditLogDAO().CountAllAuditLogs()
	if err != nil {
		logs.Error(err)
	}

	if count > 0 {
		this.Data["countPages"] = int(math.Ceil(float64(count) / float64(pageSize)))
	} else {
		this.Data["countPages"] = 0
	}

	// 读取列表数据
	ones, err := teadb.AuditLogDAO().ListAuditLogs(pageSize*(params.Page-1), pageSize)
	if err != nil {
		this.Data["logs"] = []interface{}{}
	} else {
		this.Data["logs"] = lists.Map(ones, func(k int, v interface{}) interface{} {
			log := v.(*audits.Log)

			ip, ok := log.Options["ip"]
			location := ""
			if ok && len(ip) > 0 {
				if ip == "127.0.0.1" ||
					strings.HasPrefix(ip, "192.168.") ||
					strings.HasPrefix(ip, "10.") ||
					strings.HasPrefix(ip, "172.16.") {
					location = ""
				} else {
					ipObj := net.ParseIP(ip)
					if ipObj != nil {
						record, err := teageo.DB.City(ipObj)
						if err == nil {
							if _, ok := record.Country.Names["zh-CN"]; ok {
								location = teageo.ConvertName(record.Country.Names["zh-CN"])
							}
							if _, ok := record.City.Names["zh-CN"]; ok {
								location += " " + teageo.ConvertName(record.City.Names["zh-CN"])
							}
						}
					}
				}
			}

			return maps.Map{
				"username":    log.Username,
				"action":      log.Action,
				"actionName":  log.ActionName(),
				"description": log.Description,
				"datetime":    timeutil.Format("Y-m-d H:i:s", time.Unix(log.Timestamp, 0)),
				"options":     log.Options,
				"location":    location,
			}
		})
	}

	this.Show()
}
