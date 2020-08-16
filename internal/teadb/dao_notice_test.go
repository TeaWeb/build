package teadb

import (
	"github.com/TeaWeb/build/internal/teaconfigs/notices"
	"github.com/TeaWeb/build/internal/teatesting"
	"github.com/iwind/TeaGo/assert"
	"testing"
	"time"
)

func TestNoticeDAO_InsertOne(t *testing.T) {
	if !teatesting.RequireDBAvailable() {
		return
	}

	{
		notice := notices.NewNotice()
		notice.Message = "this is test"
		notice.Receivers = []string{"EpBqPQMqpRlvFh9Q"}
		notice.Agent.AgentId = "local"
		notice.Agent.AppId = "system"
		notice.Agent.ItemId = "cpu.load"
		notice.SetTime(time.Now())
		notice.Hash()

		dao := NoticeDAO()
		err := dao.InsertOne(notice)
		if err != nil {
			t.Fatal(err)
		}
		t.Log("ok")
	}

	{
		notice := notices.NewNotice()
		notice.Message = "this is test"
		notice.Receivers = []string{"EpBqPQMqpRlvFh9Q"}
		notice.Agent.AgentId = "local"
		notice.Agent.AppId = "system"
		notice.Agent.ItemId = "cpu.load"
		notice.IsRead = true
		notice.SetTime(time.Now())
		notice.Hash()

		dao := NoticeDAO()
		err := dao.InsertOne(notice)
		if err != nil {
			t.Fatal(err)
		}
		t.Log("ok")
	}

	{
		notice := notices.NewNotice()
		notice.Message = "this is test"
		notice.Receivers = []string{"EpBqPQMqpRlvFh9Q"}
		notice.Agent.AgentId = "1TABzdF0uAIFPGkr"
		notice.Agent.AppId = "system"
		notice.Agent.ItemId = "cpu.load"
		notice.SetTime(time.Now())
		notice.Hash()

		dao := NoticeDAO()
		err := dao.InsertOne(notice)
		if err != nil {
			t.Fatal(err)
		}
		t.Log("ok")
	}
}

func TestNoticeDAO_NotifyProxyMessage(t *testing.T) {
	if !teatesting.RequireDBAvailable() {
		return
	}

	dao := NoticeDAO()
	err := dao.NotifyProxyMessage(notices.ProxyCond{
		ServerId: "test",
	}, "this is test")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("ok")
}

func TestNoticeDAO_NotifyProxyServerMessage(t *testing.T) {
	if !teatesting.RequireDBAvailable() {
		return
	}

	dao := NoticeDAO()
	err := dao.NotifyProxyServerMessage("test2", notices.NoticeLevelWarning, "Hello")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("ok")
}

func TestNoticeDAO_CountAllCountUnreadNotices(t *testing.T) {
	if !teatesting.RequireDBAvailable() {
		return
	}

	dao := NoticeDAO()
	count, err := dao.CountAllUnreadNotices()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("count:", count)
}

func TestNoticeDAO_CountAllCountReadNotices(t *testing.T) {
	if !teatesting.RequireDBAvailable() {
		return
	}

	dao := NoticeDAO()
	count, err := dao.CountAllReadNotices()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("count:", count)
}

func TestNoticeDAO_CountUnreadNoticesForAgent(t *testing.T) {
	if !teatesting.RequireDBAvailable() {
		return
	}

	dao := NoticeDAO()
	count, err := dao.CountUnreadNoticesForAgent("local")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("count:", count)
}

func TestNoticeDAO_CountReadNoticesForAgent(t *testing.T) {
	if !teatesting.RequireDBAvailable() {
		return
	}

	dao := NoticeDAO()
	count, err := dao.CountReadNoticesForAgent("local")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("count:", count)
}

func TestNoticeDAO_CountReceivedNotices(t *testing.T) {
	if !teatesting.RequireDBAvailable() {
		return
	}

	dao := NoticeDAO()
	count, err := dao.CountReceivedNotices("EpBqPQMqpRlvFh9Q", map[string]interface{}{
		"agent.agentId": "local",
	}, 86400)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("count:", count)
}

func TestNoticeDAO_ExistNoticesWithHash(t *testing.T) {
	if !teatesting.RequireDBAvailable() {
		return
	}

	dao := NoticeDAO()
	{
		b, err := dao.ExistNoticesWithHash("4604200", map[string]interface{}{
			"agent.agentId": "local",
		}, 86400*time.Minute)
		if err != nil {
			t.Fatal(err)
		}
		if b {
			t.Log("exists")
		} else {
			t.Log("not exists")
		}
	}

	{
		b, err := dao.ExistNoticesWithHash("4157704579", map[string]interface{}{
			"agent.agentId": "local",
		}, 86400*time.Minute)
		if err != nil {
			t.Fatal(err)
		}
		if b {
			t.Log("exists")
		} else {
			t.Log("not exists")
		}
	}
}

