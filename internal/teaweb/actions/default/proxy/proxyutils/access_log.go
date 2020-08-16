package proxyutils

import (
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/TeaWeb/build/internal/tealogs/accesslogs"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/types"
	"net/http"
)

// 格式化访问日志配置
func FormatAccessLog(accessLogs []*teaconfigs.AccessLogConfig) []maps.Map {
	// access log
	result := []maps.Map{}
	for _, accessLog := range accessLogs {
		m := maps.Map{
			"id":          accessLog.Id,
			"on":          accessLog.On,
			"status1":     accessLog.Status1,
			"status2":     accessLog.Status2,
			"status3":     accessLog.Status3,
			"status4":     accessLog.Status4,
			"status5":     accessLog.Status5,
			"storageOnly": accessLog.StorageOnly,
		}

		// fields
		fields := []maps.Map{}
		for _, field := range accesslogs.AccessLogFields {
			code := types.Int(field["code"])
			name := field["name"]
			isChecked := false
			if len(accessLog.Fields) == 0 {
				isChecked = lists.ContainsInt(accesslogs.AccessLogDefaultFieldsCodes, code)
			} else {
				isChecked = lists.ContainsInt(accessLog.Fields, code)
			}
			fields = append(fields, maps.Map{
				"name":      name,
				"code":      code,
				"isChecked": isChecked,
			})
		}
		m["fields"] = fields

		// 存储策略
		m["storagePolicies"] = lists.Map(teaconfigs.SharedAccessLogStoragePolicyList().FindAllPolicies(), func(k int, v interface{}) interface{} {
			policy := v.(*teaconfigs.AccessLogStoragePolicy)
			return maps.Map{
				"id":        policy.Id,
				"name":      policy.Name,
				"type":      policy.Type,
				"isChecked": accessLog.ContainsStoragePolicy(policy.Id),
			}
		})
		m["hasSelectedStoragePolicies"] = len(accessLog.StoragePolicies) > 0

		result = append(result, m)
	}

	return result
}

// 从请求中获取访问日志信息
func ParseAccessLogForm(req *http.Request) (result []*teaconfigs.AccessLogConfig) {
	indexes, ok := req.Form["accessLogIndexes"]
	if !ok {
		return nil
	}
	if len(indexes) == 0 {
		return nil
	}
	for _, index := range indexes {
		id := req.FormValue("accessLog" + index + "Id")
		on := req.FormValue("accessLog" + index + "On")
		fields, _ := req.Form["accessLog"+index+"Fields"]
		status1 := req.FormValue("accessLog" + index + "Status1")
		status2 := req.FormValue("accessLog" + index + "Status2")
		status3 := req.FormValue("accessLog" + index + "Status3")
		status4 := req.FormValue("accessLog" + index + "Status4")
		status5 := req.FormValue("accessLog" + index + "Status5")
		storagePolicyIds, _ := req.Form["accessLog"+index+"StoragePolicyIds"]
		storageOnly := req.FormValue("accessLog" + index + "StorageOnly")

		accessLog := teaconfigs.NewAccessLogConfig()
		accessLog.Id = id
		accessLog.On = on == "1"

		if len(fields) == 0 {
			accessLog.Fields = []int{}
		} else {
			for _, field := range fields {
				fieldInt := types.Int(field)
				if fieldInt > 0 {
					accessLog.Fields = append(accessLog.Fields, fieldInt)
				}
			}
		}
		if len(accessLog.Fields) == 0 {
			accessLog.Fields = []int{0}
		}

		accessLog.Status1 = status1 == "1"
		accessLog.Status2 = status2 == "1"
		accessLog.Status3 = status3 == "1"
		accessLog.Status4 = status4 == "1"
		accessLog.Status5 = status5 == "1"

		accessLog.StoragePolicies = storagePolicyIds
		accessLog.StorageOnly = storageOnly == "1"

		result = append(result, accessLog)
	}

	return result
}
