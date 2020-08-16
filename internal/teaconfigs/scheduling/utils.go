package scheduling

import "github.com/iwind/TeaGo/maps"

// 所有请求类型
func AllSchedulingTypes() []maps.Map {
	types := []maps.Map{}
	for _, s := range []SchedulingInterface{
		new(RandomScheduling),
		new(RoundRobinScheduling),
		new(HashScheduling),
		new(StickyScheduling),
	} {
		summary := s.Summary()
		summary["instance"] = s
		types = append(types, summary)
	}
	return types
}

func FindSchedulingType(code string) maps.Map {
	for _, summary := range AllSchedulingTypes() {
		if summary["code"] == code {
			return summary
		}
	}
	return nil
}
