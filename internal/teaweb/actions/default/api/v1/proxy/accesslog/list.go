package accesslog

import (
	"github.com/TeaWeb/build/internal/teadb"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/api/apiutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
	timeutil "github.com/iwind/TeaGo/utils/time"
	"time"
)

type ListAction actions.Action

// 列出一组访问日志
func (this *ListAction) RunGet(params struct {
	ServerId string
	Size     int
}) {
	size := params.Size
	if size <= 0 {
		size = 10
	}

	accessLogs, err := teadb.AccessLogDAO().ListLatestAccessLogs(timeutil.Format("Ymd"), params.ServerId, "", false, size)
	if err != nil {
		apiutils.Fail(this, err.Error())
		return
	}

	if len(accessLogs) == 0 {
		// 查找昨天的
		accessLogs, err = teadb.AccessLogDAO().ListLatestAccessLogs(timeutil.Format("Ymd", time.Now().AddDate(0, 0, -1)), params.ServerId, "", false, size)
		if err != nil {
			apiutils.Fail(this, err.Error())
			return
		}
	}

	apiutils.Success(this, maps.Map{
		"accesslogs": accessLogs,
	})
}
