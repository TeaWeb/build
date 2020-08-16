package widgets

import (
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/utils/string"
	"regexp"
	"strings"
)

var regexpNamedVariable = regexp.MustCompile(`\${([\w.\s-]+)}`)

// 线图
type LineChart struct {
	Params []string `yaml:"params" json:"params"` // deprecated: v0.1.8 使用Lines代替
	Lines  []*Line  `yaml:"lines" json:"lines"`
	Max    float64  `yaml:"max" json:"max"` // 最大值
}

// 添加线
func (this *LineChart) AddLine(line *Line) {
	this.Lines = append(this.Lines, line)
}

// 所有参数名
func (this *LineChart) AllParamNames() []string {
	result := []string{}
	for _, line := range this.Lines {
		for _, match := range regexpNamedVariable.FindAllStringSubmatch(line.Param, -1) {
			param := strings.Replace(match[1], " ", "", -1)
			param = strings.Replace(param, "\t", "", -1)
			if !lists.ContainsString(result, param) {
				result = append(result, param)
			}
		}
	}
	return result
}

// 转换为Javascript
func (this *LineChart) AsJavascript(options map[string]interface{}) (code string, err error) {
	if len(this.Lines) == 0 {
		this.Lines = []*Line{}
	}

	// 兼容老的版本
	if len(this.Params) > 0 {
		for _, param := range this.Params {
			line := NewLine()
			line.Param = param
			this.AddLine(line)
		}
	}

	options["lines"] = this.Lines
	options["params"] = this.AllParamNames()
	options["max"] = this.Max
	return `
var chart = new charts.LineChart();
chart.options = ` + stringutil.JSONEncode(options) + `;

if (chart.options.max != 0) {
	chart.max = chart.options.max;
}

var query = NewQuery();
var ones = query.past(60, time.MINUTE).avg.apply(query, chart.options.params);

var lines = [];

chart.options.lines.$each(function (k, v) {
	var line = new charts.Line();
	line.name = v.name;

	if (v.color == null || v.color.length == 0) {
		line.color = (k < colors.ARRAY.length) ? colors.ARRAY[k] : null;
	} else {
		line.color = colors[v.color]
	}
	line.isFilled = v.isFilled;
	line.values = [];
	lines.push(line);
});

ones.$each(function (k, v) {
	chart.options.lines.$each(function (k, lineOption) {
		var value = values.valueOf(v.value, lineOption.param);
		lines[k].values.push(value);

		if (k == 0) {
			chart.addLabel(v.label);
		}
	});
});

chart.addLines(lines);
chart.render();
`, nil
}
