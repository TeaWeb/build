package agents

import (
	"math/rand"
	"testing"
	"time"
)

func TestSharedAgentIds(t *testing.T) {
	for i := 0; i < 10; i ++ {
		if rand.Int()%3 == 0 {
			agentListChanged = true
		}

		before := time.Now()
		t.Log(SharedAgents())
		t.Log(time.Since(before).Seconds(), "s")
	}
}
