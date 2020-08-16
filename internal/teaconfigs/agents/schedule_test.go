package agents

import (
	"github.com/iwind/TeaGo/assert"
	"github.com/iwind/TeaGo/utils/time"
	"testing"
	"time"
)

func TestScheduleConfig_Next(t *testing.T) {
	a := assert.NewAssertion(t)

	{
		s := NewScheduleConfig()
		s.AddYearRanges(&ScheduleRangeConfig{
			Value: 2018,
		})
		s.Validate()
		_, ok := s.Next(time.Now())
		a.IsFalse(ok)
	}

	{
		s := NewScheduleConfig()
		s.AddYearRanges(&ScheduleRangeConfig{
			Value: 2018,
		})
		s.AddYearRanges(&ScheduleRangeConfig{
			Value: 2019,
		})
		s.Validate()
		next, ok := s.Next(time.Now())
		a.IsTrue(ok)
		t.Log(timeutil.Format("Y-m-d H:i:s", next))
	}
}

func TestScheduleConfig_Next2(t *testing.T) {
	a := assert.NewAssertion(t)
	{
		s := NewScheduleConfig()
		s.AddYearRanges(&ScheduleRangeConfig{
			Value: 2018,
		})
		s.AddYearRanges(&ScheduleRangeConfig{
			Value: 2019,
		})
		s.AddYearRanges(&ScheduleRangeConfig{
			Value: 2020,
		})

		s.AddMonthRanges(&ScheduleRangeConfig{
			Value: 1,
		}, &ScheduleRangeConfig{
			Value: 5,
		}, )

		s.AddDayRanges(&ScheduleRangeConfig{
			Value: 9,
		}, &ScheduleRangeConfig{
			Value: -1,
			From:  1,
			To:    30,
			Step:  2,
		}, )

		s.AddHourRanges(&ScheduleRangeConfig{
			Value: 2,
		})
		s.AddHourRanges(&ScheduleRangeConfig{
			Value: 10,
		})

		s.AddMinuteRanges(&ScheduleRangeConfig{
			Every: true,
		})

		s.AddSecondRanges(&ScheduleRangeConfig{
			Value: 4,
		})

		s.Validate()
		next, ok := s.Next(time.Date(2019, 1, 9, 10, 0, 0, 0, time.Local))
		a.IsTrue(ok)
		t.Log(timeutil.Format("Y-m-d H:i:s", next))
	}
}

func TestScheduleConfig_Next3(t *testing.T) {
	a := assert.NewAssertion(t)
	{
		s := NewScheduleConfig()
		s.AddYearRanges(&ScheduleRangeConfig{
			Value: 2019,
		})
		s.AddMonthRanges(&ScheduleRangeConfig{
			Value: 1,
		}, &ScheduleRangeConfig{
			Value: 5,
		}, )

		s.AddDayRanges(&ScheduleRangeConfig{
			Value: 9,
		}, &ScheduleRangeConfig{
			Value: -1,
			From:  1,
			To:    30,
			Step:  2,
		}, )

		s.AddWeekDayRanges(&ScheduleRangeConfig{
			Value: -1,
			From:  4,
			To:    5,
			Step:  1,
		})

		s.AddHourRanges(&ScheduleRangeConfig{
			Value: 2,
			From:  -1,
		})
		s.AddHourRanges(&ScheduleRangeConfig{
			Value: 10,
			From:  -1,
		})

		s.AddMinuteRanges(&ScheduleRangeConfig{
			Every: true,
		})

		s.AddSecondRanges(&ScheduleRangeConfig{
			Value: 4,
		})

		s.Validate()
		next, ok := s.Next(time.Date(2019, 1, 9, 10, 0, 0, 0, time.Local))
		a.IsTrue(ok)
		t.Log(timeutil.Format("Y-m-d H:i:s w", next))
	}
}

func TestScheduleConfig_NextPerformance(t *testing.T) {
	s := NewScheduleConfig()
	s.AddYearRanges(&ScheduleRangeConfig{
		Value: 2018,
	})
	s.AddYearRanges(&ScheduleRangeConfig{
		Value: 2019,
	})
	s.AddYearRanges(&ScheduleRangeConfig{
		Value: 2020,
	})

	s.AddMonthRanges(&ScheduleRangeConfig{
		Value: 1,
	}, &ScheduleRangeConfig{
		Value: 5,
	}, )

	s.AddDayRanges(&ScheduleRangeConfig{
		Value: 9,
	}, &ScheduleRangeConfig{
		Value: -1,
		From:  1,
		To:    30,
		Step:  2,
	}, )

	s.AddHourRanges(&ScheduleRangeConfig{
		Value: 2,
	})
	s.AddHourRanges(&ScheduleRangeConfig{
		Value: 10,
	})

	s.AddMinuteRanges(&ScheduleRangeConfig{
		Every: true,
	})

	s.AddSecondRanges(&ScheduleRangeConfig{
		Value: 4,
	})
	s.Validate()

	count := 100 * 10000
	before := time.Now()
	for i := 0; i < count; i ++ {
		s.Next(time.Now())
	}
	t.Log(int(float64(count) / time.Since(before).Seconds()))
}
