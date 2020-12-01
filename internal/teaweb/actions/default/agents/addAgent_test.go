package agents

import (
	"fmt"
	"github.com/TeaWeb/build/internal/teaconfigs/agents"
	"github.com/TeaWeb/build/internal/teatesting"
	"github.com/iwind/TeaGo/rands"
	"testing"
)

func TestAddManyAgents(t *testing.T) {
	if teatesting.IsGlobal() {
		return
	}

	count := 1000
	for i := 0; i < count; i++ {
		agentList, err := agents.SharedAgentList()
		if err != nil {
			t.Fatal(err)
		}

		agent := agents.NewAgentConfig()
		agent.On = false
		agent.Name = "Web" + fmt.Sprintf("%d", i)
		agent.Host = "192.168.0." + fmt.Sprintf("%d", i)
		agent.AllowAll = true
		agent.Allow = []string{}
		agent.Key = rands.HexString(32)
		//agent.GroupIds = []string{"2kMMzOcWWPFrhdaM"}
		err = agent.Save()
		if err != nil {
			t.Fatal(err)
		}

		agentList.AddAgent(agent.Filename())
		err = agentList.Save()
		if err != nil {
			t.Fatal(err)
		}
	}
}
