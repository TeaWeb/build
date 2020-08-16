package teastats

import (
	"github.com/TeaWeb/build/internal/tealogs/accesslogs"
	stringutil "github.com/iwind/TeaGo/utils/string"
	"testing"
	"time"
)

func TestCostAllPeriodFilter_Start(t *testing.T) {
	queue := new(Queue)
	queue.Start("123456")

	filter := new(CostPagePeriodFilter)
	filter.Start(queue, "cost.all.hour")

	{
		accessLog := &accesslogs.AccessLog{}
		accessLog.RequestPath = "/"
		accessLog.Timestamp = time.Now().Unix()
		accessLog.RequestTime = 0.01
		filter.Filter(accessLog)
	}

	{
		accessLog := &accesslogs.AccessLog{}
		accessLog.RequestPath = "/"
		accessLog.Timestamp = time.Now().Unix()
		accessLog.RequestTime = 0.02
		filter.Filter(accessLog)
	}

	{
		accessLog := &accesslogs.AccessLog{}
		accessLog.RequestPath = "/"
		accessLog.Timestamp = time.Now().Unix()
		accessLog.RequestTime = 0.01
		filter.Filter(accessLog)
	}

	t.Log(filter.values)
	t.Log(stringutil.JSONEncodePretty(filter.values))

	time.Sleep(1 * time.Second)
}
