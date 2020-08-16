package teaconfigs

import (
	"github.com/iwind/TeaGo/maps"
	"time"
)

// 常用时间范围
type TimePast = string

const (
	TimePast5m  TimePast = "past5m"
	TimePast15m TimePast = "past15m"
	TimePast1h  TimePast = "past1h"
	TimePast6h  TimePast = "past6h"
	TimePast12h TimePast = "past12h"
	TimePast24h TimePast = "past24h"
	TimePast1d  TimePast = "past1d"
	TimePast2d  TimePast = "past2d"
	TimePast7d  TimePast = "past7d"
	TimePast30d TimePast = "past30d"
)

// 时间单位
type TimeUnit = string

const (
	TimeUnitSecond TimeUnit = "SECOND"
	TimeUnitMinute TimeUnit = "MINUTE"
	TimeUnitHour   TimeUnit = "HOUR"
	TimeUnitDay    TimeUnit = "DAY"
	TimeUnitMonth  TimeUnit = "MONTH"
	TimeUnitYear   TimeUnit = "YEAR"
)

func AllTimePasts() []maps.Map {
	return []maps.Map{
		{
			"name":  "5分钟内",
			"value": TimePast5m,
			"unit":  TimeUnitSecond,
		},
		{
			"name":  "15分钟内",
			"value": TimePast15m,
			"unit":  TimeUnitMinute,
		},
		{
			"name":  "1小时内",
			"value": TimePast1h,
			"unit":  TimeUnitMinute,
		},
		{
			"name":  "6小时内",
			"value": TimePast6h,
			"unit":  TimeUnitHour,
		},
		{
			"name":  "12小时内",
			"value": TimePast12h,
			"unit":  TimeUnitHour,
		},
		{
			"name":  "24小时内",
			"value": TimePast24h,
			"unit":  TimeUnitHour,
		},
		{
			"name":  "当天",
			"value": TimePast1d,
			"unit":  TimeUnitHour,
		},
		{
			"name":  "2天内",
			"value": TimePast2d,
			"unit":  TimeUnitHour,
		},
		{
			"name":  "7天内",
			"value": TimePast7d,
			"unit":  TimeUnitDay,
		},
		{
			"name":  "30天内",
			"value": TimePast30d,
			"unit":  TimeUnitDay,
		},
	}
}

func TimePastUnixTime(past TimePast) (timestamp int64) {
	now := time.Now()
	switch past {
	case TimePast5m:
		return now.Unix() - 5*60
	case TimePast15m:
		return now.Unix() - 15*60
	case TimePast1h:
		return now.Unix() - 3600
	case TimePast6h:
		return now.Unix() - 6*3600
	case TimePast12h:
		return now.Unix() - 12*3600
	case TimePast24h:
		return now.Unix() - 24*3600
	case TimePast1d:
		return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).Unix()
	case TimePast2d:
		return time.Date(now.Year(), now.Month(), now.Day()-1, 0, 0, 0, 0, now.Location()).Unix()
	case TimePast7d:
		return time.Date(now.Year(), now.Month(), now.Day()-6, 0, 0, 0, 0, now.Location()).Unix()
	case TimePast30d:
		return time.Date(now.Year(), now.Month(), now.Day()-29, 0, 0, 0, 0, now.Location()).Unix()
	}
	return 0
}

func TimePastUnit(past TimePast) TimeUnit {
	for _, p := range AllTimePasts() {
		if p.GetString("value") == past {
			return p.GetString("unit")
		}
	}
	return TimeUnitMinute
}

func TimePastUnixTimeWithUnit(number int64, unit TimeUnit) int64 {
	now := time.Now()

	switch unit {
	case TimeUnitSecond:
		return now.Unix() - number
	case TimeUnitMinute:
		return now.Unix() - number*60
	case TimeUnitHour:
		return now.Unix() - number*3600
	case TimeUnitDay:
		return time.Date(now.Year(), now.Month(), now.Day()-int(number), 0, 0, 0, 0, now.Location()).Unix()
	case TimeUnitMonth:
		return time.Date(now.Year(), now.Month()-time.Month(number), 1, 0, 0, 0, 0, now.Location()).Unix()
	case TimeUnitYear:
		return time.Date(now.Year()-int(number), 1, 1, 0, 0, 0, 0, now.Location()).Unix()
	}
	return 0
}
