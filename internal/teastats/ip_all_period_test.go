package teastats

import (
	"github.com/TeaWeb/build/internal/tealogs/accesslogs"
	stringutil "github.com/iwind/TeaGo/utils/string"
	"testing"
	"time"
)

func TestIPAllPeriodFilter_Start(t *testing.T) {
	if sharedKV == nil {
		sharedKV = NewKVStorage("stat.leveldb")
		if sharedKV == nil {
			t.Fatal("sharedKV = nil")
		}
	}

	queue := new(Queue)
	queue.Start("123456")

	filter := new(IPAllPeriodFilter)
	filter.Start(queue, "ip.all.minute")
	filter.valuesSize = 0

	accessLog := &accesslogs.AccessLog{}
	accessLog.Timestamp = time.Now().Unix()
	accessLog.RemoteAddr = "133.18.203.152:1234"
	filter.Filter(accessLog)

	t.Log(filter.values)
	t.Log(stringutil.JSONEncodePretty(filter.values))

	time.Sleep(1 * time.Second)
	queue.Stop()
}
