package widgets

import "github.com/iwind/TeaGo/maps"

var AllChartTypes = []maps.Map{
	{
		"name":        "线图",
		"code":        "line",
		"description": "使用线图展示数据",
		"instance":    new(LineChart),
	},
	{
		"name":        "饼图",
		"code":        "pie",
		"description": "使用饼图展示数据",
		"instance":    new(PieChart),
	},
	{
		"name":        "HTML",
		"code":        "html",
		"description": "使用HTML写一个图表",
		"instance":    new(HTMLChart),
	},
	{
		"name":        "URL",
		"code":        "url",
		"description": "引入一个外部的URL，要注意可能产生的安全问题",
		"instance":    new(URLChart),
	},
	{
		"name":        "时钟",
		"code":        "clock",
		"description": "使用时钟展示当前时间，时间格式可以类似于 Mon Jan 21 16:46:06 CST 2019、2019-01-02 03:04:95 或一个时间戳 1548060582",
		"instance":    new(ClockChart),
	},
	{
		"name":        "Javascript",
		"code":        "javascript",
		"description": "直接使用Javascript代码来写图表",
		"instance":    new(JavascriptChart),
	},
}

var StatChartTypes = []maps.Map{
	{
		"name":        "线图",
		"code":        "line",
		"description": "使用线图展示数据",
		"instance":    new(LineChart),
	},
	{
		"name":        "饼图",
		"code":        "pie",
		"description": "使用饼图展示数据",
		"instance":    new(PieChart),
	},
	{
		"name":        "HTML",
		"code":        "html",
		"description": "使用HTML写一个图表",
		"instance":    new(HTMLChart),
	},
	{
		"name":        "URL",
		"code":        "url",
		"description": "引入一个外部的URL，要注意可能产生的安全问题",
		"instance":    new(URLChart),
	},
	{
		"name":        "时钟",
		"code":        "clock",
		"description": "使用时钟展示当前时间，时间格式可以类似于 Mon Jan 21 16:46:06 CST 2019、2019-01-02 03:04:95 或一个时间戳 1548060582",
		"instance":    new(ClockChart),
	},
	{
		"name":        "Javascript",
		"code":        "javascript",
		"description": "直接使用Javascript代码来写图表",
		"instance":    new(JavascriptChart),
	},
}

// 查找类型对应的名称
func FindChartTypeName(chartType string) string {
	for _, t := range AllChartTypes {
		if t["code"] == chartType {
			return t["name"].(string)
		}
	}
	return ""
}
