package agents

import (
	"github.com/iwind/TeaGo/assert"
	"testing"
)

func TestScheduleRangeList_Next(t *testing.T) {
	a := assert.NewAssertion(t)

	r := &ScheduleRangeList{}
	r.Every = true
	a.IsTrue(r.Next(1, -1) == 1)

	r.Every = false
	a.IsTrue(r.Next(1, -1) == -1)

	r.Every = false
	r.Ranges = []*ScheduleRangeConfig{
		{
			Value: 1,
			From:  -1,
		},
		{
			Value: 2,
			From:  -1,
		},
		{
			Value: 3,
			From:  -1,
		},
		{
			Value: 5,
			From:  -1,
		},
		{
			Value: 7,
			From:  -1,
		},
	}
	a.IsTrue(r.Next(3, -1) == 3)
	a.IsTrue(r.Next(6, -1) == 7)
	a.IsTrue(r.Next(8, -1) == -1)

	r.Every = false
	r.Ranges = []*ScheduleRangeConfig{
		{
			From: 1, To: 4, Step: 1, Value: -1,
		},
	}
	a.IsTrue(r.Next(0, -1) == 1)
	a.IsTrue(r.Next(1, -1) == 1)
	a.IsTrue(r.Next(3, -1) == 3)

	r.Every = false
	r.Ranges = []*ScheduleRangeConfig{
		{
			From: 1, To: 5, Step: 2, Value: -1,
		},
	}
	a.IsTrue(r.Next(0, -1) == 1)
	a.IsTrue(r.Next(1, -1) == 1)
	a.IsTrue(r.Next(3, -1) == 3)
	a.IsTrue(r.Next(4, -1) == 5)
	a.IsTrue(r.Next(6, -1) == -1)

	r.Every = false
	r.Ranges = []*ScheduleRangeConfig{
		{
			From: 1, To: 5, Step: 2, Value: -1,
		},
		{
			From: 5, To: 10, Step: 2, Value: -1,
		},
	}
	a.IsTrue(r.Next(6, -1) == 7)
}
