package teadb

import (
	"github.com/TeaWeb/build/internal/teaconfigs/stats"
	"github.com/TeaWeb/build/internal/teadb/shared"
	"github.com/TeaWeb/build/internal/teatesting"
	"github.com/iwind/TeaGo/assert"
	"github.com/iwind/TeaGo/logs"
	stringutil "github.com/iwind/TeaGo/utils/string"
	"testing"
	"time"
)

func TestServerValueDAO_InsertOne(t *testing.T) {
	if !teatesting.RequireDBAvailable() {
		return
	}

	{
		v := stats.NewItemValue()
		v.Timestamp = time.Now().Unix() - 10
		v.Item = "a.b.c"
		v.Period = stats.ValuePeriodSecond
		v.SetTime(time.Unix(1566651441, 0))
		v.Value = map[string]interface{}{
			"name": "lu",
			"age":  20,
		}
		err := ServerValueDAO().InsertOne("test", v)
		if err != nil {
			t.Fatal(err)
		}
		t.Log("ok")
	}

	{
		v := stats.NewItemValue()
		v.Timestamp = time.Now().Unix()
		v.Item = "a.b.c"
		v.Period = stats.ValuePeriodMinute
		v.SetTime(time.Unix(1566651441, 0))
		v.Value = map[string]interface{}{
			"name": "lu",
			"age":  20,
		}
		err := ServerValueDAO().InsertOne("test", v)
		if err != nil {
			t.Fatal(err)
		}
		t.Log("ok")
	}

	{

		v := stats.NewItemValue()
		v.Timestamp = time.Now().Unix()
		v.Item = "a.b.c"
		v.Period = stats.ValuePeriodMinute
		v.Params = map[string]string{
			"param1": "name",
			"param2": "age",
		}
		v.Value = map[string]interface{}{
			"name": "lu",
			"age":  20,
		}
		v.SetTime(time.Unix(1566651441, 0))
		err := ServerValueDAO().InsertOne("test", v)
		if err != nil {
			t.Fatal(err)
		}
		t.Log("ok")
	}
}

func TestServerValueDAO_DeleteExpiredValues(t *testing.T) {
	if !teatesting.RequireDBAvailable() {
		return
	}

	dao := ServerValueDAO()
	err := dao.DeleteExpiredValues("test", stats.ValuePeriodSecond, 10)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("ok")
}

func TestServerValueDAO_FindSameItemValue(t *testing.T) {
	if !teatesting.RequireDBAvailable() {
		return
	}

	a := assert.NewAssertion(t)

	dao := ServerValueDAO()

	oldItem := &stats.Value{
		Period: stats.ValuePeriodMinute,
		Item:   "a.b.c",
		Params: map[string]string{},
	}
	oldItem.SetTime(time.Unix(1566651441, 0))
	item, err := dao.FindSameItemValue("test", oldItem)
	if err != nil {
		t.Fatal(err)
	}
	if item == nil {
		a.Log("not found")
		return
	}
	logs.PrintAsJSON(item, t)

	a.IsTrue(item.Period == oldItem.Period)
	a.IsTrue(len(item.Params) == len(oldItem.Params))
	a.IsTrue(item.TimeFormat.Minute == oldItem.TimeFormat.Minute)
}

func TestServerValueDAO_FindSameItemValue2(t *testing.T) {
	if !teatesting.RequireDBAvailable() {
		return
	}

	a := assert.NewAssertion(t)

	dao := ServerValueDAO()

	oldItem := &stats.Value{
		Period: stats.ValuePeriodMinute,
		Item:   "a.b.c",
		Params: map[string]string{
			"param1": "name",
			"param2": "age",
		},
	}
	oldItem.SetTime(time.Unix(1566651441, 0))
	item, err := dao.FindSameItemValue("test", oldItem)
	if err != nil {
		t.Fatal(err)
	}
	if item == nil {
		a.Log("not found")
		return
	}
	logs.PrintAsJSON(item, t)

	a.IsTrue(item.Period == oldItem.Period)
	a.IsTrue(len(item.Params) == len(oldItem.Params))
	a.IsTrue(item.TimeFormat.Minute == oldItem.TimeFormat.Minute)
}

func TestServerValueDAO_UpdateItemValueAndTimestamp(t *testing.T) {
	if !teatesting.RequireDBAvailable() {
		return
	}

	dao := ServerValueDAO()
	ones, err := dao.QueryValues(NewQuery(dao.TableName("test")).Limit(1).Desc("_id"))
	if err != nil {
		t.Fatal(err)
	}
	if len(ones) == 0 {
		t.Log("not found")
		return
	}
	err = dao.UpdateItemValueAndTimestamp("test", ones[0].Id.Hex(), map[string]interface{}{
		"name": "xia",
		"age":  "21",
	}, 1566651442)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("ok")
}

func TestServerValueDAO_CreateIndex(t *testing.T) {
	if !teatesting.RequireDBAvailable() {
		return
	}

	dao := ServerValueDAO()

	fields := []*shared.IndexField{
		shared.NewIndexField("param1", true),
		shared.NewIndexField("param2", false),
	}
	err := dao.CreateIndex("test", fields)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("ok")
}

func TestServerValueDAO_QueryValues(t *testing.T) {
	if !teatesting.RequireDBAvailable() {
		return
	}

	dao := ServerValueDAO()

	q := NewQuery(dao.TableName("test"))
	values, err := dao.QueryValues(q)
	if err != nil {
		t.Fatal(err)
	}

	for _, v := range values {
		t.Log(v.Id, v.Period, v.Params, stringutil.JSONEncode(v.Value), v.TimeFormat.Minute)
	}
}

func TestServerValueDAO_FindOneWithItem(t *testing.T) {
	if !teatesting.RequireDBAvailable() {
		return
	}

	dao := ServerValueDAO()

	{
		v, err := dao.FindOneWithItem("test", "a.b.c")
		if err != nil {
			t.Fatal(err)
		}
		t.Log(stringutil.JSONEncodePretty(v))
	}

	{
		v, err := dao.FindOneWithItem("test", "a.b.c1")
		if err != nil {
			t.Fatal(err)
		}
		t.Log(stringutil.JSONEncodePretty(v))
	}
}

func TestServerValueDAO_DropServerTable(t *testing.T) {
	if !teatesting.RequireDBAvailable() {
		return
	}

	err := ServerValueDAO().DropServerTable("test")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("ok")
}
