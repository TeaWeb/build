package teastats

import (
	"fmt"
	"testing"
)

func TestFindNewFilter(t *testing.T) {
	t.Log(FindNewFilter("request.all.second"))
}

func TestRestartServerFilters(t *testing.T) {
	RestartServerFilters("123456", []string{"request.all.second", "request.all.minute", "request.all.minute"})
	RestartServerFilters("123456", []string{"request.all.second", "request.all.minute", "pv.all.second"})
}

func TestMarkdownFilters(t *testing.T) {
	for _, filter := range AllStatFilters {
		name := filter.GetString("name")
		period := filter.GetString("period")
		if len(period) > 0 {
			name = name + "（" + period + "）"
		}
		fmt.Println("### " + name + " - " + filter.GetString("code"))
		fmt.Println(filter.GetString("description") + "。")
		fmt.Println("")

		params := filter.Get("instance").(FilterInterface).ParamVariables()
		if len(params) == 0 {
			fmt.Println("参数：无。")
		} else {
			fmt.Println("参数：")
			for _, param := range params {
				fmt.Println("* `" + param.Code + "` - " + param.Description)
			}
		}

		fmt.Println("")

		values := filter.Get("instance").(FilterInterface).ValueVariables()
		if len(values) == 0 {
			fmt.Println("统计数据：无。")
		} else {
			fmt.Println("统计数据：")
			for _, param := range values {
				fmt.Println("* `" + param.Code + "` - " + param.Description)
			}
		}

		fmt.Println("")
	}
}
