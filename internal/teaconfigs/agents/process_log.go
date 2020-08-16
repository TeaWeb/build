package agents

import (
	"github.com/TeaWeb/build/internal/teadb/shared"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	timeutil "github.com/iwind/TeaGo/utils/time"
	"time"
)

// 进程日志
type ProcessLog struct {
	Id         shared.ObjectId `var:"id" bson:"_id" json:"id"` // 数据库存储的ID
	AgentId    string          `bson:"agentId" json:"agentId"`
	TaskId     string          `bson:"taskId" json:"taskId"`
	ProcessId  string          `bson:"processId" json:"processId"`
	ProcessPid int             `bson:"processPid" json:"processPid"`
	EventType  string          `bson:"eventType" json:"eventType"` // start, log, stop
	Data       string          `bson:"data" json:"data"`
	Timestamp  int64           `bson:"timestamp" json:"timestamp"` // unix时间戳，单位为秒
	TimeFormat struct {
		Year   string `bson:"year" json:"year"`
		Month  string `bson:"month" json:"month"`
		Week   string `bson:"week" json:"week"`
		Day    string `bson:"day" json:"day"`
		Hour   string `bson:"hour" json:"hour"`
		Minute string `bson:"minute" json:"minute"`
		Second string `bson:"second" json:"second"`
	} `bson:"timeFormat" json:"timeFormat"`
}

// 设置时间
func (this *ProcessLog) SetTime(t time.Time) {
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
func (this *ProcessLog) SetDBColumns(v maps.Map) {
	id, err := shared.ObjectIdFromHex(v.GetString("_id"))
	if err != nil {
		logs.Error(err)
	} else {
		this.Id = id
	}
	this.AgentId = v.GetString("agentId")
	this.TaskId = v.GetString("taskId")
	this.ProcessId = v.GetString("processId")
	this.ProcessPid = v.GetInt("processPid")
	this.EventType = v.GetString("eventType")
	this.Data = v.GetString("data")
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
func (this *ProcessLog) DBColumns() maps.Map {
	if this.Id.IsZero() {
		this.Id = shared.NewObjectId()
	}
	return maps.Map{
		"_id":               this.Id.Hex(),
		"agentId":           this.AgentId,
		"taskId":            this.TaskId,
		"processId":         this.ProcessId,
		"processPid":        this.ProcessPid,
		"eventType":         this.EventType,
		"data":              this.Data,
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
