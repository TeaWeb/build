package notices

import (
	"fmt"
	"github.com/TeaWeb/build/internal/teadb/shared"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"hash/crc32"
	"strings"
	"time"
)

// 通知
type Notice struct {
	Id          shared.ObjectId `bson:"_id" json:"id"`                  // 数据库存储的ID
	Proxy       ProxyCond       `bson:"proxy" json:"proxy"`             // 代理相关参数
	Agent       AgentCond       `bson:"agent" json:"agent"`             // 主机相关参数
	Timestamp   int64           `bson:"timestamp" json:"timestamp"`     // 时间戳
	Message     string          `bson:"message" json:"message"`         // 消息内容
	MessageHash string          `bson:"messageHash" json:"messageHash"` // 消息内容Hash：crc32(message)
	IsRead      bool            `bson:"isRead" json:"isRead"`           // 已读
	IsNotified  bool            `bson:"isNotified" json:"isNotified"`   // 是否发送通知
	Receivers   []string        `bson:"receivers" json:"receivers"`     // 接收人ID列表
}

// Proxy条件
type ProxyCond struct {
	ServerId   string `bson:"serverId" json:"serverId"`
	Websocket  bool   `bson:"websocket" json:"websocket"`
	LocationId string `bson:"locationId" json:"serverId"`
	RewriteId  string `bson:"rewriteId" json:"serverId"`
	BackendId  string `bson:"backendId" json:"serverId"`
	FastcgiId  string `bson:"fastcgiId" json:"serverId"`
	Level      uint8  `bson:"level" json:"level"`
}

// Agent条件
type AgentCond struct {
	AgentId   string `bson:"agentId" json:"agentId"`
	AppId     string `bson:"appId" json:"appId"`
	TaskId    string `bson:"taskId" json:"taskId"`
	ItemId    string `bson:"itemId" json:"itemId"`
	Level     uint8  `bson:"level" json:"level"`
	Threshold string `bson:"threshold" json:"threshold"`
}

// 获取通知对象
func NewNotice() *Notice {
	return &Notice{
		Id: shared.NewObjectId(),
	}
}

// 设置时间
func (this *Notice) SetTime(t time.Time) {
	this.Timestamp = t.Unix()
}

// 计算Hash
func (this *Notice) Hash() {
	this.MessageHash = fmt.Sprintf("%d", crc32.ChecksumIEEE([]byte(this.Message)))
}

// 设置数据库列值
func (this *Notice) SetDBColumns(v maps.Map) {
	id, err := shared.ObjectIdFromHex(v.GetString("_id"))
	if err != nil {
		logs.Error(err)
	} else {
		this.Id = id
	}

	this.Timestamp = v.GetInt64("timestamp")
	this.Message = v.GetString("message")
	this.MessageHash = v.GetString("messageHash")
	this.IsRead = v.GetInt8("isRead") > 0
	this.IsNotified = v.GetInt8("isNotified") > 0

	receivers := v.GetString("receivers")
	if len(receivers) > 0 {
		this.Receivers = strings.Split(receivers, ",")
	}

	this.Proxy.ServerId = v.GetString("proxyServerId")
	this.Proxy.Websocket = v.GetInt8("proxyWebsocket") > 0
	this.Proxy.LocationId = v.GetString("proxyLocationId")
	this.Proxy.RewriteId = v.GetString("proxyRewriteId")
	this.Proxy.BackendId = v.GetString("proxyBackendId")
	this.Proxy.FastcgiId = v.GetString("proxyFastcgiId")
	this.Proxy.Level = v.GetUint8("level")

	this.Agent.AgentId = v.GetString("agentId")
	this.Agent.AppId = v.GetString("agentAppId")
	this.Agent.TaskId = v.GetString("agentTaskId")
	this.Agent.ItemId = v.GetString("agentItemId")
	this.Agent.Threshold = v.GetString("agentThreshold")
	this.Agent.Level = v.GetUint8("level")
}

// 获取数据库列值
func (this *Notice) DBColumns() maps.Map {
	if this.Id.IsZero() {
		this.Id = shared.NewObjectId()
	}
	isRead := 0
	if this.IsRead {
		isRead = 1
	}
	isNotified := 0
	if this.IsNotified {
		isNotified = 1
	}
	proxyWebsocket := 0
	if this.Proxy.Websocket {
		proxyWebsocket = 1
	}
	level := this.Proxy.Level
	if len(this.Agent.AgentId) > 0 {
		level = this.Agent.Level
	}
	return maps.Map{
		"_id":             this.Id.Hex(),
		"timestamp":       this.Timestamp,
		"message":         this.Message,
		"messageHash":     this.MessageHash,
		"isRead":          isRead,
		"isNotified":      isNotified,
		"receivers":       strings.Join(this.Receivers, ","),
		"proxyServerId":   this.Proxy.ServerId,
		"proxyWebsocket":  proxyWebsocket,
		"proxyLocationId": this.Proxy.LocationId,
		"proxyRewriteId":  this.Proxy.RewriteId,
		"proxyBackendId":  this.Proxy.BackendId,
		"proxyFastcgiId":  this.Proxy.FastcgiId,
		"level":           level,
		"agentId":         this.Agent.AgentId,
		"agentAppId":      this.Agent.AppId,
		"agentTaskId":     this.Agent.TaskId,
		"agentItemId":     this.Agent.ItemId,
		"agentThreshold":  this.Agent.Threshold,
	}
}