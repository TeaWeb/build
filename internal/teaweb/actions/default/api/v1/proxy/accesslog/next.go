package accesslog

import (
	"github.com/TeaWeb/build/internal/teadb"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/api/apiutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
	timeutil "github.com/iwind/TeaGo/utils/time"
	"time"
)

type NextAction actions.Action

// 读取下一条日志
func (this *NextAction) RunGet(params struct {
	ServerId string
	LastId   string
}) {
	// 查找昨天的
	accessLogs, err := teadb.AccessLogDAO().ListLatestAccessLogs(timeutil.Format("Ymd", time.Now().AddDate(0, 0, -1)), params.ServerId, params.LastId, false, 1)
	if err != nil {
		apiutils.Fail(this, err.Error())
		return
	}

	if len(accessLogs) == 0 {
		// 查找今天的
		accessLogs, err = teadb.AccessLogDAO().ListLatestAccessLogs(timeutil.Format("Ymd"), params.ServerId, params.LastId, false, 1)
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
