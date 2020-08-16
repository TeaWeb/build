package teadb

import (
	"encoding/json"
	"errors"
	"github.com/TeaWeb/build/internal/teaconfigs/stats"
	"github.com/TeaWeb/build/internal/teadb/shared"
	"github.com/iwind/TeaGo/logs"
	"strings"
	"time"
)

type SQLServerValueDAO struct {
	BaseDAO
}

// 初始化
func (this *SQLServerValueDAO) Init() {

}

// 表名
func (this *SQLServerValueDAO) TableName(serverId string) string {
	table := "teaweb_values_server_" + serverId
	this.initTable(table)
	return table
}

// 插入新数据
func (this *SQLServerValueDAO) InsertOne(serverId string, value *stats.Value) error {
	return NewQuery(this.TableName(serverId)).
		InsertOne(value)
}

// 删除过期的数据
func (this *SQLServerValueDAO) DeleteExpiredValues(serverId string, period stats.ValuePeriod, life int) error {
	return NewQuery(this.TableName(serverId)).
		Attr("period", period).
		Lt("timestamp", time.Now().Unix()-int64(life)).
		Delete()
}

// 查询相同的数值记录
func (this *SQLServerValueDAO) FindSameItemValue(serverId string, item *stats.Value) (*stats.Value, error) {
	query := NewQuery(this.TableName(serverId))
	query.Attr("item", item.Item)
	query.Attr("period", item.Period)

	switch item.Period {
	case stats.ValuePeriodSecond:
		query.Attr("timestamp", item.Timestamp)
	case stats.ValuePeriodMinute:
		query.Attr("timeFormat_minute", item.TimeFormat.Minute)
	case stats.ValuePeriodHour:
		query.Attr("timeFormat_hour", item.TimeFormat.Hour)
	case stats.ValuePeriodDay:
		query.Attr("timeFormat_day", item.TimeFormat.Day)
	case stats.ValuePeriodWeek:
		query.Attr("timeFormat_week", item.TimeFormat.Week)
	case stats.ValuePeriodMonth:
		query.Attr("timeFormat_month", item.TimeFormat.Month)
	case stats.ValuePeriodYear:
		query.Attr("timeFormat_year", item.TimeFormat.Year)
	}

	// 参数
	if len(item.Params) > 0 {
		for k, v := range item.Params {
			query.Attr(this.driver.(SQLDriverInterface).JSONExtract("params", k), v)
		}
	} else {
		switch sharedDBType {
		case "mysql":
			query.Attr("JSON_LENGTH(params)", 0)
		case "postgres":
			query.Attr("params::\"text\"", "{}")
		default:
			return nil, errors.New("unknown database type")
		}
	}

	one, err := query.FindOne(new(stats.Value))
	if err != nil {
		return nil, err
	}
	if one == nil {
		return nil, nil
	}
	return one.(*stats.Value), nil
}

// 修改值和时间戳
func (this *SQLServerValueDAO) UpdateItemValueAndTimestamp(serverId string, valueId string, value map[string]interface{}, timestamp int64) error {
	query := NewQuery(this.TableName(serverId)).
		Attr("_id", valueId)
	valuesJSON, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return this.driver.(SQLDriverInterface).UpdateOnes(query, map[string]interface{}{
		"value": valuesJSON,
	})
}

// 创建索引
func (this *SQLServerValueDAO) CreateIndex(serverId string, fields []*shared.IndexField) error {
	// 无法为JSON字段创建索引，除非是虚拟字段，所以暂时不实现
	return nil
}

// 查询数据
func (this *SQLServerValueDAO) QueryValues(query *Query) ([]*stats.Value, error) {
	query.fieldMapping = this.mapField
	ones, err := query.FindOnes(new(stats.Value))
	if err != nil {
		return nil, err
	}
	result := []*stats.Value{}
	for _, one := range ones {
		result = append(result, one.(*stats.Value))
	}
	return result, err
}

// 根据item查找一条数据
func (this *SQLServerValueDAO) FindOneWithItem(serverId string, item string) (*stats.Value, error) {
	one, err := NewQuery(this.TableName(serverId)).
		Attr("item", item).
		FindOne(new(stats.Value))
	if err != nil {
		return nil, err
	}
	if one == nil {
		return nil, err
	}
	return one.(*stats.Value), nil
}

