package teadb

import (
	"context"
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/TeaWeb/build/internal/teadb/shared"
	"github.com/TeaWeb/build/internal/tealogs/accesslogs"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	timeutil "github.com/iwind/TeaGo/utils/time"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"strings"
	"time"
)

type MongoAccessLogDAO struct {
	BaseDAO
}

func (this *MongoAccessLogDAO) Init() {

}

func (this *MongoAccessLogDAO) TableName(day string) string {
	table := "logs." + day
	return table
}

func (this *MongoAccessLogDAO) InsertingTableName(day string) string {
	table := "logs." + day
	this.initTable(table)
	return table
}

// 写入一条日志
func (this *MongoAccessLogDAO) InsertOne(accessLog *accesslogs.AccessLog) error {
	if accessLog.Id.IsZero() {
		accessLog.Id = shared.NewObjectId()
	}
	return NewQuery(this.InsertingTableName(timeutil.Format("Ymd"))).
		InsertOne(accessLog)
}

// 写入一组日志
func (this *MongoAccessLogDAO) InsertAccessLogs(accessLogList []interface{}) error {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	coll, err := this.driver.(*MongoDriver).SelectColl(this.InsertingTableName(timeutil.Format("Ymd")))
	if err != nil {
		return err
	}
	_, err = coll.InsertMany(ctx, accessLogList)
	return err
}

func (this *MongoAccessLogDAO) FindAccessLogCookie(day string, logId string) (*accesslogs.AccessLog, error) {
	idObject, err := shared.ObjectIdFromHex(logId)
	if err != nil {
		return nil, err
	}

	one, err := NewQuery(this.TableName(day)).
		Attr("_id", idObject).
		Result("cookie").
		FindOne(new(accesslogs.AccessLog))
	if err != nil {
		return nil, err
	}
	if one == nil {
		return nil, nil
	}
	return one.(*accesslogs.AccessLog), nil
}

func (this *MongoAccessLogDAO) FindRequestHeaderAndBody(day string, logId string) (*accesslogs.AccessLog, error) {
	idObject, err := shared.ObjectIdFromHex(logId)
	if err != nil {
		return nil, err
	}
	one, err := NewQuery(this.TableName(day)).
		Attr("_id", idObject).
		Result("header", "requestData").
		FindOne(new(accesslogs.AccessLog))
	if err != nil {
		return nil, err
	}
	if one == nil {
		return nil, nil
	}
	return one.(*accesslogs.AccessLog), nil
}

func (this *MongoAccessLogDAO) FindResponseHeaderAndBody(day string, logId string) (*accesslogs.AccessLog, error) {
	idObject, err := shared.ObjectIdFromHex(logId)
	if err != nil {
		return nil, err
	}
	one, err := NewQuery(this.TableName(day)).
		Attr("_id", idObject).
		Result("sentHeader", "responseBodyData").
		FindOne(new(accesslogs.AccessLog))
	if err != nil {
		return nil, err
	}
	if one == nil {
		return nil, nil
	}
	return one.(*accesslogs.AccessLog), nil
}

