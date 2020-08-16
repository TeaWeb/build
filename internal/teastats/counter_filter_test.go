package teastats

import (
	"github.com/TeaWeb/build/internal/teaconfigs/stats"
	"github.com/TeaWeb/build/internal/tealogs/accesslogs"
	"github.com/iwind/TeaGo/logs"
	"runtime"
	"testing"
	"time"
)

func TestCounterFilter_EncodeParams(t *testing.T) {
	filter := &CounterFilter{}
	t.Log(filter.encodeParams(map[string]string{
		"a": "1",
		"b": "2",
		"c": "3",
	}))
	t.Log(filter.encodeParams(map[string]string{
		"a": "1",
		"b": "2",
		"c": "3",
	}))
	t.Log(filter.encodeParams(map[string]string{
		"a": "1",
		"b": "2",
		"c": "3",
	}))
}

func TestCounterFilter_EncodeParams2(t *testing.T) {
	filter := &CounterFilter{}
	t.Log(filter.encodeParams(map[string]string{}))
}

func TestCounterFilter_ApplyFilter(t *testing.T) {
	runtime.GOMAXPROCS(1)

	filter := &CounterFilter{}
	filter.StartFilter("test", stats.ValuePeriodSecond)
	time.Sleep(1 * time.Second)
	before := time.Now()
	count := 10000
	for i := 0; i < count; i++ {
		filter.ApplyFilter(&accesslogs.AccessLog{
			Timestamp: before.Unix(),
		}, map[string]string{
			"a": "1",
			"b": "2",
		}, map[string]interface{}{
			"count": 1,
		})
	}
	t.Log(time.Since(before).Seconds()*1000, "ms")
	logs.PrintAsJSON(filter.values, t)
}

func BenchmarkCounterFilter_ApplyFilter(b *testing.B) {
	runtime.GOMAXPROCS(1)

	filter := &CounterFilter{}
	filter.StartFilter("test", stats.ValuePeriodSecond)
	now := time.Now()
	accessLog := &accesslogs.AccessLog{
		Timestamp: now.Unix(),
	}

	for i := 0; i < b.N; i++ {
		filter.ApplyFilter(accessLog, map[string]string{
			"a": "1",
			"b": "2",
		}, map[string]interface{}{
			"count": 1,
		})
	}
}

func BenchmarkCounterFilter_EncodeParams(b *testing.B) {
	runtime.GOMAXPROCS(1)
	filter := &CounterFilter{}
	for i := 0; i < b.N; i++ {
		filter.encodeParams(map[string]string{
			"a": "1",
			"b": "2",
			"c": "3",
		})
	}
}

func BenchmarkCounterFilter_EncodeParams2(b *testing.B) {
	runtime.GOMAXPROCS(1)
	filter := &CounterFilter{}
	for i := 0; i < b.N; i++ {
		filter.encodeParams(map[string]string{})
	}
}
