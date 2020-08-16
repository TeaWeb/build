package teadb

import (
	"errors"
	"github.com/TeaWeb/build/internal/teaconfigs/stats"
	"github.com/TeaWeb/build/internal/teadb/shared"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/types"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"golang.org/x/net/context"
	"strings"
	"time"
)

type MongoServerValueDAO struct {
	BaseDAO
}

func (this *MongoServerValueDAO) Init() {
}

func (this *MongoServerValueDAO) TableName(serverId string) string {
	return "values.server." + serverId
}

func (this *MongoServerValueDAO) InsertOne(serverId string, value *stats.Value) error {
	coll, err := this.selectColl(this.TableName(serverId))
	if err != nil {
		return err
	}

	if value.Id.IsZero() {
		value.Id = shared.NewObjectId()
	}

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	_, err = coll.InsertOne(ctx, value)
	return err
}

func (this *MongoServerValueDAO) DeleteExpiredValues(serverId string, period stats.ValuePeriod, life int) error {
	return NewQuery(this.TableName(serverId)).
		Attr("period", period).
		Lt("timestamp", time.Now().Unix()-int64(life)).
		Delete()
}

// 根据参数查询已有的数据
func (this *MongoServerValueDAO) FindSameItemValue(serverId string, item *stats.Value) (*stats.Value, error) {
	query := NewQuery(this.TableName(serverId))
	query.Attr("item", item.Item)
	query.Attr("period", item.Period)

	switch item.Period {
	case stats.ValuePeriodSecond:
		query.Attr("timestamp", item.Timestamp)
	case stats.ValuePeriodMinute:
		query.Attr("timeFormat.minute", item.TimeFormat.Minute)
	case stats.ValuePeriodHour:
		query.Attr("timeFormat.hour", item.TimeFormat.Hour)
	case stats.ValuePeriodDay:
		query.Attr("timeFormat.day", item.TimeFormat.Day)
	case stats.ValuePeriodWeek:
		query.Attr("timeFormat.week", item.TimeFormat.Week)
	case stats.ValuePeriodMonth:
		query.Attr("timeFormat.month", item.TimeFormat.Month)
	case stats.ValuePeriodYear:
		query.Attr("timeFormat.year", item.TimeFormat.Year)
	}

	// 参数
	if len(item.Params) > 0 {
		for k, v := range item.Params {
			query.Attr("params."+k, v)
		}
	} else {
		query.Attr("params", map[string]string{})
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

func (this *MongoServerValueDAO) UpdateItemValueAndTimestamp(serverId string, valueId string, value map[string]interface{}, timestamp int64) error {
	objectId, err := shared.ObjectIdFromHex(valueId)
	if err != nil {
		return err
	}

	coll, err := this.selectColl(this.TableName(serverId))
	if err != nil {
		return err
	}

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	_, err = coll.UpdateOne(ctx, map[string]interface{}{
		"_id": objectId,
	}, map[string]interface{}{
		"$set": map[string]interface{}{
			"value":     value,
			"timestamp": timestamp,
		},
	})
	return err
}

func (this *MongoServerValueDAO) CreateIndex(serverId string, fields []*shared.IndexField) error {
	if len(fields) == 0 {
		return nil
	}

	coll, err := this.selectColl(this.TableName(serverId))
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

func (this *MongoServerValueDAO) QueryValues(query *Query) ([]*stats.Value, error) {
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

func (this *MongoServerValueDAO) FindOneWithItem(serverId string, item string) (*stats.Value, error) {
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
func (this *MongoServerValueDAO) DropServerTable(serverId string) error {
	coll, err := this.selectColl(this.TableName(serverId))
	if err != nil {
		return err
	}
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	return coll.Drop(ctx)
}

func (this *MongoServerValueDAO) selectColl(collName string) (*MongoCollection, error) {
	coll, err := this.driver.(*MongoDriver).SelectColl(collName)
	if err != nil {
		return nil, err
	}

	if isInitializedTable(collName) {
		return coll, nil
	}

	for _, fields := range [][]*shared.IndexField{
		{
			shared.NewIndexField("item", true),
			shared.NewIndexField("timestamp", true),
		},
		{
			shared.NewIndexField("item", true),
			shared.NewIndexField("timeFormat.second", true),
		},
		{
			shared.NewIndexField("item", true),
			shared.NewIndexField("timeFormat.minute", true),
		},
		{
			shared.NewIndexField("item", true),
			shared.NewIndexField("timeFormat.hour", true),
		},
		{
			shared.NewIndexField("item", true),
			shared.NewIndexField("timeFormat.day", true),
		},
		{
			shared.NewIndexField("item", true),
			shared.NewIndexField("timeFormat.week", true),
		},
		{
			shared.NewIndexField("item", true),
			shared.NewIndexField("timeFormat.month", true),
		},
		{
			shared.NewIndexField("item", true),
			shared.NewIndexField("timeFormat.year", true),
		},
	} {
		err := coll.CreateIndex(fields...)
		if err != nil {
			logs.Error(errors.New("create index: " + err.Error()))
		}
	}

	return coll, nil
}

func (this *MongoServerValueDAO) checkIndexEqual(index1 map[string]interface{}, index2 map[string]interface{}) bool {
	if len(index1) != len(index2) {
		return false
	}
	for k, v := range index1 {
		v2, ok := index2[k]
		if !ok {
			return false
		}
		if types.Int(v) != types.Int(v2) {
			return false
		}
	}
	return true
}
