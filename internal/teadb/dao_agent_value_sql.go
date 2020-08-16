package teadb

import (
	"errors"
	"github.com/TeaWeb/build/internal/teaconfigs/agents"
	"github.com/TeaWeb/build/internal/teaconfigs/notices"
	"github.com/iwind/TeaGo/logs"
)

type SQLAgentValueDAO struct {
	BaseDAO
}

// 初始化
func (this *SQLAgentValueDAO) Init() {

}

// 获取表格
func (this *SQLAgentValueDAO) TableName(agentId string) string {
	table := "teaweb_values_agent_" + agentId
	this.initTable(table)
	return table
}

// 插入数据
func (this *SQLAgentValueDAO) Insert(agentId string, value *agents.Value) error {
	if value == nil {
		return errors.New("value should not be nil")
	}
	if len(agentId) == 0 {
		if len(value.AgentId) > 0 {
			agentId = value.AgentId
		} else {
			return errors.New("AgentId should be set")
		}
	}

	if value.Value == nil {
		value.Value = 0
	}

	return NewQuery(this.TableName(agentId)).InsertOne(value)
}

// 清除数值
func (this *SQLAgentValueDAO) ClearItemValues(agentId string, appId string, itemId string, level notices.NoticeLevel) error {
	if len(agentId) == 0 {
		return errors.New("agentId should not be empty")
	}
	query := NewQuery(this.TableName(agentId)).
		Attr("appId", appId).
		Attr("itemId", itemId)
	if level > 0 {
		query.Attr("noticeLevel", level)
	}
	return query.Delete()
}

// 查找最近的一条记录
func (this *SQLAgentValueDAO) FindLatestItemValue(agentId string, appId string, itemId string) (*agents.Value, error) {
	query := NewQuery(this.TableName(agentId)).
		Attr("itemId", itemId).
		Node().
		Desc("createdAt")
	if len(appId) > 0 {
		query.Attr("appId", appId)
	}
	v, err := query.FindOne(new(agents.Value))
	if err != nil {
		return nil, err
	}
	if v == nil {
		return nil, nil
	}
	return v.(*agents.Value), nil
}

// 查找最近的一条非错误的记录
func (this *SQLAgentValueDAO) FindLatestItemValueNoError(agentId string, appId string, itemId string) (*agents.Value, error) {
	query := NewQuery(this.TableName(agentId)).
		Attr("itemId", itemId).
		Attr("error", "").
		Node().
		Desc("createdAt")
	if len(appId) > 0 {
		query.Attr("appId", appId)
	}
	v, err := query.FindOne(new(agents.Value))
	if err != nil {
		return nil, err
	}
	if v == nil {
		return nil, nil
	}
	return v.(*agents.Value), nil
}

// 取得最近的数值记录
func (this *SQLAgentValueDAO) FindLatestItemValues(agentId string, appId string, itemId string, noticeLevel notices.NoticeLevel, lastId string, size int) ([]*agents.Value, error) {
	query := NewQuery(this.TableName(agentId))
	if len(appId) > 0 {
		query.Attr("appId", appId)
	}
	if len(itemId) > 0 {
		query.Attr("itemId", itemId)
	}
	query.Node()
	query.Limit(size)
	query.Desc("createdAt")

	if noticeLevel > 0 {
		if noticeLevel == notices.NoticeLevelInfo {
			query.Attr("noticeLevel", []interface{}{notices.NoticeLevelInfo, notices.NoticeLevelNone})
		} else {
			query.Attr("noticeLevel", noticeLevel)
		}
	}

	if len(lastId) > 0 {
		query.Gt("_id", lastId)
	}

	ones, err := query.FindOnes(new(agents.Value))
	if err != nil {
		return nil, err
	}
	result := []*agents.Value{}
	for _, one := range ones {
		result = append(result, one.(*agents.Value))
	}
	return result, nil
}

// 列出数值
func (this *SQLAgentValueDAO) ListItemValues(agentId string, appId string, itemId string, noticeLevel notices.NoticeLevel, lastId string, offset int, size int) ([]*agents.Value, error) {
	query := NewQuery(this.TableName(agentId))
	query.Attr("appId", appId)
	query.Attr("itemId", itemId)
	query.Node()
	query.Offset(offset)
	query.Limit(size)
	query.Desc("createdAt")

	if noticeLevel > 0 {
		if noticeLevel == notices.NoticeLevelInfo {
			query.Attr("noticeLevel", []interface{}{notices.NoticeLevelInfo, notices.NoticeLevelNone})
		} else {
			query.Attr("noticeLevel", noticeLevel)
		}
	}

	if len(lastId) > 0 {
		query.Lt("_id", lastId)
	}

	ones, err := query.FindOnes(new(agents.Value))
	if err != nil {
		return nil, err
	}
	result := []*agents.Value{}
	for _, one := range ones {
		result = append(result, one.(*agents.Value))
	}
	return result, nil
}

// 分组查询
func (this *SQLAgentValueDAO) QueryValues(query *Query) ([]*agents.Value, error) {
	if len(query.table) > 0 {
		this.initTable(query.table)
	}

	query.FieldMap(this.mapField)
	ones, err := query.FindOnes(new(agents.Value))
	if err != nil {
		return nil, err
	}

	result := []*agents.Value{}
	for _, one := range ones {
		result = append(result, one.(*agents.Value))
	}
	return result, nil
}

