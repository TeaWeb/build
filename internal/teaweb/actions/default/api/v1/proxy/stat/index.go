package stat

import (
	"github.com/TeaWeb/build/internal/teadb"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/api/apiutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
	"strings"
)

type IndexAction actions.Action

// 统计
func (this *IndexAction) RunGet(params struct {
	ServerId string
	Item     string

	// 时间
	Year   string
	Month  string
	Week   string
	Day    string
	Hour   string
	Minute string
	Second string

	// 排序
	Sort string
}) {
	query := teadb.NewQuery(teadb.ServerValueDAO().TableName(params.ServerId))
	query.Attr("item", params.Item)

	if len(params.Year) > 0 {
		query.Attr("timeFormat.year", params.Year)
	}
	if len(params.Month) > 0 {
		query.Attr("timeFormat.month", params.Month)
	}
	if len(params.Week) > 0 {
		query.Attr("timeFormat.week", params.Week)
	}
	if len(params.Hour) > 0 {
		query.Attr("timeFormat.hour", params.Hour)
	}
	if len(params.Minute) > 0 {
		query.Attr("timeFormat.minute", params.Minute)
	}
	if len(params.Second) > 0 {
		query.Attr("timeFormat.second", params.Second)
	}

	// 参数 params.*
	for k, v := range this.Request.URL.Query() {
		if len(v) == 0 || !strings.HasPrefix(k, "params.") {
			continue
		}
		query.Attr(k, v[0])
	}

	if len(params.Sort) > 0 {
		if strings.HasPrefix(params.Sort, "-") {
			query.Desc("value." + params.Sort[1:])
		} else {
			query.Asc("value." + params.Sort)
		}
	}

	values, err := teadb.ServerValueDAO().QueryValues(query)
	if err != nil {
		apiutils.Fail(this, err.Error())
		return
	}

	result := []maps.Map{}
	for _, v := range values {
		result = append(result, maps.Map{
			"value":  v.Value,
			"params": v.Params,
			"time":   v.TimeFormat,
		})
	}

	apiutils.Success(this, result)
}
