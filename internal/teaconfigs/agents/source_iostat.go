package agents

import (
	"github.com/TeaWeb/build/internal/teaconfigs/forms"
	"github.com/TeaWeb/build/internal/teaconfigs/notices"
	"github.com/TeaWeb/build/internal/teaconfigs/widgets"
	"github.com/iwind/TeaGo/maps"
	"github.com/shirou/gopsutil/disk"
	"time"
)

// IO统计
type IOStatSource struct {
	Source `yaml:",inline"`

	lastStatMap map[string]disk.IOCountersStat
	lastTime    time.Time
}

// 获取新对象
func NewIOStatSource() *IOStatSource {
	return &IOStatSource{}
}

// 名称
func (this *IOStatSource) Name() string {
	return "IOStat"
}

// 代号
func (this *IOStatSource) Code() string {
	return "iostat"
}

// 描述
func (this *IOStatSource) Description() string {
	return "IO统计信息"
}

// 执行
func (this *IOStatSource) Execute(params map[string]string) (value interface{}, err error) {
	valueMap := maps.Map{}

	if this.lastStatMap == nil {
		this.lastStatMap = map[string]disk.IOCountersStat{}
	}

	statMap, err := disk.IOCounters()
	if err != nil {
		return valueMap, err
	}

	all := disk.IOCountersStat{
		Name: "ALL",
	}

	seconds := uint64(time.Now().Unix() - this.lastTime.Unix())
	if seconds == 0 {
		seconds = 1
	}

	for name, stat := range statMap {
		lastStat, ok := this.lastStatMap[name]
		if !ok {
			lastStat = stat
		}
		valueMap[name] = maps.Map{
			"readCount":           stat.ReadCount,
			"avgReadCount":        (stat.ReadCount - lastStat.ReadCount) / seconds,
			"mergedReadCount":     stat.MergedReadCount,
			"avgMergedReadCount":  (stat.MergedReadCount - lastStat.MergedReadCount) / seconds,
			"writeCount":          stat.WriteCount,
			"avgWriteCount":       (stat.WriteCount - lastStat.WriteCount) / seconds,
			"mergedWriteCount":    stat.MergedWriteCount,
			"avgMergedWriteCount": (stat.MergedWriteCount - lastStat.MergedWriteCount) / seconds,
			"readBytes":           stat.ReadBytes,
			"avgReadBytes":        (stat.ReadBytes - lastStat.ReadBytes) / seconds,
			"writeBytes":          stat.WriteBytes,
			"avgWriteBytes":       (stat.WriteBytes - lastStat.WriteBytes) / seconds,
			"readTime":            stat.ReadTime,
			"avgReadTime":         (stat.ReadTime - lastStat.ReadTime) / seconds,
			"writeTime":           stat.WriteTime,
			"avgWriteTime":        (stat.WriteTime - lastStat.WriteTime) / seconds,
			"ioTime":              stat.IoTime,
			"avgIOTime":           (stat.IoTime - lastStat.IoTime) / seconds,
			"name":                name,
		}

		this.lastStatMap[name] = stat

		all.ReadCount += stat.ReadCount
		all.MergedReadCount += stat.MergedReadCount
		all.WriteCount += stat.WriteCount
		all.MergedWriteCount += stat.MergedWriteCount
		all.ReadBytes += stat.ReadBytes
		all.WriteBytes += stat.WriteBytes
		all.ReadTime += stat.ReadTime
		all.WriteTime += stat.WriteTime
		all.IoTime += stat.IoTime
	}

	lastStat, ok := this.lastStatMap["ALL"]
	if !ok {
		lastStat = all
	}
	this.lastStatMap["ALL"] = all
	valueMap["ALL"] = maps.Map{
		"readCount":           all.ReadCount,
		"avgReadCount":        (all.ReadCount - lastStat.ReadCount) / seconds,
		"mergedReadCount":     all.MergedReadCount,
		"avgMergedReadCount":  (all.MergedReadCount - lastStat.MergedReadCount) / seconds,
		"writeCount":          all.WriteCount,
		"avgWriteCount":       (all.WriteCount - lastStat.WriteCount) / seconds,
		"mergedWriteCount":    all.MergedWriteCount,
		"avgMergedWriteCount": (all.MergedWriteCount - lastStat.MergedWriteCount) / seconds,
		"readBytes":           all.ReadBytes,
		"avgReadBytes":        (all.ReadBytes - lastStat.ReadBytes) / seconds,
		"writeBytes":          all.WriteBytes,
		"avgWriteBytes":       (all.WriteBytes - lastStat.WriteBytes) / seconds,
		"readTime":            all.ReadTime,
		"avgReadTime":         (all.ReadTime - lastStat.ReadTime) / seconds,
		"writeTime":           all.WriteTime,
		"avgWriteTime":        (all.WriteTime - lastStat.WriteTime) / seconds,
		"ioTime":              all.IoTime,
		"avgIOTime":           (all.IoTime - lastStat.IoTime) / seconds,
		"name":                "ALL",
	}

	this.lastTime = time.Now()

	return valueMap, nil
}

// 表单信息
func (this *IOStatSource) Form() *forms.Form {
	form := forms.NewForm(this.Code())
	return form
}

