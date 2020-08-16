package api

import (
	"github.com/TeaWeb/build/internal/teatesting"
	"github.com/iwind/TeaGo/assert"
	"testing"
	"time"
)

func TestAPITestPlan(t *testing.T) {
	if teatesting.IsGlobal() {
		return
	}

	a := assert.NewAssertion(t).Quiet()

	plan := NewAPITestPlan()
	now := time.Now()
	plan.Hour = now.Hour()
	plan.Minute = now.Minute()
	plan.Second = now.Second()
	plan.Weekdays = []int{1, 2, 3}

	t.Logf("%#v", plan)

	a.IsTrue(plan.MatchTime(time.Now()))
}