// 删除代理服务相关表
func (this *SQLServerValueDAO) DropServerTable(serverId string) error {
	return this.driver.DropTable(this.TableName(serverId))
}

func (this *SQLServerValueDAO) initTable(table string) {
	if isInitializedTable(table) {
		return
	}

	logs.Println("[db]check table '" + table + "'")

	switch sharedDBType {
	case "mysql":
		err := this.driver.(SQLDriverInterface).CreateTable(table, "CREATE TABLE `"+table+"` ("+
			"`id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,"+
			"`_id` varchar(24) DEFAULT NULL,"+
			"`item` varchar(256) DEFAULT NULL,"+
			"`period` varchar(64) DEFAULT NULL,"+
			"`value` json DEFAULT NULL,"+
			"`params` json DEFAULT NULL,"+
			"`timestamp` int(11) unsigned DEFAULT '0',"+
			"`timeFormat_year` varchar(4) DEFAULT NULL,"+
			"`timeFormat_month` varchar(6) DEFAULT NULL,"+
			"`timeFormat_week` varchar(6) DEFAULT NULL,"+
			"`timeFormat_day` varchar(8) DEFAULT NULL,"+
			"`timeFormat_hour` varchar(10) DEFAULT NULL,"+
			"`timeFormat_minute` varchar(12) DEFAULT NULL,"+
			"`timeFormat_second` varchar(14) DEFAULT NULL,"+
			"PRIMARY KEY (`id`),"+
			"UNIQUE KEY `_id` (`_id`),"+
			"KEY `item_timestamp` (`item`,`timestamp`),"+
			"KEY `item_second` (`item`,`timeFormat_second`),"+
			"KEY `item_minute` (`item`,`timeFormat_minute`),"+
			"KEY `item_hour` (`item`,`timeFormat_hour`),"+
			"KEY `item_day` (`item`,`timeFormat_day`),"+
			"KEY `item_week` (`item`,`timeFormat_week`),"+
			"KEY `item_month` (`item`,`timeFormat_month`),"+
			"KEY `item_year` (`item`,`timeFormat_year`)"+
			") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;")
		if err != nil {
			logs.Error(err)
			removeInitializedTable(table)
		}

	case "postgres":
		err := this.driver.(SQLDriverInterface).CreateTable(table, `CREATE TABLE "public"."`+table+`" (
			"id" serial8 primary key,
			"_id" varchar(24),
			"item" varchar(256),
			"period" varchar(64),
			"value" json,
			"params" json,
			"timestamp" int4 default 0,
			"timeFormat_year" varchar(4),
			"timeFormat_month" varchar(6),
			"timeFormat_week" varchar(6),
			"timeFormat_day" varchar(8),
			"timeFormat_hour" varchar(10),
			"timeFormat_minute" varchar(12),
			"timeFormat_second" varchar(14)
		)
		;

		CREATE UNIQUE INDEX "`+table+`_id" ON "public"."`+table+`" ("_id");
		CREATE INDEX "`+table+`_item_timestamp" ON "public"."`+table+`" ("item", "timestamp");
		CREATE INDEX "`+table+`_item_timeFormat_year" ON "public"."`+table+`" ("item", "timeFormat_year");
		CREATE INDEX "`+table+`_item_timeFormat_month" ON "public"."`+table+`" ("item", "timeFormat_month");
		CREATE INDEX "`+table+`_item_timeFormat_week" ON "public"."`+table+`" ("item", "timeFormat_week");
		CREATE INDEX "`+table+`_item_timeFormat_day" ON "public"."`+table+`" ("item", "timeFormat_day");
		CREATE INDEX "`+table+`_item_timeFormat_hour" ON "public"."`+table+`" ("item", "timeFormat_hour");
		CREATE INDEX "`+table+`_item_timeFormat_minute" ON "public"."`+table+`" ("item", "timeFormat_minute");
		CREATE INDEX "`+table+`_item_timeFormat_second" ON "public"."`+table+`" ("item", "timeFormat_second");
		`)
		if err != nil {
			logs.Error(err)
			removeInitializedTable(table)
		}
	}
}

func (this *SQLServerValueDAO) mapField(field string) string {
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

	// params
	if strings.HasPrefix(field, "params.") {
		return this.driver.(SQLDriverInterface).JSONExtract("params", field[len("params."):])
	}

	return field
}
