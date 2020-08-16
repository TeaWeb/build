package audits

import (
	"encoding/json"
	"github.com/TeaWeb/build/internal/teadb/shared"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"time"
)

// 动作
type Action = string

const (
	ActionLogin = "LOGIN" // 登录 {ip}
)

// 审计日志
type Log struct {
	Id          shared.ObjectId   `bson:"_id" json:"id"` // 数据库存储的ID
	Username    string            `bson:"username" json:"username"`
	Action      Action            `bson:"action" json:"action"`           // 类型
	Description string            `bson:"description" json:"description"` // 描述
	Options     map[string]string `bson:"options" json:"options"`         // 选项
	Timestamp   int64             `bson:"timestamp" json:"timestamp"`     // 时间戳
}

// 获取新审计日志对象
func NewLog(username string, action Action, description string, options map[string]string) *Log {
	return &Log{
		Id:          shared.NewObjectId(),
		Username:    username,
		Action:      action,
		Description: description,
		Timestamp:   time.Now().Unix(),
		Options:     options,
	}
}

// 设置数据库列值
func (this *Log) SetDBColumns(v maps.Map) {
	this.Action = v.GetString("action")
	id, err := shared.ObjectIdFromHex(v.GetString("_id"))
	if err != nil {
		logs.Error(err)
	} else {
		this.Id = id
	}
	this.Username = v.GetString("username")
	this.Description = v.GetString("description")
	this.Timestamp = v.GetInt64("timestamp")
	this.Options = map[string]string{}

	options := v.GetString("options")
	if len(options) > 0 {
		err = json.Unmarshal([]byte(v.GetString("options")), &this.Options)
		if err != nil {
			logs.Error(err)
		}
	}
}

// 获取数据库列值
func (this *Log) DBColumns() maps.Map {
	if this.Id.IsZero() {
		this.Id = shared.NewObjectId()
	}
	optionsData, _ := json.Marshal(this.Options)
	return maps.Map{
		"_id":         this.Id.Hex(),
		"action":      this.Action,
		"username":    this.Username,
		"description": this.Description,
		"timestamp":   this.Timestamp,
		"options":     optionsData,
	}
}

// 审计日志类型
func (this *Log) ActionName() string {
	switch this.Action {
	case ActionLogin:
		return "登录"
	}
	return ""
}
