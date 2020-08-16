package accesslog

import (
	"github.com/TeaWeb/build/internal/teadb"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/api/apiutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
	timeutil "github.com/iwind/TeaGo/utils/time"
	"time"
)

type LatestAction actions.Action

// 读取最后一条访问日志
func (this *LatestAction) RunGet(params struct {
	ServerId string
}) {
	accessLogs, err := teadb.AccessLogDAO().ListLatestAccessLogs(timeutil.Format("Ymd"), params.ServerId, "", false, 1)
	if err != nil {
		apiutils.Fail(this, err.Error())
		return
	}

	if len(accessLogs) == 0 {
		// 查找昨天的
		accessLogs, err = teadb.AccessLogDAO().ListLatestAccessLogs(timeutil.Format("Ymd", time.Now().AddDate(0, 0, -1)), params.ServerId, "", false, 1)
		if err != nil {
			apiutils.Fail(this, err.Error())
			return
		}

		if len(accessLogs) == 0 {
			apiutils.Success(this, maps.Map{
				"accesslog": nil,
			})
			return
		}
	}

	apiutils.Success(this, maps.Map{
		"accesslog": accessLogs[0],
	})
}
