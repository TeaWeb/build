package accesslog

import (
	"github.com/TeaWeb/build/internal/teadb"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/api/apiutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
	timeutil "github.com/iwind/TeaGo/utils/time"
	"time"
)

type NextListAction actions.Action

// 列表下N条数据
func (this *NextListAction) RunGet(params struct {
	ServerId string
	Size     int
	LastId   string
}) {
	size := params.Size
	if size <= 0 {
		size = 10
	}

	// 查找昨天的
	accessLogs, err := teadb.AccessLogDAO().ListLatestAccessLogs(timeutil.Format("Ymd", time.Now().AddDate(0, 0, -1)), params.ServerId, params.LastId, false, size)
	if err != nil {
		apiutils.Fail(this, err.Error())
		return
	}

	if len(accessLogs) == 0 {
		// 查找今天的
		accessLogs, err = teadb.AccessLogDAO().ListLatestAccessLogs(timeutil.Format("Ymd"), params.ServerId, params.LastId, false, size)
		if err != nil {
			apiutils.Fail(this, err.Error())
			return
		}
	}

	apiutils.Success(this, maps.Map{
		"accesslogs": accessLogs,
	})
}