// 变量
func (this *IOStatSource) Variables() []*SourceVariable {
	return []*SourceVariable{
		{
			Code:        "ALL.readCount",
			Description: "读取总次数",
		},
		{
			Code:        "ALL.avgReadCount",
			Description: "每秒平均读取次数",
		},
		{
			Code:        "ALL.writeCount",
			Description: "写入总次数",
		},
		{
			Code:        "ALL.avgWriteCount",
			Description: "每秒平均写入次数",
		},
		{
			Code:        "ALL.mergedReadCount",
			Description: "合并后的读取次数",
		},
		{
			Code:        "ALL.avgMergedReadCount",
			Description: "合并后的每秒平均读取次数",
		},
		{
			Code:        "ALL.mergedWriteCount",
			Description: "合并后的写入次数",
		},
		{
			Code:        "ALL.avgMergedWriteCount",
			Description: "合并后的每秒平均写入次数",
		},
		{
			Code:        "ALL.readBytes",
			Description: "读取的总字节数",
		},
		{
			Code:        "ALL.avgReadBytes",
			Description: "平均每秒读取的字节数",
		},
		{
			Code:        "ALL.writeBytes",
			Description: "写入的总字节数",
		},
		{
			Code:        "ALL.avgWriteBytes",
			Description: "平均每秒写入的字节数",
		},
		{
			Code:        "ALL.readTime",
			Description: "读取的总时间",
		},
		{
			Code:        "ALL.avgReadTime",
			Description: "平均每秒读取的时间",
		},
		{
			Code:        "ALL.writeTime",
			Description: "写入的总时间",
		},
		{
			Code:        "ALL.avgWriteTime",
			Description: "平均每秒写入的时间",
		},
		{
			Code:        "ALL.ioTime",
			Description: "IO操作的总时间",
		},
		{
			Code:        "ALL.avgIOTime",
			Description: "平均每秒IO操作的时间",
		},
		{
			Code:        "ALL.name",
			Description: "总体统计，恒为ALL",
		},
		{
			Code:        "$.readCount",
			Description: "某个磁盘读取总次数",
		},
		{
			Code:        "$.avgReadCount",
			Description: "某个磁盘每秒平均读取次数",
		},
		{
			Code:        "$.writeCount",
			Description: "某个磁盘写入总次数",
		},
		{
			Code:        "$.avgWriteCount",
			Description: "某个磁盘每秒平均写入次数",
		},
		{
			Code:        "$.mergedReadCount",
			Description: "某个磁盘合并后的读取次数",
		},
		{
			Code:        "$.avgMergedReadCount",
			Description: "某个磁盘合并后的每秒平均读取次数",
		},
		{
			Code:        "$.mergedWriteCount",
			Description: "某个磁盘合并后的写入次数",
		},
		{
			Code:        "$.avgMergedWriteCount",
			Description: "某个磁盘合并后的每秒平均写入次数",
		},
		{
			Code:        "$.readBytes",
			Description: "某个磁盘读取的总字节数",
		},
		{
			Code:        "$.avgReadBytes",
			Description: "某个磁盘平均每秒读取的字节数",
		},
		{
			Code:        "$.writeBytes",
			Description: "某个磁盘写入的总字节数",
		},
		{
			Code:        "$.avgWriteBytes",
			Description: "某个磁盘平均每秒写入的字节数",
		},
		{
			Code:        "$.readTime",
			Description: "某个磁盘读取的总时间",
		},
		{
			Code:        "$.avgReadTime",
			Description: "某个磁盘平均每秒读取的时间",
		},
		{
			Code:        "$.writeTime",
			Description: "某个磁盘写入的总时间",
		},
		{
			Code:        "$.avgWriteTime",
			Description: "某个磁盘平均每秒写入的时间",
		},
		{
			Code:        "$.ioTime",
			Description: "某个磁盘IO操作的总时间",
		},
		{
			Code:        "$.avgIOTime",
			Description: "某个磁盘平均每秒IO操作的时间",
		},
		{
			Code:        "$.name",
			Description: "磁盘卷名",
		},
	}
}

// 阈值
func (this *IOStatSource) Thresholds() []*Threshold {
	result := []*Threshold{}
	{
		t := NewThreshold()
		t.Param = "${ALL.avgReadBytes}"
		t.Value = "20971520" // 20M
		t.NoticeLevel = notices.NoticeLevelWarning
		t.Operator = ThresholdOperatorGte
		t.NoticeMessage = "IO读取速度超过每秒20M"
		result = append(result, t)
	}
	{
		t := NewThreshold()
		t.Param = "${ALL.avgWriteBytes}"
		t.Value = "20971520" // 20M
		t.NoticeLevel = notices.NoticeLevelWarning
		t.Operator = ThresholdOperatorGte
		t.NoticeMessage = "IO写入速度超过每秒20M"
		result = append(result, t)
	}
	return result
}

// 图表
func (this *IOStatSource) Charts() []*widgets.Chart {
	charts := []*widgets.Chart{}

	{
		// chart
		chart := widgets.NewChart()
		chart.Id = "disk.stat"
		chart.Name = "IO统计"
		chart.Columns = 2
		chart.Type = "javascript"
		chart.SupportsTimeRange = true
		chart.Options = maps.Map{
			"code": `var chart = new charts.LineChart();

var ones = NewQuery().past(60, time.MINUTE).avg("ALL.avgReadBytes", "ALL.avgWriteBytes");

var lines = [];

{
	var line = new charts.Line();
	line.name = "读（MBytes）";
	line.color = colors.ARRAY[0];
	line.isFilled = true;
	lines.push(line);
}

{
	var line = new charts.Line();
	line.name = "写（MBytes）";
	line.color = colors.BROWN;
	line.isFilled = true;
	lines.push(line);
}

ones.$each(function (k, v) {
	lines[0].addValue(v.value.ALL.avgReadBytes / 1024 / 1024 );
	lines[1].addValue(v.value.ALL.avgWriteBytes / 1024 / 1024 );
	chart.addLabel(v.label);
});

chart.addLines(lines);
chart.render();
`,
		}

		charts = append(charts, chart)
	}

	return charts
}
