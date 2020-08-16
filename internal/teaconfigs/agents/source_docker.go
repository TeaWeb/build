package agents

import (
	"github.com/TeaWeb/build/internal/teaconfigs/forms"
	"github.com/TeaWeb/build/internal/teaconfigs/notices"
	"github.com/TeaWeb/build/internal/teaconfigs/widgets"
	"github.com/TeaWeb/build/internal/teautils"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/types"
	"strings"
)

// Docker相关数据源
type DockerSource struct {
	Source `yaml:",inline"`
}

// 获取新对象
func NewDockerSource() *DockerSource {
	return &DockerSource{}
}

// 名称
func (this *DockerSource) Name() string {
	return "Docker"
}

// 代号
func (this *DockerSource) Code() string {
	return "docker"
}

// 描述
func (this *DockerSource) Description() string {
	return "Docker状态和统计（使用docker ps和docker stats）"
}

// 执行
func (this *DockerSource) Execute(params map[string]string) (value interface{}, err error) {
	containerList := []maps.Map{}

	// status
	cmd := teautils.NewCommandExecutor()
	fields := []string{
		"ID", "Image", "Command", "CreatedAt", "RunningFor", "Ports", "Status", "Size", "Names", "Labels", "Mounts", "Networks",
	}
	formats := []string{}
	for _, field := range fields {
		formats = append(formats, "{{."+field+"}}")
	}

	cmd.Add("docker", "ps", "-a", "--no-trunc", "--format", strings.Join(formats, "||"))
	output, err := cmd.Run()
	if err != nil {
		return containerList, err
	}

	if len(output) == 0 {
		return containerList, nil
	}

	for _, line := range strings.Split(output, "\n") {
		values := strings.Split(line, "||")
		container := maps.Map{}
		for index, field := range fields {
			field = this.lowerCaseFirst(field)
			if index < len(values) {
				container[field] = values[index]
			} else {
				container[field] = ""
			}
		}
		containerList = append(containerList, container)
	}

	{
		cmd := teautils.NewCommandExecutor()
		fields := []string{
			"ID", "CPUPerc", "MemUsage", "NetIO", "BlockIO", "MemPerc", "PIDs",
		}
		formats := []string{}
		for _, field := range fields {
			formats = append(formats, "{{."+field+"}}")
		}

		cmd.Add("docker", "stats", "-a", "--no-stream", "--no-trunc", "--format", strings.Join(formats, "||"))
		output, err := cmd.Run()
		if err != nil {
			return containerList, err
		}
		for _, line := range strings.Split(output, "\n") {
			values := strings.Split(line, "||")
			var container = maps.Map{}
			for index, field := range fields {
				field = this.lowerCaseFirst(field)
				if index < len(values) {
					container[field] = values[index]

					if field == "cpuPercent" {
						container["cpuPercentString"] = values[index]
						container["cpuPercent"] = types.Float64(strings.Replace(values[index], "%", "", -1))
					} else if field == "memPercent" {
						container["memPercentString"] = values[index]
						container["memPercent"] = types.Float64(strings.Replace(values[index], "%", "", -1))
					} else if field == "pids" {
						container["pids"] = types.Int(values[index])
					}
				}
			}

			// 复制到上一个列表中
			for _, c := range containerList {
				if c.GetString("id") == container.GetString("id") {
					for k, v := range container {
						c[k] = v
					}
					break
				}
			}
		}
	}

	return containerList, nil
}

// 表单信息
func (this *DockerSource) Form() *forms.Form {
	form := forms.NewForm(this.Code())
	return form
}

// 变量
func (this *DockerSource) Variables() []*SourceVariable {
	return []*SourceVariable{
		{
			Code:        "$.blockIO",
			Description: "Block IO",
		},
		{
			Code:        "$.command",
			Description: "命令行",
		},
		{
			Code:        "$.cpuPercent",
			Description: "CPU使用比例",
		},
		{
			Code:        "$.cpuPercentString",
			Description: "CPU使用比例（字符串）",
		},
		{
			Code:        "$.createdAt",
			Description: "容器创建时间",
		},
		{
			Code:        "$.id",
			Description: "容器ID",
		},
		{
			Code:        "$.image",
			Description: "镜像Image",
		},
		{
			Code:        "$.labels",
			Description: "标签",
		},
		{
			Code:        "$.memPercent",
			Description: "内存使用比例",
		},
		{
			Code:        "$.memPercentString",
			Description: "内存使用比例（字符串）",
		},
		{
			Code:        "$.memUsage",
			Description: "内容使用",
		},
		{
			Code:        "$.mounts",
			Description: "挂载分卷",
		},
		{
			Code:        "$.names",
			Description: "容器名",
		},
		{
			Code:        "$.netIO",
			Description: "网络IO",
		},
		{
			Code:        "$.networks",
			Description: "绑定网络",
		}, {
			Code:        "$.pids",
			Description: "PId数量（在Windows上不可用）",
		},
		{
			Code:        "$.ports",
			Description: "开放端口",
		},
		{
			Code:        "$.runningFor",
			Description: "已运行时间",
		},
		{
			Code:        "$.size",
			Description: "容器硬盘空间占用",
		},
		{
			Code:        "$.status",
			Description: "容器状态",
		},
	}
}

// 阈值
func (this *DockerSource) Thresholds() []*Threshold {
	result := []*Threshold{}

	{
		t := NewThreshold()
		t.Param = "${$.status}"
		t.NoticeLevel = notices.NoticeLevelWarning
		t.Operator = ThresholdOperatorRegexp
		t.Value = "(?i)Exited"
		t.NoticeMessage = "${ROW.names}容器处于Exited状态"
		result = append(result, t)
	}

	return result
}

// 图表
func (this *DockerSource) Charts() []*widgets.Chart {
	charts := []*widgets.Chart{}

	{
		chart := widgets.NewChart()
		chart.Columns = 2
		chart.Name = "Docker状态"
		chart.Type = "javascript"
		chart.Options = maps.Map{
			"code": `var chart = new charts.TableChart();
var query = NewQuery();
query.desc();
query.limit(1);
var ones = query.findAll();
if (ones != null && ones.length > 0) {
	ones[0].value.$each(function (k, one) {
		if (one.status.match(/Exit/i)) {
			one.id = "<span class=\"red\">" + one.id + "</span>";
			one.names = "<span class=\"red\">" + one.names + "</span>";
			one.image = "<span class=\"red\">" + one.image + "</span>";
			one.status = "<span class=\"red\">" + one.status + "</span>";	
		} else if (one.status.match(/Up /i)) {
			one.id = "<span class=\"green\">" + one.id + "</span>";
			one.names = "<span class=\"green\">" + one.names + "</span>";
			one.image = "<span class=\"green\">" + one.image + "</span>";
			one.status = "<span class=\"green\">" + one.status + "</span>";	
		}
		chart.addRow(one.id, one.image, one.names, one.status);
	});	
}
chart.setWidth(0, "eight wide");
chart.render();`,
		}

		charts = append(charts, chart)
	}

	return charts
}

func (this *DockerSource) lowerCaseFirst(s string) string {
	if len(s) == 0 {
		return s
	}
	if s == "ID" {
		return "id"
	}
	if s == "CPUPerc" {
		return "cpuPercent"
	}
	if s == "MemPerc" {
		return "memPercent"
	}
	if s == "PIDs" {
		return "pids"
	}
	return strings.ToLower(s[0:1]) + s[1:]
}
