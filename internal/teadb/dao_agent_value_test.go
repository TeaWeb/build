package teadb

import (
	"encoding/json"
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/TeaWeb/build/internal/teaconfigs/agents"
	"github.com/TeaWeb/build/internal/teaconfigs/notices"
	"github.com/TeaWeb/build/internal/teatesting"
	"github.com/iwind/TeaGo/logs"
	timeutil "github.com/iwind/TeaGo/utils/time"
	"testing"
	"time"
)

func TestAgentValueDAO_Insert(t *testing.T) {
	dao := AgentValueDAO()
	node := teaconfigs.SharedNodeConfig()

	{
		value := agents.NewValue()
		value.AgentId = "local"
		value.AppId = "mysql"
		value.SetTime(time.Now())
		value.Value = 4
		value.ItemId = "ping"
		value.Error = "error"
		value.NoticeLevel = notices.NoticeLevelWarning
		value.IsNotified = true
		value.Threshold = "${0} gt 10"
		value.ThresholdId = "abc"
		value.CreatedAt = time.Now().Unix()

		if node != nil {
			value.NodeId = node.Id
		}

		err := dao.Insert(value.AgentId, value)
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		value := agents.NewValue()
		value.AgentId = "local"
		value.AppId = "system"
		value.SetTime(time.Now())
		value.Value = map[string]interface{}{
			"load1":  1.23,
			"load5":  2.54,
			"load15": 2.12,
		}
		value.ItemId = "cpu.load"
		value.Error = ""
		value.IsNotified = true
		value.CreatedAt = time.Now().Unix()

		if node != nil {
			value.NodeId = node.Id
		}

		err := dao.Insert(value.AgentId, value)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestAgentValueDAO_Insert2(t *testing.T) {
	jsonString := `
{
    "code": 500,
    "message": "\u8bf7\u8f93\u5165\u6b63\u786e\u7684\u4ee4\u724c\uff08001\uff09",
    "data": {},
    "next": null,
    "errors": []
}`
	v := map[string]interface{}{}
	err := json.Unmarshal([]byte(jsonString), &v)
	if err != nil {
		t.Fatal(err)
	}
	value := &agents.Value{
		AppId:       "1",
		AgentId:     "1",
		ItemId:      "1",
		Value:       v,
		Error:       "",
		NoticeLevel: notices.NoticeLevelWarning,
		CreatedAt:   time.Now().Unix(),
	}
	value.SetTime(time.Now())

	node := teaconfigs.SharedNodeConfig()
	if node != nil {
		value.NodeId = node.Id
	}

	err = AgentValueDAO().Insert("local", value)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("ok")
}

func TestAgentValueDAO_ClearItemValues(t *testing.T) {
	dao := AgentValueDAO()
	err := dao.ClearItemValues("local", "1", "1", 0)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("ok")
}

func TestAgentValuedAO_FindLatestItemValue(t *testing.T) {
	dao := AgentValueDAO()
	v, err := dao.FindLatestItemValue("local", "system", "cpu.load")
	if err != nil {
		t.Fatal(err)
	}
	if v == nil {
		t.Log("not found")
		return
	}
	logs.PrintAsJSON(v, t)
	t.Log("createdTime:", timeutil.Format("Y-m-d H:i:s", time.Unix(v.CreatedAt, 0)))
}

func TestAgentValuedAO_FindLatestItemValueNoError(t *testing.T) {
	dao := AgentValueDAO()
	v, err := dao.FindLatestItemValueNoError("local", "system", "cpu.load")
	if err != nil {
		t.Fatal(err)
	}
	if v == nil {
		t.Log("not found")
		return
	}
	logs.PrintAsJSON(v, t)
	t.Log("createdTime:", timeutil.Format("Y-m-d H:i:s", time.Unix(v.CreatedAt, 0)))
}

func TestAgentValueDAO_FindLatestItemValues(t *testing.T) {
	if !teatesting.RequireDBAvailable() {
		return
	}

	{
		dao := AgentValueDAO()
		values, err := dao.FindLatestItemValues("local", "system", "cpu.load", 0, "", 10)
		if err != nil {
			t.Fatal(err)
		}
		for _, v := range values {
			t.Log(v.Id, v.Value, v.NoticeLevel)
		}
	}

	t.Log("===all local===")
	{
		dao := AgentValueDAO()
		values, err := dao.FindLatestItemValues("local", "", "", 0, "", 10)
		if err != nil {
			t.Fatal(err)
		}
		for _, v := range values {
			t.Log(v.Id, v.Value, v.NoticeLevel)
		}
	}

	t.Log("===level & lastId===")
	{
		dao := AgentValueDAO()
		values, err := dao.FindLatestItemValues("local", "", "", notices.NoticeLevelInfo, "123", 10)
		if err != nil {
			t.Fatal(err)
		}
		for _, v := range values {
			t.Log(v.Id, v.Value, v.NoticeLevel)
		}
	}
}

func TestAgentValueDAO_ListItemValues(t *testing.T) {
	dao := AgentValueDAO()
	values, err := dao.ListItemValues("local", "system", "cpu.load", 0, "", 0, 5)
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range values {
		t.Log(v.Id, v.Value, v.NoticeLevel)
	}
}

func TestAgentValueDAO_QueryValues(t *testing.T) {
	dao := AgentValueDAO()
	q := NewQuery("teaweb.values.agent.local")
	//q.Attr("timeFormat.year", timeutil.Format("Y"))
	q.Limit(10)
	values, err := dao.QueryValues(q)
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range values {
		t.Log(v)
	}
}

func TestAgentValueDAO_GroupValues(t *testing.T) {
	dao := AgentValueDAO()

	q := NewQuery("teaweb.values.agent.local").
		Attr("itemId", "cpu.load")

	values, err := dao.GroupValuesByTime(q, "day", map[string]Expr{
		"load1":  NewAvgExpr("value.load1"),
		"load5":  "value.load5",
		"load15": "value.load15",
	})
	if err != nil {
		t.Fatal(err)
	}

	for _, v := range values {
		t.Log(v.TimeFormat.Day, v.Value)
	}
}

func TestAgentValueDAO_DropAgentTable(t *testing.T) {
	dao := AgentValueDAO()
	_, err := dao.FindLatestItemValue("test", "", "")
	if err != nil {
		t.Fatal(err)
	}

	err = dao.DropAgentTable("test")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("ok")
}
