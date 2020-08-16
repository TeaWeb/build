package teadb

import (
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/TeaWeb/build/internal/teadb/shared"
	"github.com/TeaWeb/build/internal/tealogs/accesslogs"
	"github.com/go-sql-driver/mysql"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	timeutil "github.com/iwind/TeaGo/utils/time"
	"github.com/lib/pq"
	"strings"
)

type SQLAccessLogDAO struct {
	BaseDAO
}

// 初始化
func (this *SQLAccessLogDAO) Init() {
	return
}

// 获取表名
func (this *SQLAccessLogDAO) TableName(day string) string {
	if day == timeutil.Format("Ymd") {
		this.initTable("teaweb_logs_" + day)
	}
	return "teaweb_logs_" + day
}

// 获取当前时间表名
func (this *SQLAccessLogDAO) TodayTableName() string {
	return this.TableName(timeutil.Format("Ymd"))
}

// 写入一条日志
func (this *SQLAccessLogDAO) InsertOne(accessLog *accesslogs.AccessLog) error {
	if accessLog.Id.IsZero() {
		accessLog.Id = shared.NewObjectId()
	}
	return NewQuery(this.TodayTableName()).
		InsertOne(accessLog)
}

// 写入一组日志
func (this *SQLAccessLogDAO) InsertAccessLogs(accessLogList []interface{}) error {
	return NewQuery(this.TodayTableName()).
		InsertOnes(accessLogList)
}

