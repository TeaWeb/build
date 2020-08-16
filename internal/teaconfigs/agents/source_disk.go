package agents

import (
	"github.com/TeaWeb/build/internal/teaconfigs/forms"
	"github.com/TeaWeb/build/internal/teaconfigs/notices"
	"github.com/TeaWeb/build/internal/teaconfigs/widgets"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"github.com/shirou/gopsutil/disk"
	"runtime"
	"strings"
)

// 文件系统信息
type DiskSource struct {
	Source `yaml:",inline"`

	ContainsAllMountPoints bool `yaml:"containsAllMountPoints" json:"containsAllMountPoints"`
}

// 获取新对象
func NewDiskSource() *DiskSource {
	return &DiskSource{}
}

// 名称
func (this *DiskSource) Name() string {
	return "文件系统信息"
}

// 代号
func (this *DiskSource) Code() string {
	return "disk"
}

// 描述
func (this *DiskSource) Description() string {
	return "文件系统信息"
}

// 执行
func (this *DiskSource) Execute(params map[string]string) (value interface{}, err error) {
	partitions, err := disk.Partitions(false)
	if err != nil {
		logs.Error(err)
		return
	}
	lists.Sort(partitions, func(i int, j int) bool {
		p1 := partitions[i]
		p2 := partitions[j]
		return p1.Mountpoint > p2.Mountpoint
	})

	// 当前TeaWeb所在的fs
	rootFS := ""
	if lists.ContainsString([]string{"darwin", "linux", "freebsd"}, runtime.GOOS) {
		for _, p := range partitions {
			if p.Mountpoint == "/" {
				rootFS = p.Fstype
				break
			}
		}
	}

	result := []maps.Map{}
	for _, partition := range partitions {
		if runtime.GOOS != "windows" && !strings.Contains(partition.Device, "/") && !strings.Contains(partition.Device, "\\") {
			continue
		}

		// 跳过不同fs的
		if !this.ContainsAllMountPoints {
			if len(rootFS) > 0 && rootFS != partition.Fstype {
				continue
			}
		}

		usage, err := disk.Usage(partition.Mountpoint)
		if err != nil {
			continue
		}
		result = append(result, maps.Map{
			"name":          partition.Mountpoint,
			"used":          usage.Used,
			"free":          usage.Free,
			"total":         usage.Total,
			"percent":       usage.UsedPercent,
			"inodesUsed":    usage.InodesUsed,
			"inodesFree":    usage.InodesFree,
			"inodesTotal":   usage.InodesTotal,
			"inodesPercent": usage.InodesUsedPercent,
			"fstype":        usage.Fstype,
		})
	}

	value = maps.Map{
		"partitions": result,
	}

	return
}

// 表单信息
func (this *DiskSource) Form() *forms.Form {
	form := forms.NewForm(this.Code())
	{
		group := form.NewGroup()

		{
			field := forms.NewCheckBox("包含所有挂载点", "")
			field.Value = "1"
			field.Comment = "如果没有选中，则会自动去掉不同文件系统类型的一些挂载点"
			field.Code = "containsAllMountPoints"

			group.Add(field)
		}
	}
	return form
}

// 显示信息
func (this *DiskSource) Presentation() *forms.Presentation {
	p := forms.NewPresentation()
	p.HTML = `
<tr>
	<td>包含所有挂载点</td>
	<td>
		<span class="green" v-if="source.containsAllMountPoints">Y</span>
		<span v-if="!source.containsAllMountPoints">N</span>
	</td>
</tr>
`
	return p
}

// 变量
func (this *DiskSource) Variables() []*SourceVariable {
	return []*SourceVariable{
		{
			Code:        "partitions",
			Description: "分区信息",
		},
		{
			Code:        "partitions.$.name",
			Description: "分区名",
		},
		{
			Code:        "partitions.$.total",
			Description: "总空间尺寸（字节）",
		},
		{
			Code:        "partitions.$.used",
			Description: "已使用空间尺寸（字节）",
		},
		{
			Code:        "partitions.$.free",
			Description: "剩余空间尺寸（字节）",
		},
		{
			Code:        "partitions.$.percent",
			Description: "已使用百分比",
		},
		{
			Code:        "partitions.$.inodesTotal",
			Description: "inodes总数",
		},
		{
			Code:        "partitions.$.inodesUsed",
			Description: "已使用inodes数量",
		},
		{
			Code:        "partitions.$.inodesFree",
			Description: "剩余inodes数量",
		},
		{
			Code:        "partitions.$.inodesPercent",
			Description: "inodes使用百分比",
		},
		{
			Code:        "partitions.$.fstype",
			Description: "文件系统类型",
		},
	}
}

// 阈值
func (this *DiskSource) Thresholds() []*Threshold {
	result := []*Threshold{}

	{
		t := NewThreshold()
		t.Param = "${partitions.$.percent}"
		t.Operator = ThresholdOperatorGt
		t.Value = "80"
		t.NoticeLevel = notices.NoticeLevelWarning
		t.NoticeMessage = "${ROW.name}分区使用已达到${ROW.percent|round(2)}%"
		result = append(result, t)
	}

	return result
}

// 图表
func (this *DiskSource) Charts() []*widgets.Chart {
	charts := []*widgets.Chart{}

	{
		chart := widgets.NewChart()
		chart.Id = "disk.usage.chart1"
		chart.Name = "文件系统"
		chart.Columns = 2
		chart.Type = "javascript"
		chart.Options = maps.Map{
			"code": `
var chart = new charts.StackBarChart();
chart.values = [];
chart.labels = [];

var latest = NewQuery().cache(120).latest(1);
if (latest.length > 0) {
	var partitions = latest[0].value.partitions;
	partitions.$each(function (k, v) {
		chart.values.push([v.used, v.total - v.used]);
		chart.labels.push(v.name + "（" + (Math.round(v.used / 1024 / 1024 / 1024 * 100) / 100)+ "G/" + (Math.round(v.total / 1024 / 1024 / 1024 * 100) / 100) +"G）");
	});

	chart.options.height = partitions.length * 4;
}

chart.colors = [ colors.BROWN, colors.GREEN ];
chart.render();
`,
		}

		charts = append(charts, chart)
	}

	return charts
}
