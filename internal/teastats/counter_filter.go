package teastats

import (
	"github.com/TeaWeb/build/internal/teaconfigs/stats"
	"github.com/TeaWeb/build/internal/tealogs/accesslogs"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/timers"
	"github.com/iwind/TeaGo/types"
	"github.com/iwind/TeaGo/utils/time"
	"net"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

// 数值增长型的统计
type CounterFilter struct {
	looper     *timers.Looper
	queue      *Queue
	code       string
	Period     stats.ValuePeriod
	valuesSize int
	values     map[string]*CounterValue // param_time => value
	locker     sync.Mutex

	IncreaseFunc func(value maps.Map, inc maps.Map) maps.Map

	sortedKeys []string
}

// 启动筛选器
func (this *CounterFilter) StartFilter(code string, period stats.ValuePeriod) {
	this.code = code
	this.Period = period
	this.values = map[string]*CounterValue{}
	this.valuesSize = 100 // 缓存中不能超过一定数目，防止一次性提交过多

	// 自动导入
	duration := 1 * time.Second
	switch this.Period {
	case stats.ValuePeriodSecond:
		duration = 1 * time.Second
	case stats.ValuePeriodMinute:
		duration = 30 * time.Second
	case stats.ValuePeriodHour:
		duration = 5 * time.Minute
	case stats.ValuePeriodDay:
		duration = 10 * time.Minute
	case stats.ValuePeriodWeek:
		duration = 15 * time.Minute
	case stats.ValuePeriodMonth:
		duration = 20 * time.Minute
	case stats.ValuePeriodYear:
		duration = 30 * time.Minute
	}
	this.looper = timers.Loop(duration, func(looper *timers.Looper) {
		this.locker.Lock()
		defer this.locker.Unlock()

		this.Commit()
	})
}

// 应用筛选器
func (this *CounterFilter) ApplyFilter(accessLog *accesslogs.AccessLog, params map[string]string, incrValue map[string]interface{}) {
	this.locker.Lock()
	defer this.locker.Unlock()

	key := this.encodeParams(params)
	key.WriteString("@")

	switch this.Period {
	case stats.ValuePeriodSecond:
		key.WriteString(strconv.FormatInt(accessLog.Timestamp, 10))
	case stats.ValuePeriodMinute:
		key.WriteString(strconv.Itoa(int(accessLog.Timestamp / 60)))
	case stats.ValuePeriodHour:
		key.WriteString(strconv.Itoa(int(accessLog.Timestamp / 3600)))
	case stats.ValuePeriodDay:
		t := accessLog.Time()
		key.WriteString(strconv.Itoa(t.Year()))
		key.WriteByte('_')
		key.WriteString(strconv.Itoa(int(t.Month())))
		key.WriteByte('_')
		key.WriteString(strconv.Itoa(t.Day()))
	case stats.ValuePeriodWeek:
		t := accessLog.Time()
		year, week := t.ISOWeek()
		key.WriteString(strconv.Itoa(year))
		key.WriteByte('_')
		key.WriteString(strconv.Itoa(week))
	case stats.ValuePeriodMonth:
		t := accessLog.Time()
		key.WriteString(strconv.Itoa(t.Year()))
		key.WriteByte('_')
		key.WriteString(strconv.Itoa(int(t.Month())))
	case stats.ValuePeriodYear:
		t := accessLog.Time()
		key.WriteString(strconv.Itoa(t.Year()))
	}

	keyString := key.String()

	value, found := this.values[keyString]
	if found {
		for k, v := range incrValue {
			v1, ok := value.Value[k]
			if !ok {
				value.Value[k] = v
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
			}
			value.Value[k] = v1
		}
	} else {
		this.values[keyString] = &CounterValue{
			Timestamp: accessLog.Timestamp,
			Params:    params,
			Value:     incrValue,
		}
	}

	if len(this.values) > this.valuesSize {
		this.Commit()
	}
}

// 停止筛选器
func (this *CounterFilter) StopFilter() {
	if this.looper != nil {
		this.looper.Stop()
		this.looper = nil
	}

	this.locker.Lock()
	defer this.locker.Unlock()
	this.Commit()
}

// 检查新UV
func (this *CounterFilter) CheckNewUV(accessLog *accesslogs.AccessLog, attachKey string) bool {
	// 从cookie中检查UV是否存在
	uid, ok := accessLog.Cookie["TeaUID"]
	if !ok || len(uid) == 0 {
		ip := accessLog.RemoteAddr
		if len(ip) == 0 {
			return false
		}
		host, _, err := net.SplitHostPort(ip)
		if err == nil {
			ip = host
		}
		uid = ip
	}
	l := len(uid)
	if l == 0 {
		return false
	}
	if l > 32 {
		uid = uid[:32]
	}

	key := ""
	life := time.Second
	switch this.Period {
	case stats.ValuePeriodSecond:
		key = this.code + "_" + this.queue.ServerId + "_" + uid + "_" + timeutil.Format("YmdHis", accessLog.Time())
	case stats.ValuePeriodMinute:
		key = this.code + "_" + this.queue.ServerId + "_" + uid + "_" + timeutil.Format("YmdHi", accessLog.Time())
		life = 2 * time.Minute
	case stats.ValuePeriodHour:
		key = this.code + "_" + this.queue.ServerId + "_" + uid + "_" + timeutil.Format("YmdH", accessLog.Time())
		life = 2 * time.Hour
	case stats.ValuePeriodDay:
		key = this.code + "_" + this.queue.ServerId + "_" + uid + "_" + timeutil.Format("Ymd", accessLog.Time())
		life = 2 * 24 * time.Hour
	case stats.ValuePeriodWeek:
		key = this.code + "_" + this.queue.ServerId + "_" + uid + "_" + timeutil.Format("YW", accessLog.Time())
		life = 8 * 24 * time.Hour
	case stats.ValuePeriodMonth:
		key = this.code + "_" + this.queue.ServerId + "_" + uid + "_" + timeutil.Format("Ym", accessLog.Time())
		life = 32 * 24 * time.Hour
	case stats.ValuePeriodYear:
		key = this.code + "_" + this.queue.ServerId + "_" + uid + "_" + timeutil.Format("Y", accessLog.Time())
		life = 370 * 24 * time.Hour
	}

	if len(attachKey) > 0 {
		key += attachKey
	}
	hasKey, err := sharedKV.Has(key)
	if err != nil {
		logs.Error(err)
		return false
	}
	if hasKey {
		return false
	}

	err = sharedKV.Set(key, "1", life)
	if err != nil {
		logs.Error(err)
		return false
	}
	return true
}

// 检查新IP
func (this *CounterFilter) CheckNewIP(accessLog *accesslogs.AccessLog, attachKey string) bool {
	// IP是否存在
	ip := accessLog.RemoteAddr
	if len(ip) == 0 {
		return false
	}
	lastIndex := strings.LastIndex(ip, ":")
	if lastIndex > -1 {
		ip = ip[:lastIndex]
	}
	key := ""
	life := time.Second
	switch this.Period {
	case stats.ValuePeriodSecond:
		key = this.code + "_" + this.queue.ServerId + "_" + ip + "_" + timeutil.Format("YmdHis", accessLog.Time())
	case stats.ValuePeriodMinute:
		key = this.code + "_" + this.queue.ServerId + "_" + ip + "_" + timeutil.Format("YmdHi", accessLog.Time())
		life = 2 * time.Minute
	case stats.ValuePeriodHour:
		key = this.code + "_" + this.queue.ServerId + "_" + ip + "_" + timeutil.Format("YmdH", accessLog.Time())
		life = 2 * time.Hour
	case stats.ValuePeriodDay:
		key = this.code + "_" + this.queue.ServerId + "_" + ip + "_" + timeutil.Format("Ymd", accessLog.Time())
		life = 2 * 24 * time.Hour
	case stats.ValuePeriodWeek:
		key = this.code + "_" + this.queue.ServerId + "_" + ip + "_" + timeutil.Format("YW", accessLog.Time())
		life = 8 * 24 * time.Hour
	case stats.ValuePeriodMonth:
		key = this.code + "_" + this.queue.ServerId + "_" + ip + "_" + timeutil.Format("Ym", accessLog.Time())
		life = 32 * 24 * time.Hour
	case stats.ValuePeriodYear:
		key = this.code + "_" + this.queue.ServerId + "_" + ip + "_" + timeutil.Format("Y", accessLog.Time())
		life = 370 * 24 * time.Hour
	}

	if len(attachKey) > 0 {
		key += attachKey
	}
	hasKey, err := sharedKV.Has(key)
	if err != nil {
		logs.Error(err)
		return false
	}
	if hasKey {
		return false
	}

	err = sharedKV.Set(key, "1", life)
	if err != nil {
		logs.Error(err)
		return false
	}

	return true
}

// 提交
func (this *CounterFilter) Commit() {
	if len(this.values) > 0 {
		for _, v := range this.values {
			if this.IncreaseFunc != nil {
				v.Value["$increase"] = this.IncreaseFunc
			}
			if this.queue != nil {
				this.queue.Add(this.code, time.Unix(v.Timestamp, 0), this.Period, v.Params, v.Value)
			}
		}
		this.values = map[string]*CounterValue{}
	}
}

// 判断参数是否相等
func (this *CounterFilter) equalParams(params1 map[string]string, params2 map[string]string) bool {
	if params1 == nil && params2 == nil {
		return true
	}
	if len(params1) != len(params2) {
		return false
	}
	for k, v := range params1 {
		v1, ok := params2[k]
		if !ok {
			return false
		}
		if v != v1 {
			return false
		}
	}
	return true
}

// 对参数进行编码
func (this *CounterFilter) encodeParams(params map[string]string) *strings.Builder {
	result := &strings.Builder{}

	if len(params) == 0 {
		return result
	}

	if len(this.sortedKeys) == 0 {
		for key := range params {
			this.sortedKeys = append(this.sortedKeys, key)
		}
		sort.Slice(this.sortedKeys, func(i, j int) bool {
			return this.sortedKeys[i] < this.sortedKeys[j]
		})
	}

	for _, k := range this.sortedKeys {
		result.WriteString(k)
		result.WriteByte(':')
		result.WriteString(params[k])
		result.WriteRune('|')
	}
	return result
}