func TestNoticeDAO_ListNotices(t *testing.T) {
	if !teatesting.RequireDBAvailable() {
		return
	}

	t.Log("===read===")
	{
		result, err := NoticeDAO().ListNotices(true, 0, 5)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(len(result), "notices")
		for _, n := range result {
			t.Log(n.Id, n.Message)
		}
	}

	t.Log("===unread===")
	{
		result, err := NoticeDAO().ListNotices(false, 0, 5)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(len(result), "notices")
		for _, n := range result {
			t.Log(n.Id, n.Message, n.Agent.AgentId, n.Agent.AppId)
		}
	}
}

func TestNoticeDAO_ListNotices_Unread(t *testing.T) {
	if !teatesting.RequireDBAvailable() {
		return
	}

	result, err := NoticeDAO().ListNotices(false, 0, 5)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(len(result), "notices")
	for _, n := range result {
		t.Log(n.Id, n.Message)
	}
}

func TestNoticeDAO_ListAgentNotices(t *testing.T) {
	if !teatesting.RequireDBAvailable() {
		return
	}

	result, err := NoticeDAO().ListAgentNotices("local", true, 0, 5)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(len(result), "notices")
	for _, n := range result {
		t.Log(n.Id, n.Message)
	}
}

func TestNoticeDAO_ListAgentNotices_Notfound(t *testing.T) {
	if !teatesting.RequireDBAvailable() {
		return
	}

	a := assert.NewAssertion(t)
	result, err := NoticeDAO().ListAgentNotices("local123", true, 0, 5)
	if err != nil {
		t.Fatal(err)
	}
	a.IsTrue(len(result) == 0)
}

func TestNoticeDAO_DeleteNoticesForAgent(t *testing.T) {
	if !teatesting.RequireDBAvailable() {
		return
	}

	dao := NoticeDAO()
	err := dao.DeleteNoticesForAgent("1TABzdF0uAIFPGkr")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("ok")
}

func TestNoticeDAO_UpdateNoticeReceivers(t *testing.T) {
	if !teatesting.RequireDBAvailable() {
		return
	}

	dao := NoticeDAO()
	ones, err := dao.ListNotices(false, 0, 1)
	if err != nil {
		t.Fatal(err)
	}
	if len(ones) == 0 {
		ones, err = dao.ListNotices(false, 0, 1)
		if err != nil {
			t.Fatal(err)
		}
		if len(ones) == 0 {
			t.Log("not found")
			return
		}
	}

	err = dao.UpdateNoticeReceivers(ones[0].Id.Hex(), []string{"a", "b", "c"})
	if err != nil {
		t.Fatal(err)
	}
	t.Log("ok")
}

func TestNoticeDAO_UpdateAllNoticesRead(t *testing.T) {
	if !teatesting.RequireDBAvailable() {
		return
	}

	dao := NoticeDAO()
	err := dao.UpdateAllNoticesRead()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("ok")
}

func TestNoticeDAO_UpdateNoticesRead(t *testing.T) {
	if !teatesting.RequireDBAvailable() {
		return
	}

	dao := NoticeDAO()
	ones, err := dao.ListNotices(false, 0, 2)
	if err != nil {
		t.Fatal(err)
	}
	if len(ones) == 0 {
		t.Log("ok")
		return
	}
	ids := []string{}
	for _, one := range ones {
		ids = append(ids, one.Id.Hex())
	}
	t.Log("ids:", ids)
	err = dao.UpdateNoticesRead(ids)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("ok")
}

func TestNoticeDAO_UpdateAgentNoticesRead(t *testing.T) {
	if !teatesting.RequireDBAvailable() {
		return
	}

	dao := NoticeDAO()
	ones, err := dao.ListAgentNotices("local", false, 0, 2)
	if err != nil {
		t.Fatal(err)
	}
	if len(ones) == 0 {
		t.Log("ok")
		return
	}
	ids := []string{}
	for _, one := range ones {
		ids = append(ids, one.Id.Hex())
	}
	t.Log("ids:", ids)
	err = dao.UpdateAgentNoticesRead("local", ids)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("ok")
}

func TestNoticeDAO_UpdateAllAgentNoticesRead(t *testing.T) {
	if !teatesting.RequireDBAvailable() {
		return
	}

	err := NoticeDAO().UpdateAllAgentNoticesRead("local")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("ok")
}
