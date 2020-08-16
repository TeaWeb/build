package teastats

import (
	"github.com/TeaWeb/build/internal/teaconfigs/stats"
	"github.com/TeaWeb/build/internal/teadb"
	"github.com/TeaWeb/build/internal/teadb/shared"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/timers"
	"github.com/iwind/TeaGo/types"
	"reflect"
	"sync"
	"time"
)

// 入库队列
type Queue struct {
	ServerId    string
	c           chan *stats.Value
	looper      *timers.Looper
	indexes     [][]string // { { page }, { region, province, city }, ... }
	indexLocker sync.Mutex
}

// 获取新对象
func NewQueue() *Queue {
	return &Queue{}
}

// 队列单例
func (this *Queue) Start(serverId string) {
	this.ServerId = serverId
	this.c = make(chan *stats.Value, 4096)

	// 测试连接，如果有错误则重新连接
	if teadb.SharedDB().IsAvailable() {
		err := teadb.SharedDB().Test()
		if err != nil {
			if teadb.SharedDB().IsAvailable() {
				logs.Println("[stat]queue start failed: can not connect to database, will reconnect to database")
			}
			time.Sleep(5 * time.Second)
			this.Start(serverId)
			return
		}
	} else {
		time.Sleep(5 * time.Second)
		this.Start(serverId)
		return
	}

	// 导入数据
	go func() {
		// 延时等待数据库准备好
		time.Sleep(5 * time.Second)

		for {
			item := <-this.c

			if item == nil {
				break
			}
			//logs.Println("[stat]dump item '" + item.Item + "' for server '" + this.ServerId + "'")

			// 是否已存在
			oneValue, err := teadb.ServerValueDAO().FindSameItemValue(serverId, item)
			if err != nil {
				logs.Error(err)
				continue
			}
			if oneValue == nil {
				// 是否有自定义运算函数
				increaseFunc, found := item.Value["$increase"]
				if found {
					item.Value = increaseFunc.(func(value maps.Map, inc maps.Map) maps.Map)(nil, item.Value)
					delete(item.Value, "$increase")
				}
				err := teadb.ServerValueDAO().InsertOne(serverId, item)
				if err != nil {
					logs.Error(err)
				}
			} else {
				// 是否有自定义运算函数
				increaseFunc, found := item.Value["$increase"]
				if found {
					item.Value = increaseFunc.(func(value maps.Map, inc maps.Map) maps.Map)(oneValue.Value, item.Value)
					delete(item.Value, "$increase")
				} else {
					// 简单的增长
					item.Value = this.increase(oneValue.Value, item.Value)
				}
				err = teadb.ServerValueDAO().UpdateItemValueAndTimestamp(serverId, oneValue.Id.Hex(), item.Value, item.Timestamp)
				if err != nil {
					logs.Error(err)
				}
			}
		}
	}()

	// 清理数据
	go func() {
		this.looper = timers.Loop(1*time.Hour, func(looper *timers.Looper) {
			// 清除24小时之前的second
			err := teadb.ServerValueDAO().DeleteExpiredValues(serverId, "second", 24*3600)
			if err != nil {
				logs.Error(err)
			}

			// 清除24小时之前的minute
			err = teadb.ServerValueDAO().DeleteExpiredValues(serverId, "minute", 24*3600)
			if err != nil {
				logs.Error(err)
			}

			// 清除48小时之前的hour
			err = teadb.ServerValueDAO().DeleteExpiredValues(serverId, "hour", 48*3600)
			if err != nil {
				logs.Error(err)
			}
		})
	}()
}

// 添加指标值
func (this *Queue) Add(itemCode string, t time.Time, period stats.ValuePeriod, params map[string]string, value maps.Map) {
	if params == nil {
		params = map[string]string{}
	}
	if value == nil {
		value = maps.Map{}
	}
	item := stats.NewItemValue()
	item.Id = shared.NewObjectId()
	item.Item = itemCode
	item.Period = period
	item.Value = value
	item.Params = params
	item.SetTime(t)

	if this.c != nil {
		this.c <- item
	}
}

// 停止
func (this *Queue) Stop() {
	// 等待数据完成
	if len(this.c) > 0 {
		time.Sleep(200 * time.Millisecond)
	}

	close(this.c)
	this.c = nil

	if this.looper != nil {
		this.looper.Stop()
		this.looper = nil
	}
}

// 添加索引
func (this *Queue) Index(index []string) {
	if len(index) == 0 {
		return
	}

	this.indexLocker.Lock()
	defer this.indexLocker.Unlock()

	// 是否已存在
	for _, i := range this.indexes {
		if this.equalStrings(index, i) {
			return
		}
	}

	fields := []*shared.IndexField{shared.NewIndexField("item", true)}
	for _, i := range index {
		fields = append(fields, shared.NewIndexField("params."+i, true))
	}
	err := teadb.ServerValueDAO().CreateIndex(this.ServerId, fields)
	if err != nil {
		logs.Error(err)
	}

	this.indexes = append(this.indexes, index)
}

// 增加值
// 只支持int, int32, int64, float32, float64
func (this *Queue) increase(value maps.Map, inc maps.Map) maps.Map {
	if inc == nil {
		return maps.Map{}
	}
	if value == nil {
		return inc
	}
	for k, v := range inc {
		v1, ok := value[k]
		if !ok {
			value[k] = v
			continue
		}
		switch v2 := v1.(type) {
		case int:
			v1 = v2 + types.Int(v)
		case int32:
			v1 = v2 + types.Int32(v)
		case int64:
			v1 = v2 + types.Int64(v)
		case float32:
			v1 = v2 + types.Float32(v)
		case float64:
			v1 = v2 + types.Float64(v)
		default:
			logs.Println("[teastats]queue increase not match:", reflect.TypeOf(v1).Kind())
		}
		value[k] = v1
	}

	return value
}

// 对比字符串数组看是否相等
func (this *Queue) equalStrings(strings1 []string, strings2 []string) bool {
	for _, s1 := range strings1 {
		if !lists.ContainsString(strings2, s1) {
			return false
		}
	}
	for _, s2 := range strings2 {
		if !lists.ContainsString(strings1, s2) {
			return false
		}
	}
	return true
}