// 查找某条访问日志的cookie信息
func (this *SQLAccessLogDAO) FindAccessLogCookie(day string, logId string) (*accesslogs.AccessLog, error) {
	one, err := NewQuery(this.TableName(day)).
		Attr("_id", logId).
		Result("_id", "cookie").
		FindOne(new(accesslogs.AccessLog))
	if err != nil {
		if this.tableNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	if one == nil {
		return nil, nil
	}
	return one.(*accesslogs.AccessLog), nil
}

// 查找某条访问日志的请求信息
func (this *SQLAccessLogDAO) FindRequestHeaderAndBody(day string, logId string) (*accesslogs.AccessLog, error) {
	one, err := NewQuery(this.TableName(day)).
		Attr("_id", logId).
		Result("_id", "header", "requestData").
		FindOne(new(accesslogs.AccessLog))
	if err != nil {
		if this.tableNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	if one == nil {
		return nil, nil
	}
	return one.(*accesslogs.AccessLog), nil
}

// 查找某条访问日志的响应信息
func (this *SQLAccessLogDAO) FindResponseHeaderAndBody(day string, logId string) (*accesslogs.AccessLog, error) {
	one, err := NewQuery(this.TableName(day)).
		Attr("_id", logId).
		Result("_id", "sentHeader", "responseBodyData").
		FindOne(new(accesslogs.AccessLog))
	if err != nil {
		if this.tableNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	if one == nil {
		return nil, nil
	}
	return one.(*accesslogs.AccessLog), nil
}

// 列出日志
func (this *SQLAccessLogDAO) ListAccessLogs(day string, serverId string, fromId string, onlyErrors bool, searchIP string, offset int, size int) ([]*accesslogs.AccessLog, error) {
	query := NewQuery(this.TableName(day))
	query.Attr("serverId", serverId)
	if len(fromId) > 0 {
		query.Lt("_id", fromId)
	}
	if onlyErrors {
		query.Or([]*OperandList{
			NewOperandList().Add("hasErrors", NewOperand(OperandEq, 1)),
			NewOperandList().Add("status", NewOperand(OperandGte, 400)),
		})
	}
	if len(searchIP) > 0 {
		query.Attr("remoteAddr", searchIP)
	}
	query.Offset(offset)
	query.Limit(size)
	query.Desc("_id")
	ones, err := query.FindOnes(new(accesslogs.AccessLog))
	if err != nil {
		if this.tableNotFound(err) {
			return nil, nil
		}
		return nil, err
	}

	result := []*accesslogs.AccessLog{}
	for _, one := range ones {
		result = append(result, one.(*accesslogs.AccessLog))
	}
	return result, nil
}

// 检查是否有下一条日志
func (this *SQLAccessLogDAO) HasNextAccessLog(day string, serverId string, fromId string, onlyErrors bool, searchIP string) (bool, error) {
	query := NewQuery(this.TableName(day))
	query.Attr("serverId", serverId).
		Result("_id")
	if len(fromId) > 0 {
		query.Lt("_id", fromId)
	}
	if onlyErrors {
		query.Or([]*OperandList{
			NewOperandList().Add("hasErrors", NewOperand(OperandEq, 1)),
			NewOperandList().Add("status", NewOperand(OperandGte, 400)),
		})
	}
	if len(searchIP) > 0 {
		query.Attr("remoteAddr", searchIP)
	}

	one, err := query.FindOne(new(accesslogs.AccessLog))
	if err != nil {
		if this.tableNotFound(err) {
			return false, nil
		}
		return false, err
	}
	return one != nil, nil
}

// 判断某个代理服务是否有日志
func (this *SQLAccessLogDAO) HasAccessLog(day string, serverId string) (bool, error) {
	query := NewQuery(this.TableName(day))
	one, err := query.Attr("serverId", serverId).
		Result("_id").
		FindOne(new(accesslogs.AccessLog))
	if err != nil && this.tableNotFound(err) {
		return false, nil
	}
	return one != nil, err
}

// 列出WAF日志
func (this *SQLAccessLogDAO) ListAccessLogsWithWAF(day string, wafId string, fromId string, onlyErrors bool, searchIP string, offset int, size int) ([]*accesslogs.AccessLog, error) {
	query := NewQuery(this.TableName(day))
	query.Attr(this.driver.(SQLDriverInterface).JSONExtract("attrs", "waf_id"), wafId)
	if len(fromId) > 0 {
		query.Lt("_id", fromId)
	}
	if onlyErrors {
		query.Or([]*OperandList{
			NewOperandList().Add("hasErrors", NewOperand(OperandEq, 1)),
			NewOperandList().Add("status", NewOperand(OperandGte, 400)),
		})
	}
	if len(searchIP) > 0 {
		query.Attr("remoteAddr", searchIP)
	}
	query.Offset(offset)
	query.Limit(size)
	query.Desc("_id")
	ones, err := query.FindOnes(new(accesslogs.AccessLog))
	if err != nil {
		if this.tableNotFound(err) {
			return nil, nil
		}
		return nil, err
	}

	result := []*accesslogs.AccessLog{}
	for _, one := range ones {
		result = append(result, one.(*accesslogs.AccessLog))
	}
	return result, nil
}

// 检查是否有下一条日志
func (this *SQLAccessLogDAO) HasNextAccessLogWithWAF(day string, wafId string, fromId string, onlyErrors bool, searchIP string) (bool, error) {
	query := NewQuery(this.TableName(day))
	query.Attr(this.driver.(SQLDriverInterface).JSONExtract("attrs", "waf_id"), wafId).
		Result("_id")
	if len(fromId) > 0 {
		query.Lt("_id", fromId)
	}
	if onlyErrors {
		query.Or([]*OperandList{
			NewOperandList().Add("hasErrors", NewOperand(OperandEq, 1)),
			NewOperandList().Add("status", NewOperand(OperandGte, 400)),
		})
	}
	if len(searchIP) > 0 {
		query.Attr("remoteAddr", searchIP)
	}

	one, err := query.FindOne(new(accesslogs.AccessLog))
	if err != nil {
		if this.tableNotFound(err) {
			return false, nil
		}
		return false, err
	}
	return one != nil, nil
}

// 判断某个WAF是否有日志
func (this *SQLAccessLogDAO) HasAccessLogWithWAF(day string, wafId string) (bool, error) {
	query := NewQuery(this.TableName(day))
	one, err := query.Attr(this.driver.(SQLDriverInterface).JSONExtract("attrs", "waf_id"), wafId).
		Result("_id").
		FindOne(new(accesslogs.AccessLog))
	if err != nil && this.tableNotFound(err) {
		return false, nil
	}
	return one != nil, err
}

func (this *SQLAccessLogDAO) GroupWAFRuleGroups(day string, wafId string) ([]maps.Map, error) {
	waf := teaconfigs.SharedWAFList().FindWAF(wafId)
	if waf == nil {
		return []maps.Map{}, nil
	}

	driver := this.driver.(SQLDriverInterface)
	query := NewQuery(this.TableName(day))
	ones, err := query.
		Attr(driver.JSONExtract("attrs", "waf_id"), wafId).
		Group(driver.JSONExtract("attrs", "waf_group"), map[string]Expr{
			"groupId": "attrs.waf_group",
			"count":   "COUNT(_id)",
		})
	if err != nil {
		return []maps.Map{}, err
	}

	result := []maps.Map{}
	for _, one := range ones {
		groupId := strings.Trim(one.GetString("groupId"), "\"")
		group := waf.FindRuleGroup(groupId)
		if group == nil {
			continue
		}

		result = append(result, maps.Map{
			"name":  group.Name,
			"count": one.GetInt("count"),
		})
	}

	return result, err
}

// 列出最近的某些日志
func (this *SQLAccessLogDAO) ListLatestAccessLogs(day string, serverId string, fromId string, onlyErrors bool, size int) ([]*accesslogs.AccessLog, error) {
	query := NewQuery(this.TableName(day))

	shouldReverse := true
	query.Attr("serverId", serverId)
	if len(fromId) > 0 {
		query.Gt("_id", fromId)
		query.Asc("_id")
	} else {
		query.Desc("_id")
		shouldReverse = false
	}
	if onlyErrors {
		query.Or([]*OperandList{
			NewOperandList().Add("hasErrors", NewOperand(OperandEq, 1)),
			NewOperandList().Add("status", NewOperand(OperandGte, 400)),
		})
	}
	query.Limit(size)
	ones, err := query.FindOnes(new(accesslogs.AccessLog))
	if err != nil {
		if this.tableNotFound(err) {
			return nil, nil
		}
		return nil, err
	}

	if shouldReverse {
		lists.Reverse(ones)
	}

	result := []*accesslogs.AccessLog{}
	for _, one := range ones {
		result = append(result, one.(*accesslogs.AccessLog))
	}

	return result, nil
}

// 列出某天的一些日志
func (this *SQLAccessLogDAO) ListTopAccessLogs(day string, size int) ([]*accesslogs.AccessLog, error) {
	ones, err := NewQuery(this.TableName(day)).
		Limit(size).
		Desc("_id").
		FindOnes(new(accesslogs.AccessLog))
	if err != nil {
		if this.tableNotFound(err) {
			return nil, nil
		}
		return nil, err
	}

	result := []*accesslogs.AccessLog{}
	for _, one := range ones {
		result = append(result, one.(*accesslogs.AccessLog))
	}
	return result, nil
}

// 根据查询条件来查找日志
func (this *SQLAccessLogDAO) QueryAccessLogs(day string, serverId string, query *Query) ([]*accesslogs.AccessLog, error) {
	query.table = this.TableName(day)
	ones, err := query.
		Attr("serverId", serverId).
		FindOnes(new(accesslogs.AccessLog))
	if err != nil {
		if this.tableNotFound(err) {
			return nil, nil
		}
		return nil, err
	}

	result := []*accesslogs.AccessLog{}
	for _, one := range ones {
		result = append(result, one.(*accesslogs.AccessLog))
	}
	return result, nil
}

func (this *SQLAccessLogDAO) initTable(table string) {
	if isInitializedTable(table) {
		return
	}

	logs.Println("[db]check table '" + table + "'")
	switch sharedDBType {
	case "mysql":
		err := this.driver.(SQLDriverInterface).CreateTable(table, "CREATE TABLE `"+table+"` ("+
			"`id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,"+
			"`_id` varchar(24) DEFAULT NULL,"+
			"`serverId` varchar(64) DEFAULT NULL,"+
			"`backendId` varchar(64) DEFAULT NULL,"+
			"`locationId` varchar(64) DEFAULT NULL,"+
			"`fastcgiId` varchar(64) DEFAULT NULL,"+
			"`rewriteId` varchar(64) DEFAULT NULL,"+
			"`teaVersion` varchar(32) DEFAULT NULL,"+
			"`remoteAddr` varchar(64) DEFAULT NULL,"+
			"`remotePort` int(11) unsigned DEFAULT '0',"+
			"`remoteUser` varchar(128) DEFAULT NULL,"+
			"`requestURI` varchar(2048) DEFAULT NULL,"+
			"`requestPath` varchar(2048) DEFAULT NULL,"+
			"`requestLength` bigint(20) unsigned DEFAULT '0',"+
			"`requestTime` decimal(20,6) unsigned DEFAULT '0.000000',"+
			"`requestMethod` varchar(16) DEFAULT NULL,"+
			"`requestFilename` varchar(1024) DEFAULT NULL,"+
			"`scheme` varchar(16) DEFAULT NULL,"+
			"`proto` varchar(16) DEFAULT NULL,"+
			"`bytesSent` bigint(20) unsigned DEFAULT '0',"+
			"`bodyBytesSent` bigint(20) unsigned DEFAULT '0',"+
			"`status` int(11) unsigned DEFAULT '0',"+
			"`statusMessage` varchar(1024) DEFAULT NULL,"+
			"`sentHeader` json DEFAULT NULL,"+
			"`timeISO8601` varchar(128) DEFAULT NULL,"+
			"`timeLocal` varchar(128) DEFAULT NULL,"+
			"`msec` decimal(20,6) unsigned DEFAULT '0.000000',"+
			"`timestamp` int(11) unsigned DEFAULT '0',"+
			"`host` varchar(128) DEFAULT NULL,"+
			"`referer` varchar(2048) DEFAULT NULL,"+
			"`userAgent` varchar(1024) DEFAULT NULL,"+
			"`request` varchar(2048) DEFAULT NULL,"+
			"`contentType` varchar(256) DEFAULT NULL,"+
			"`cookie` json DEFAULT NULL,"+
			"`arg` json DEFAULT NULL,"+
			"`args` text,"+
			"`queryString` text,"+
			"`header` json DEFAULT NULL,"+
			"`serverName` varchar(256) DEFAULT NULL,"+
			"`serverPort` int(11) unsigned DEFAULT '0',"+
			"`serverProtocol` varchar(16) DEFAULT NULL,"+
			"`backendAddress` varchar(256) DEFAULT NULL,"+
			"`fastcgiAddress` varchar(256) DEFAULT NULL,"+
			"`requestData` longblob,"+
			"`responseHeaderData` longblob,"+
			"`responseBodyData` longblob,"+
			"`errors` json DEFAULT NULL,"+
			"`hasErrors` tinyint(1) unsigned DEFAULT '0',"+
			"`extend` json DEFAULT NULL,"+
			"`attrs` json DEFAULT NULL,"+
			"PRIMARY KEY (`id`),"+
			"UNIQUE KEY `_id` (`_id`),"+
			"KEY `serverId` (`serverId`),"+
			"KEY `serverId_status` (`serverId`,`status`),"+
			"KEY `serverId_remoteAddr` (`serverId`,`remoteAddr`),"+
			"KEY `serverId_hasErrors` (`serverId`,`hasErrors`)"+
			") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;")
		if err != nil {
			logs.Error(err)
			removeInitializedTable(table)
		}
	case "postgres":
		err := this.driver.(SQLDriverInterface).CreateTable(table, `CREATE TABLE "public"."`+table+`" (
			"id" serial8 primary key,
			"_id" varchar(24),
			"serverId" varchar(64),
			"backendId" varchar(64),
			"locationId" varchar(64),
			"fastcgiId" varchar(64),
			"rewriteId" varchar(64),
			"teaVersion" varchar(32),
			"remoteAddr" varchar(64),
			"remotePort" int4 default 0,
			"remoteUser" varchar(128),
			"requestURI" varchar(2048),
			"requestPath" varchar(2048),
			"requestLength" int8 default 0,
			"requestTime" float8 default 0,
			"requestMethod" varchar(16),
			"requestFilename" varchar(1024),
			"scheme" varchar(16),
			"proto" varchar(16),
			"bytesSent" int8 default 0,
			"bodyBytesSent" int8 default 0,
			"status" int4 default 0,
			"statusMessage" varchar(1024),
			"sentHeader" json,
			"timeISO8601" varchar(128),
			"timeLocal" varchar(128),
			"msec" float8 default 0,
			"timestamp" int4 default 0,
			"host" varchar(128),
			"referer" varchar(2048),
			"userAgent" varchar(1024),
			"request" varchar(2048),
			"contentType" varchar(256),
			"cookie" json,
			"arg" json,
			"args" text,
			"queryString" text,
			"header" json,
			"serverName" varchar(256),
			"serverPort" int4 default 0,
			"serverProtocol" varchar(16),
			"backendAddress" varchar(256),
			"fastcgiAddress" varchar(256),
			"requestData" bytea,
			"responseHeaderData" bytea,
			"responseBodyData" bytea,
			"errors" json,
			"hasErrors" int2 default 0,
			"extend" json,
			"attrs" json
		)
		;

		CREATE UNIQUE INDEX "`+table+`_id" ON "public"."`+table+`" ("_id");
		CREATE INDEX "`+table+`_serverId_status" ON "public"."`+table+`" ("serverId", "status");
		CREATE INDEX "`+table+`_serverId_remoteAddr" ON "public"."`+table+`" ("serverId", "remoteAddr");
		CREATE INDEX "`+table+`_serverId_hasErrors" ON "public"."`+table+`" ("serverId", "hasErrors");
		`)
		if err != nil {
			logs.Error(err)
			removeInitializedTable(table)
		}
	}
}

func (this *SQLAccessLogDAO) tableNotFound(err error) bool {
	if err == nil {
		return false
	}

	mysqlErr, ok := err.(*mysql.MySQLError)
	if ok {
		return mysqlErr.Number == 1146
	}

	pqErr, ok := err.(*pq.Error)
	if ok {
		return pqErr.Code == "42P01"
	}

	return false
}
