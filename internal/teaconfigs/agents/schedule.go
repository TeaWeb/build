package agents

import (
	"fmt"
	"strings"
	"time"
)

// 定时
type ScheduleConfig struct {
	SecondRanges  []*ScheduleRangeConfig `yaml:"secondRanges" json:"secondRanges"`   // 秒
	MinuteRanges  []*ScheduleRangeConfig `yaml:"minuteRanges" json:"minuteRanges"`   // 分
	HourRanges    []*ScheduleRangeConfig `yaml:"hourRanges" json:"hourRanges"`       // 小时
	DayRanges     []*ScheduleRangeConfig `yaml:"dayRanges" json:"dayRanges"`         // 天
	MonthRanges   []*ScheduleRangeConfig `yaml:"monthRanges" json:"monthRanges"`     // 月
	YearRanges    []*ScheduleRangeConfig `yaml:"yearRanges" json:"yearRanges"`       // 年份
	WeekDayRanges []*ScheduleRangeConfig `yaml:"weekDayRanges" json:"weekDayRanges"` // 一周中的某天，1-7

	secondList  *ScheduleRangeList
	minuteList  *ScheduleRangeList
	hourList    *ScheduleRangeList
	dayList     *ScheduleRangeList
	monthList   *ScheduleRangeList
	yearList    *ScheduleRangeList
	weekDayList *ScheduleRangeList
}

// 定时新对象
func NewScheduleConfig() *ScheduleConfig {
	return &ScheduleConfig{}
}

// 校验
func (this *ScheduleConfig) Validate() error {
	this.secondList = NewScheduleRangeList()
	this.secondList.Convert(this.SecondRanges)

	this.minuteList = NewScheduleRangeList()
	this.minuteList.Convert(this.MinuteRanges)

	this.hourList = NewScheduleRangeList()
	this.hourList.Convert(this.HourRanges)

	this.dayList = NewScheduleRangeList()
	this.dayList.Convert(this.DayRanges)

	this.monthList = NewScheduleRangeList()
	this.monthList.Convert(this.MonthRanges)

	this.yearList = NewScheduleRangeList()
	this.yearList.Convert(this.YearRanges)

	this.weekDayList = NewScheduleRangeList()
	this.weekDayList.Convert(this.WeekDayRanges)
	return nil
}

// 添加秒设置
func (this *ScheduleConfig) AddSecondRanges(r ... *ScheduleRangeConfig) {
	this.SecondRanges = append(this.SecondRanges, r ...)
}

// 添加分钟设置
func (this *ScheduleConfig) AddMinuteRanges(r ... *ScheduleRangeConfig) {
	this.MinuteRanges = append(this.MinuteRanges, r...)
}

// 添加小时设置
func (this *ScheduleConfig) AddHourRanges(r ... *ScheduleRangeConfig) {
	this.HourRanges = append(this.HourRanges, r...)
}

// 添加天设置
func (this *ScheduleConfig) AddDayRanges(r ... *ScheduleRangeConfig) {
	this.DayRanges = append(this.DayRanges, r...)
}

// 添加月设置
func (this *ScheduleConfig) AddMonthRanges(r ... *ScheduleRangeConfig) {
	this.MonthRanges = append(this.MonthRanges, r...)
}

// 添加年设置
func (this *ScheduleConfig) AddYearRanges(r ... *ScheduleRangeConfig) {
	this.YearRanges = append(this.YearRanges, r...)
}

// 添加周N设置
func (this *ScheduleConfig) AddWeekDayRanges(r ... *ScheduleRangeConfig) {
	this.WeekDayRanges = append(this.WeekDayRanges, r...)
}

