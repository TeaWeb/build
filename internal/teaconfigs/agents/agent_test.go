package agents

import (
	"github.com/iwind/TeaGo/assert"
	"testing"
)

func TestAgentConfig_BelongsToGroup(t *testing.T) {
	a := assert.NewAssertion(t)

	{
		agent := NewAgentConfig()
		agent.GroupIds = []string{}
		a.IsTrue(agent.BelongsToGroup("default"))
		a.IsTrue(agent.BelongsToGroups([]string{"default"}))
		a.IsFalse(agent.BelongsToGroup("web001"))
	}

	{
		agent := NewAgentConfig()
		agent.GroupIds = []string{"web001"}
		a.IsFalse(agent.BelongsToGroup("default"))
		a.IsFalse(agent.BelongsToGroups([]string{"default"}))
		a.IsTrue(agent.BelongsToGroup("web001"))
	}

	{
		agent := NewAgentConfig()
		agent.GroupIds = []string{"default"}
		a.IsTrue(agent.BelongsToGroup(""))
		a.IsTrue(agent.BelongsToGroups([]string{"default"}))
		a.IsTrue(agent.BelongsToGroups([]string{}))
	}
}
