package agents

import (
	"github.com/iwind/TeaGo/assert"
	"testing"
	"time"
)

func TestTaskConfig_Schedule(t *testing.T) {
	a := assert.NewAssertion(t)

	task := NewTaskConfig()

	{
		schedule := NewScheduleConfig()
		schedule.AddSecondRanges(&ScheduleRangeConfig{
			Every: true,
		})
		task.Schedule = []*ScheduleConfig{schedule}

		task.Validate()

		now, ok := task.Next(time.Date(2019, 2, 1, 10, 20, 30, 0, time.Local))
		a.IsTrue(ok).Then(func() {
			t.Log(now)
		})
	}

	{
		schedule := NewScheduleConfig()
		schedule.AddSecondRanges(&ScheduleRangeConfig{
			Every: false,
			From:  0,
			To:    59,
			Step:  4,
		})
		task.Schedule = []*ScheduleConfig{schedule}
		task.Validate()

		now, ok := task.Next(time.Date(2019, 2, 1, 10, 20, 30, 0, time.Local))
		a.IsTrue(ok).Then(func() {
			a.IsTrue(now.Second() == 32)
			t.Log(now)
		})
	}
}

func TestTaskConfig_ScheduleNextHour(t *testing.T) {
	a := assert.NewAssertion(t)

	task := NewTaskConfig()

	{
		schedule := NewScheduleConfig()
		schedule.AddMinuteRanges(&ScheduleRangeConfig{
			From:  1,
			To:    10,
			Step:  1,
			Value: -1,
		})
		schedule.AddSecondRanges(&ScheduleRangeConfig{
			Every: false,
			From:  2,
			To:    59,
			Step:  4,
			Value: -1,
		})
		task.Schedule = []*ScheduleConfig{schedule}
		task.Validate()

		now, ok := task.Next(time.Date(2019, 2, 1, 10, 20, 30, 0, time.Local))
		a.IsTrue(ok).Then(func() {
			t.Log(now)
		})
	}

	{
		schedule := NewScheduleConfig()
		schedule.AddMinuteRanges(&ScheduleRangeConfig{
			From:  1,
			To:    10,
			Step:  1,
			Value: -1,
		})
		schedule.AddSecondRanges(&ScheduleRangeConfig{
			Every: false,
			From:  2,
			To:    59,
			Step:  4,
			Value: -1,
		})
		task.Schedule = []*ScheduleConfig{schedule}
		task.Validate()

		now, ok := task.Next(time.Date(2019, 2, 1, 10, 9, 33, 0, time.Local))
		a.IsTrue(ok).Then(func() {
			t.Log(now)
		})
	}
}

func TestTaskConfig_ScheduleNextMinute(t *testing.T) {
	a := assert.NewAssertion(t)

	task := NewTaskConfig()

	{
		schedule := NewScheduleConfig()
		schedule.AddMinuteRanges(&ScheduleRangeConfig{
			From:  1,
			To:    20,
			Step:  1,
			Value: -1,
		})
		schedule.AddSecondRanges(&ScheduleRangeConfig{
			Every: false,
			From:  2,
			To:    20,
			Step:  4,
			Value: -1,
		})
		task.Schedule = []*ScheduleConfig{schedule}
		task.Validate()

		now, ok := task.Next(time.Date(2019, 2, 1, 10, 20, 30, 0, time.Local))
		a.IsTrue(ok).Then(func() {
			t.Log(now)
		})
	}
}

func TestTaskConfig_Next(t *testing.T) {
	a := assert.NewAssertion(t)

	task := NewTaskConfig()

	{
		schedule := NewScheduleConfig()
		schedule.AddYearRanges(&ScheduleRangeConfig{
			From:  -1,
			Value: 2019,
		})
		schedule.AddMonthRanges(&ScheduleRangeConfig{
			From:  -1,
			Value: 2,
		})
		schedule.AddDayRanges(&ScheduleRangeConfig{
			From:  -1,
			Value: 3,
		})
		schedule.AddHourRanges(&ScheduleRangeConfig{
			From:  0,
			To:    10,
			Value: -1,
		})
		schedule.AddMinuteRanges(&ScheduleRangeConfig{
			From:  2,
			To:    20,
			Step:  2,
			Value: -1,
		})
		schedule.AddSecondRanges(&ScheduleRangeConfig{
			Every: false,
			From:  -1,
			Value: 22,
		})
		schedule.AddSecondRanges(&ScheduleRangeConfig{
			Every: false,
			From:  5,
			To:    30,
			Value: -1,
		})
		task.Schedule = []*ScheduleConfig{schedule}
		task.Validate()

		now, ok := task.Next(time.Date(2019, 2, 3, 10, 20, 30, 0, time.Local))
		a.IsTrue(ok).Then(func() {
			t.Log(now)
		})
	}
}