// 根据时间对值进行分组查询
func (this *SQLAgentValueDAO) GroupValuesByTime(query *Query, timeField string, result map[string]Expr) ([]*agents.Value, error) {
	query.FieldMap(this.mapField)
	query.Asc("timeFormat." + timeField)
	result["timeFormat.year"] = "timeFormat_year"
	result["timeFormat.month"] = "timeFormat_month"
	result["timeFormat.week"] = "timeFormat_week"
	result["timeFormat.day"] = "timeFormat_day"
	result["timeFormat.hour"] = "timeFormat_hour"
	result["timeFormat.minute"] = "timeFormat_minute"
	result["timeFormat.second"] = "timeFormat_second"
	ones, err := query.Group("timeFormat."+timeField, result)
	if err != nil {
		return nil, err
	}

	values := []*agents.Value{}
	for _, one := range ones {
		value := agents.NewValue()
		timeFormat := one.GetMap("timeFormat")
		one.Delete("_id", "timeFormat")
		value.Value = one
		value.TimeFormat.Year = timeFormat.GetString("year")
		value.TimeFormat.Month = timeFormat.GetString("month")
		value.TimeFormat.Week = timeFormat.GetString("week")
		value.TimeFormat.Day = timeFormat.GetString("day")
		value.TimeFormat.Hour = timeFormat.GetString("hour")
		value.TimeFormat.Minute = timeFormat.GetString("minute")
		value.TimeFormat.Second = timeFormat.GetString("second")
		values = append(values, value)
	}
	return values, nil
}

// 删除Agent相关表
func (this *SQLAgentValueDAO) DropAgentTable(agentId string) error {
	return this.driver.DropTable(this.TableName(agentId))
}

func (this *SQLAgentValueDAO) initTable(table string) {
	if isInitializedTable(table) {
		return
	}

	logs.Println("[db]check table '" + table + "'")

	switch sharedDBType {
	case "mysql":
		err := this.driver.(SQLDriverInterface).CreateTable(table, "CREATE TABLE `"+table+"` ("+
			"`id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,"+
			"`_id` varchar(24) DEFAULT NULL,"+
			"`nodeId` varchar(64) DEFAULT NULL,"+
			"`agentId` varchar(64) DEFAULT NULL,"+
			"`appId` varchar(64) DEFAULT NULL,"+
			"`itemId` varchar(64) DEFAULT NULL,"+
			"`timestamp` int(11) unsigned DEFAULT '0',"+
			"`createdAt` int(11) unsigned DEFAULT '0',"+
			"`value` json DEFAULT NULL,"+
			"`error` varchar(1024) DEFAULT NULL,"+
			"`noticeLevel` tinyint(1) unsigned DEFAULT '0',"+
			"`isNotified` tinyint(1) unsigned DEFAULT '0',"+
			"`thresholdId` varchar(64) DEFAULT NULL,"+
			"`threshold` varchar(1024) DEFAULT NULL,"+
			"`timeFormat_year` varchar(4) DEFAULT NULL,"+
			"`timeFormat_month` varchar(6) DEFAULT NULL,"+
			"`timeFormat_week` varchar(6) DEFAULT NULL,"+
			"`timeFormat_day` varchar(8) DEFAULT NULL,"+
			"`timeFormat_hour` varchar(10) DEFAULT NULL,"+
			"`timeFormat_minute` varchar(12) DEFAULT NULL,"+
			"`timeFormat_second` varchar(14) DEFAULT NULL,"+
			"`costMs` decimal(20,6) unsigned DEFAULT '0.000000',"+
			"PRIMARY KEY (`id`),"+
			"UNIQUE KEY `_id` (`_id`),"+
			"KEY `appId_itemId` (`appId`,`itemId`),"+
			"KEY `nodeId_appId_itemId` (`nodeId`,`appId`,`itemId`)"+
			") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;")
		if err != nil {
			logs.Error(err)
			removeInitializedTable(table)
		}

	case "postgres":
		err := this.driver.(SQLDriverInterface).CreateTable(table, `CREATE TABLE "public"."`+table+`" (
			"id" serial8 primary key,
			"_id" varchar(24),
			"nodeId" varchar(64),
			"agentId" varchar(64),
			"appId" varchar(64),
			"itemId" varchar(64),
			"timestamp" int4 default 0,
			"createdAt" int4 default 0,
			"value" json,
			"error" varchar(1024),
			"noticeLevel" int2 default 0,
			"isNotified" int2 default 0,
			"thresholdId" varchar(64),
			"threshold" varchar(1024),
			"timeFormat_year" varchar(4),
			"timeFormat_month" varchar(6),
			"timeFormat_week" varchar(6),
			"timeFormat_day" varchar(8),
			"timeFormat_hour" varchar(10),
			"timeFormat_minute" varchar(12),
			"timeFormat_second" varchar(14),
			"costMs" float8 default 0
		)
		;

		CREATE UNIQUE INDEX "`+table+`_id" ON "public"."`+table+`" ("_id");
		CREATE INDEX "`+table+`_appId_itemId" ON "public"."`+table+`" ("appId", "itemId");
		CREATE INDEX "`+table+`_nodeId_appId_itemId" ON "public"."`+table+`" ("nodeId", "appId", "itemId");
		`)
		if err != nil {
			logs.Error(err)
			removeInitializedTable(table)
		}
	}
}

func (this *SQLAgentValueDAO) mapField(field string) string {
	switch field {
	case "timeFormat.year":
		return "timeFormat_year"
	case "timeFormat.month":
		return "timeFormat_month"
	case "timeFormat.week":
		return "timeFormat_week"
	case "timeFormat.day":
		return "timeFormat_day"
	case "timeFormat.hour":
		return "timeFormat_hour"
	case "timeFormat.minute":
		return "timeFormat_minute"
	case "timeFormat.second":
		return "timeFormat_second"
	}
	return field
}