// Summary
func (this *ScheduleConfig) Summary() string {
	var summaryStrings []string
	var secondStrings []string

	if len(this.SecondRanges) == 0 {
		secondStrings = append(secondStrings, "每秒钟")
	} else {
		for _, r := range this.SecondRanges {
			if r.Every {
				secondStrings = append(secondStrings, "每秒钟")
			}
		}

		for _, r := range this.SecondRanges {
			if r.Value > -1 {
				secondStrings = append(secondStrings, fmt.Sprintf("%d秒", r.Value))
			}
		}
		for _, r := range this.SecondRanges {
			if r.From > -1 && r.To > -1 && r.Step > -1 {
				secondStrings = append(secondStrings, fmt.Sprintf("%d秒-%d秒/每%d秒", r.From, r.To, r.Step))
			}
		}
	}

	if len(secondStrings) > 0 {
		summaryStrings = append(summaryStrings, "["+strings.Join(secondStrings, "，")+"]")
	}

	var minuteStrings []string
	if len(this.MinuteRanges) == 0 {
		minuteStrings = append(minuteStrings, "每分钟")
	} else {
		for _, r := range this.MinuteRanges {
			if r.Every {
				minuteStrings = append(minuteStrings, "每分钟")
			}
		}
		for _, r := range this.MinuteRanges {
			if r.Value > -1 {
				minuteStrings = append(minuteStrings, fmt.Sprintf("%d分", r.Value))
			}
		}
		for _, r := range this.MinuteRanges {
			if r.From > -1 && r.To > -1 && r.Step > -1 {
				minuteStrings = append(minuteStrings, fmt.Sprintf("%d分钟-%d分钟/每%d分钟", r.From, r.To, r.Step))
			}
		}
	}
	if len(minuteStrings) > 0 {
		summaryStrings = append(summaryStrings, "["+strings.Join(minuteStrings, "，")+"]")
	}

	var hourStrings []string
	if len(this.HourRanges) == 0 {
		hourStrings = append(hourStrings, "每小时")
	} else {
		for _, r := range this.HourRanges {
			if r.Every {
				hourStrings = append(hourStrings, "每小时")
			}
		}
		for _, r := range this.HourRanges {
			if r.Value > -1 {
				hourStrings = append(hourStrings, fmt.Sprintf("%d小时", r.Value))
			}
		}
		for _, r := range this.HourRanges {
			if r.From > -1 && r.To > -1 && r.Step > -1 {
				hourStrings = append(hourStrings, fmt.Sprintf("%d小时-%d小时/每%d小时", r.From, r.To, r.Step))
			}
		}
	}
	if len(hourStrings) > 0 {
		summaryStrings = append(summaryStrings, "["+strings.Join(hourStrings, "，")+"]")
	}

	var dayStrings []string
	for _, r := range this.DayRanges {
		if r.Value > -1 {
			dayStrings = append(dayStrings, fmt.Sprintf("%d日", r.Value))
		}
	}
	for _, r := range this.DayRanges {
		if r.From > -1 && r.To > -1 && r.Step > -1 {
			dayStrings = append(dayStrings, fmt.Sprintf("%d日-%d日/每%d天", r.From, r.To, r.Step))
		}
	}
	if len(dayStrings) > 0 {
		summaryStrings = append(summaryStrings, "["+strings.Join(dayStrings, "，")+"]")
	}

	var monthStrings []string
	for _, r := range this.MonthRanges {
		if r.Value > -1 {
			monthStrings = append(monthStrings, fmt.Sprintf("%d月", r.Value))
		}
	}
	for _, r := range this.MonthRanges {
		if r.From > -1 && r.To > -1 && r.Step > -1 {
			monthStrings = append(monthStrings, fmt.Sprintf("%d月-%d月/每%d月", r.From, r.To, r.Step))
		}
	}
	if len(monthStrings) > 0 {
		summaryStrings = append(summaryStrings, "["+strings.Join(monthStrings, "，")+"]")
	}

	var yearStrings []string
	for _, r := range this.YearRanges {
		if r.Value > -1 {
			yearStrings = append(yearStrings, fmt.Sprintf("%d年", r.Value))
		}
	}
	for _, r := range this.YearRanges {
		if r.From > -1 && r.To > -1 && r.Step > -1 {
			yearStrings = append(yearStrings, fmt.Sprintf("%d年-%d年/每%d年", r.From, r.To, r.Step))
		}
	}
	if len(yearStrings) > 0 {
		summaryStrings = append(summaryStrings, "["+strings.Join(yearStrings, "，")+"]")
	}

	var weekDayStrings []string
	for _, r := range this.WeekDayRanges {
		if r.Value > -1 {
			weekDayStrings = append(weekDayStrings, fmt.Sprintf("周%d", r.Value))
		}
	}
	for _, r := range this.WeekDayRanges {
		if r.From > -1 && r.To > -1 && r.Step > -1 {
			weekDayStrings = append(weekDayStrings, fmt.Sprintf("周%d-周%d/每%d天", r.From, r.To, r.Step))
		}
	}

	if len(weekDayStrings) > 0 {
		summaryStrings = append(summaryStrings, "["+strings.Join(weekDayStrings, "，")+"]")
	}

	return strings.Join(summaryStrings, " ")
}