func (this *MongoAccessLogDAO) ListAccessLogs(day string, serverId string, fromId string, onlyErrors bool, searchIP string, offset int, size int) ([]*accesslogs.AccessLog, error) {
	query := NewQuery(this.TableName(day))
	query.Attr("serverId", serverId)
	if len(fromId) > 0 {
		fromIdObject, err := shared.ObjectIdFromHex(fromId)
		if err != nil {
			return nil, err
		}
		query.Lt("_id", fromIdObject)
	}
	if onlyErrors {
		query.Or([]*OperandList{
			NewOperandList().Add("hasErrors", NewOperand(OperandEq, true)),
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
		return nil, err
	}

	result := []*accesslogs.AccessLog{}
	for _, one := range ones {
		result = append(result, one.(*accesslogs.AccessLog))
	}
	return result, nil
}

func (this *MongoAccessLogDAO) HasNextAccessLog(day string, serverId string, fromId string, onlyErrors bool, searchIP string) (bool, error) {
	query := NewQuery(this.TableName(day))
	query.Attr("serverId", serverId).
		Result("_id")
	if len(fromId) > 0 {
		fromIdObject, err := shared.ObjectIdFromHex(fromId)
		if err != nil {
			return false, err
		}
		query.Lt("_id", fromIdObject)
	}
	if onlyErrors {
		query.Or([]*OperandList{
			NewOperandList().Add("hasErrors", NewOperand(OperandEq, true)),
			NewOperandList().Add("status", NewOperand(OperandGte, 400)),
		})
	}
	if len(searchIP) > 0 {
		query.Attr("remoteAddr", searchIP)
	}

	one, err := query.FindOne(new(accesslogs.AccessLog))
	if err != nil {
		return false, err
	}
	return one != nil, nil
}

func (this *MongoAccessLogDAO) HasAccessLog(day string, serverId string) (bool, error) {
	query := NewQuery(this.TableName(day))
	one, err := query.Attr("serverId", serverId).
		Result("_id").
		FindOne(new(accesslogs.AccessLog))
	return one != nil, err
}

func (this *MongoAccessLogDAO) ListAccessLogsWithWAF(day string, wafId string, fromId string, onlyErrors bool, searchIP string, offset int, size int) ([]*accesslogs.AccessLog, error) {
	query := NewQuery(this.TableName(day))
	query.Attr("attrs.waf_id", wafId)
	if len(fromId) > 0 {
		fromIdObject, err := shared.ObjectIdFromHex(fromId)
		if err != nil {
			return nil, err
		}
		query.Lt("_id", fromIdObject)
	}
	if onlyErrors {
		query.Or([]*OperandList{
			NewOperandList().Add("hasErrors", NewOperand(OperandEq, true)),
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
		return nil, err
	}

	result := []*accesslogs.AccessLog{}
	for _, one := range ones {
		result = append(result, one.(*accesslogs.AccessLog))
	}
	return result, nil
}

func (this *MongoAccessLogDAO) HasNextAccessLogWithWAF(day string, wafId string, fromId string, onlyErrors bool, searchIP string) (bool, error) {
	query := NewQuery(this.TableName(day))
	query.Attr("attrs.waf_id", wafId).
		Result("_id")
	if len(fromId) > 0 {
		fromIdObject, err := shared.ObjectIdFromHex(fromId)
		if err != nil {
			return false, err
		}
		query.Lt("_id", fromIdObject)
	}
	if onlyErrors {
		query.Or([]*OperandList{
			NewOperandList().Add("hasErrors", NewOperand(OperandEq, true)),
			NewOperandList().Add("status", NewOperand(OperandGte, 400)),
		})
	}
	if len(searchIP) > 0 {
		query.Attr("remoteAddr", searchIP)
	}

	one, err := query.FindOne(new(accesslogs.AccessLog))
	if err != nil {
		return false, err
	}
	return one != nil, nil
}

func (this *MongoAccessLogDAO) HasAccessLogWithWAF(day string, wafId string) (bool, error) {
	query := NewQuery(this.TableName(day))
	one, err := query.Attr("attrs.waf_id", wafId).
		Result("_id").
		FindOne(new(accesslogs.AccessLog))
	return one != nil, err
}

func (this *MongoAccessLogDAO) GroupWAFRuleGroups(day string, wafId string) ([]maps.Map, error) {
	waf := teaconfigs.SharedWAFList().FindWAF(wafId)
	if waf == nil {
		return []maps.Map{}, nil
	}

	query := NewQuery(this.TableName(day))
	ones, err := query.
		Attr("attrs.waf_id", wafId).
		Group("attrs.waf_group", map[string]Expr{
			"groupId": "attrs.waf_group",
			"count": maps.Map{
				"$sum": 1,
			},
		})
	if err != nil {
		return nil, err
	}

	result := []maps.Map{}
	for _, one := range ones {
		groupId := one.GetString("groupId")
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

func (this *MongoAccessLogDAO) ListLatestAccessLogs(day string, serverId string, fromId string, onlyErrors bool, size int) ([]*accesslogs.AccessLog, error) {
	query := NewQuery(this.TableName(day))

	shouldReverse := true
	query.Attr("serverId", serverId)
	if len(fromId) > 0 {
		fromIdObject, err := shared.ObjectIdFromHex(fromId)
		if err != nil {
			return nil, err
		}
		query.Gt("_id", fromIdObject)
		query.Asc("_id")
	} else {
		query.Desc("_id")
		shouldReverse = false
	}
	if onlyErrors {
		query.Or([]*OperandList{
			NewOperandList().Add("hasErrors", NewOperand(OperandEq, true)),
			NewOperandList().Add("status", NewOperand(OperandGte, 400)),
		})
	}
	query.Limit(size)
	ones, err := query.FindOnes(new(accesslogs.AccessLog))
	if err != nil {
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

func (this *MongoAccessLogDAO) ListTopAccessLogs(day string, size int) ([]*accesslogs.AccessLog, error) {
	ones, err := NewQuery(this.TableName(day)).
		Limit(size).
		Desc("_id").
		FindOnes(new(accesslogs.AccessLog))
	if err != nil {
		return nil, err
	}

	result := []*accesslogs.AccessLog{}
	for _, one := range ones {
		result = append(result, one.(*accesslogs.AccessLog))
	}
	return result, nil
}

func (this *MongoAccessLogDAO) QueryAccessLogs(day string, serverId string, query *Query) ([]*accesslogs.AccessLog, error) {
	query.table = this.TableName(day)
	ones, err := query.
		Attr("serverId", serverId).
		FindOnes(new(accesslogs.AccessLog))
	if err != nil {
		return nil, err
	}

	result := []*accesslogs.AccessLog{}
	for _, one := range ones {
		result = append(result, one.(*accesslogs.AccessLog))
	}
	return result, nil
}

func (this *MongoAccessLogDAO) initTable(table string) {
	if isInitializedTable(table) {
		return
	}

	// 异步执行，防止阻塞进程
	go func() {
		for _, fields := range [][]*shared.IndexField{
			{
				shared.NewIndexField("serverId", true),
			},
			{
				shared.NewIndexField("status", true),
				shared.NewIndexField("serverId", true),
			},
			{
				shared.NewIndexField("remoteAddr", true),
				shared.NewIndexField("serverId", true),
			},
			{
				shared.NewIndexField("hasErrors", true),
				shared.NewIndexField("serverId", true),
			},
			{
				shared.NewIndexField("attrs.waf_id", true),
			},
		} {
			err := this.createIndex(table, fields)
			if err != nil {
				logs.Error(err)
			}
		}
	}()
}

func (this *MongoAccessLogDAO) createIndex(table string, fields []*shared.IndexField) error {
	if len(fields) == 0 {
		return nil
	}

	coll, err := this.driver.(*MongoDriver).SelectColl(table)
	if err != nil {
		return err
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	// 创建新的
	bsonDoc := bsonx.Doc{}
	for _, field := range fields {
		if field.Asc {
			bsonDoc = bsonDoc.Append(field.Name, bsonx.Int32(1))
		} else {
			bsonDoc = bsonDoc.Append(field.Name, bsonx.Int32(-1))
		}
	}

	_, err = coll.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bsonDoc,
		Options: options.Index().SetBackground(true),
	})

	// 忽略可能产生的冲突错误
	if err != nil && strings.Contains(err.Error(), "existing") {
		err = nil
	}

	return err
}
