package stats

import (
	"encoding/json"
	"github.com/TeaWeb/build/internal/teadb/shared"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/utils/time"
	"github.com/pquerna/ffjson/ffjson"
	"time"
)

// 值周期
type ValuePeriod = string

const (
	ValuePeriodSecond = "second"
	ValuePeriodMinute = "minute"
	ValuePeriodHour   = "hour"
	ValuePeriodDay    = "day"
	ValuePeriodWeek   = "week"
	ValuePeriodMonth  = "month"
	ValuePeriodYear   = "year"
)

// 统计指标值定义
type Value struct {
	Id         shared.ObjectId        `bson:"_id" json:"id"`              // 数据库存储的ID
	Item       string                 `bson:"item" json:"item"`           // 指标代号
	Period     ValuePeriod            `bson:"period" json:"period"`       // 周期
	Value      map[string]interface{} `bson:"value" json:"value"`         // 数据内容
	Params     map[string]string      `bson:"params" json:"params"`       // 参数
	Timestamp  int64                  `bson:"timestamp" json:"timestamp"` // 时间戳
	TimeFormat struct {
		Year   string `bson:"year" json:"year"`
		Month  string `bson:"month" json:"month"`
		Week   string `bson:"week" json:"week"`
		Day    string `bson:"day" json:"day"`
		Hour   string `bson:"hour" json:"hour"`
		Minute string `bson:"minute" json:"minute"`
		Second string `bson:"second" json:"second"`
	} `bson:"timeFormat" json:"timeFormat"`                               // 时间信息
}

// 获取新对象
func NewItemValue() *Value {
	return &Value{
		Params: map[string]string{},
	}
}

func (this *Value) SetTime(t time.Time) {
	this.Timestamp = t.Unix()
	this.TimeFormat.Year = timeutil.Format("Y", t)
	this.TimeFormat.Month = timeutil.Format("Ym", t)
	this.TimeFormat.Week = timeutil.Format("YW", t)
	this.TimeFormat.Day = timeutil.Format("Ymd", t)
	this.TimeFormat.Hour = timeutil.Format("YmdH", t)
	this.TimeFormat.Minute = timeutil.Format("YmdHi", t)
	this.TimeFormat.Second = timeutil.Format("YmdHis", t)
}

// 设置数据库列值
func (this *Value) SetDBColumns(v maps.Map) {
	id, err := shared.ObjectIdFromHex(v.GetString("_id"))
	if err != nil {
		logs.Error(err)
	} else {
		this.Id = id
	}
	this.Item = v.GetString("item")
	this.Period = v.GetString("period")
	this.jsonDecode(v.Get("value"), &this.Value)
	this.jsonDecode(v.Get("params"), &this.Params)
	this.Timestamp = v.GetInt64("timestamp")
	this.TimeFormat.Year = v.GetString("timeFormat_year")
	this.TimeFormat.Month = v.GetString("timeFormat_month")
	this.TimeFormat.Week = v.GetString("timeFormat_week")
	this.TimeFormat.Day = v.GetString("timeFormat_day")
	this.TimeFormat.Hour = v.GetString("timeFormat_hour")
	this.TimeFormat.Minute = v.GetString("timeFormat_minute")
	this.TimeFormat.Second = v.GetString("timeFormat_second")
}

// 获取数据库列值
func (this *Value) DBColumns() maps.Map {
	if this.Id.IsZero() {
		this.Id = shared.NewObjectId()
	}
	valueJSON, err := json.Marshal(this.Value)
	if err != nil {
		logs.Error(err)
	}
	paramsJSON, err := json.Marshal(this.Params)
	if err != nil {
		logs.Error(err)
	}
	return maps.Map{
		"_id":               this.Id.Hex(),
		"item":              this.Item,
		"period":            this.Period,
		"value":             valueJSON,
		"params":            paramsJSON,
		"timestamp":         this.Timestamp,
		"timeFormat_year":   this.TimeFormat.Year,
		"timeFormat_month":  this.TimeFormat.Month,
		"timeFormat_week":   this.TimeFormat.Week,
		"timeFormat_day":    this.TimeFormat.Day,
		"timeFormat_hour":   this.TimeFormat.Hour,
		"timeFormat_minute": this.TimeFormat.Minute,
		"timeFormat_second": this.TimeFormat.Second,
	}
}

func (this *Value) jsonDecode(data interface{}, vPtr interface{}) {
	if data == nil {
		return
	}
	b, ok := data.([]byte)
	if ok {
		_ = ffjson.Unmarshal(b, vPtr)
	}
	s, ok := data.(string)
	if ok {
		_ = ffjson.Unmarshal([]byte(s), vPtr)
	}
}