// 下次运行时间
func (this *ScheduleConfig) Next(now time.Time) (t time.Time, ok bool) {
	var year, month, day, hour, minute, second int

	// weekday
	currentWeekday := int(now.Weekday())
	if currentWeekday == 0 {
		currentWeekday = 7
	}
	weekday := this.weekDayList.Next(currentWeekday, -1)
	if weekday == - 1 {
		return
	}
	if weekday != currentWeekday {
		return this.Next(time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location()))
	}

	year = this.yearList.Next(now.Year(), -1)
	if year == -1 {
		return
	}

	if year == now.Year() {
		month = this.monthList.Next(int(now.Month()), -1)
		if month == -1 {
			now = time.Date(now.Year()+1, 1, 1, 0, 0, 0, 0, now.Location())
			return this.Next(now)
		}
	} else {
		month = this.monthList.Next(1, -1)
	}
	if month == -1 {
		return
	}

	// day
	if year == now.Year() && month == int(now.Month()) {
		day = this.dayList.Next(now.Day(), -1)
		if day == -1 {
			now = time.Date(now.Year(), now.Month()+1, 1, 0, 0, 0, 0, now.Location())
			return this.Next(now)
		}
	} else {
		day = this.dayList.Next(1, -1)
	}
	if day == -1 {
		return
	}

	// hour
	if year == now.Year() && month == int(now.Month()) && day == now.Day() {
		hour = this.hourList.Next(now.Hour(), -1)
		if hour == -1 {
			now = time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
			return this.Next(now)
		}
	} else {
		hour = this.hourList.Next(0, -1)
	}
	if hour == -1 {
		return
	}

	if year == now.Year() && month == int(now.Month()) && day == now.Day() && hour == now.Hour() {
		minute = this.minuteList.Next(now.Minute(), -1)
		if minute == -1 {
			now = time.Date(now.Year(), now.Month(), now.Day(), now.Hour()+1, 0, 0, 0, now.Location())
			return this.Next(now)
		}
	} else {

		minute = this.minuteList.Next(0, -1)
	}
	if minute == -1 {
		return
	}

	// second
	if year == now.Year() && month == int(now.Month()) && day == now.Day() && hour == now.Hour() && minute == now.Minute() {
		second = this.secondList.Next(now.Second(), -1)
		if second == -1 {
			now = time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute()+1, 0, 0, now.Location())
			return this.Next(now)
		}
	} else {
		second = this.secondList.Next(0, -1)
	}
	if second == -1 {
		return
	}

	t = time.Date(year, time.Month(month), day, hour, minute, second, 0, now.Location())
	ok = true
	return
}

// 判断时间是否匹配
func (this *ScheduleConfig) ShouldRun(t time.Time) bool {
	return false
}
