package teadb

import (
	"github.com/TeaWeb/build/internal/teaconfigs/agents"
	"github.com/TeaWeb/build/internal/teatesting"
	"github.com/iwind/TeaGo/assert"
	"github.com/iwind/TeaGo/rands"
	stringutil "github.com/iwind/TeaGo/utils/string"
	"testing"
	"time"
)

func TestAgentLogDAO_InsertOne(t *testing.T) {
	if !teatesting.RequireDBAvailable() {
		return
	}

	{
		log := new(agents.ProcessLog)
		log.AgentId = "test"
		log.Data = "abcdefg"
		log.TaskId = "abc"
		log.SetTime(time.Now())
		log.EventType = "start"
		log.ProcessId = rands.HexString(16)
		log.ProcessPid = 1024
		err := AgentLogDAO().InsertOne("test", log)
		if err != nil {
			t.Fatal(err)
		}
		t.Log("ok")
	}

	{
		log := new(agents.ProcessLog)
		log.AgentId = "test"
		log.Data = "abcdefg1"
		log.TaskId = "abc"
		log.EventType = "run"
		log.ProcessPid = 1025
		log.SetTime(time.Now())
		err := AgentLogDAO().InsertOne("test", log)
		if err != nil {
			t.Fatal(err)
		}
		t.Log("ok")
	}

	{
		log := new(agents.ProcessLog)
		log.AgentId = "test"
		log.Data = "abcdefg2"
		log.TaskId = "abc"
		log.EventType = "run"
		log.ProcessPid = 1026
		log.SetTime(time.Now())
		err := AgentLogDAO().InsertOne("test", log)
		if err != nil {
			t.Fatal(err)
		}
		t.Log("ok")
	}
}

func TestAgentLogDAO_ListTaskLogs(t *testing.T) {
	if !teatesting.RequireDBAvailable() {
		return
	}

	a := assert.NewAssertion(t)

	taskId := "abc"
	taskLogs, err := AgentLogDAO().FindLatestTaskLogs("test", taskId, "", 2)
	if err != nil {
		t.Fatal(err)
	}

	for _, log := range taskLogs {
		a.IsTrue(log.TaskId == taskId)
		t.Log(log.Id, log.TaskId, log.Data)
	}

	if len(taskLogs) > 0 {
		t.Log("=======")
		taskLogs, err := AgentLogDAO().FindLatestTaskLogs("test", taskId, taskLogs[len(taskLogs)-1].Id.Hex(), 2)
		if err != nil {
			t.Fatal(err)
		}

		for _, log := range taskLogs {
			a.IsTrue(log.TaskId == taskId)
			t.Log(log.Id, log.TaskId, log.Data)
		}
	}
}

func TestAgentLogDAO_FindLatestTaskLog(t *testing.T) {
	if !teatesting.RequireDBAvailable() {
		return
	}

	a := assert.NewAssertion(t)

	taskId := "abc"
	taskLog, err := AgentLogDAO().FindLatestTaskLog("test", taskId)
	if err != nil {
		t.Fatal(err)
	}
	if taskLog == nil {
		t.Log("not found")
		return
	}
	a.IsTrue(taskLog.TaskId == taskId)
	t.Log(stringutil.JSONEncodePretty(taskLog))
}
