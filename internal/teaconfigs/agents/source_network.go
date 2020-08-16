package agents

import (
	"github.com/TeaWeb/build/internal/teaconfigs/forms"
	"github.com/TeaWeb/build/internal/teaconfigs/notices"
	"github.com/TeaWeb/build/internal/teaconfigs/widgets"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"github.com/shirou/gopsutil/net"
	"time"
)

// 网络带宽等信息
type NetworkSource struct {
	Source `yaml:",inline"`

	InterfaceNames []string `yaml:"interfaceNames" json:"interfaceNames"`

	lastTime    time.Time
	lastStats   []net.IOCountersStat
	lastSumStat net.IOCountersStat
}

// 获取新对象
func NewNetworkSource() *NetworkSource {
	return &NetworkSource{}
}

// 名称
func (this *NetworkSource) Name() string {
	return "网络信息"
}

// 代号
func (this *NetworkSource) Code() string {
	return "network"
}

// 描述
func (this *NetworkSource) Description() string {
	return "网络接口、带宽信息"
}

// 执行
func (this *NetworkSource) Execute(params map[string]string) (value interface{}, err error) {
	interfaces := []map[string]interface{}{}
	interfaceStats, err := net.Interfaces()
	interfaceNames := []string{}
	if err != nil {
		logs.Error(err)
	} else {
		for _, i := range interfaceStats {
			if len(i.HardwareAddr) == 0 {
				continue
			}
			if len(this.InterfaceNames) > 0 && !lists.ContainsString(this.InterfaceNames, i.Name) {
				continue
			}
			interfaceNames = append(interfaceNames, i.Name)
			interfaces = append(interfaces, map[string]interface{}{
				"name":         i.Name,
				"hardwareAddr": i.HardwareAddr,
				"mtu":          i.MTU,
				"flags":        i.Flags,
				"addrs": lists.Map(i.Addrs, func(k int, v interface{}) interface{} {
					addr := v.(net.InterfaceAddr)
					return addr.Addr
				}),
			})
		}
	}

	value = map[string]interface{}{
		"interfaces": interfaces,
		"stat": map[string]interface{}{
			"avgSentBytes":         0,
			"avgSentPackets":       0,
			"avgReceivedBytes":     0,
			"avgReceivedPackets":   0,
			"totalSentBytes":       0,
			"totalSentPackets":     0,
			"totalReceivedBytes":   0,
			"totalReceivedPackets": 0,
		},
	}

	for _, interfaceName := range interfaceNames {
		value.(map[string]interface{})[interfaceName] = map[string]interface{}{
			"avgSentBytes":         0,
			"avgSentPackets":       0,
			"avgReceivedBytes":     0,
			"avgReceivedPackets":   0,
			"totalSentBytes":       0,
			"totalSentPackets":     0,
			"totalReceivedBytes":   0,
			"totalReceivedPackets": 0,
		}
	}

	stats, err := net.IOCounters(true)
	if err != nil {
		return value, err
	}

	now := time.Now()
	seconds := uint64(now.Unix() - this.lastTime.Unix())
	if seconds == 0 {
		seconds = 1 // 避免被除数为0
	}

	sumStat := net.IOCountersStat{}
	sumStat.Name = "ALL"
	for _, interfaceName := range interfaceNames {
		for _, stat := range stats {
			if stat.Name == interfaceName {
				m := map[string]interface{}{
					"avgSentBytes":         0,
					"avgSentPackets":       0,
					"avgReceivedBytes":     0,
					"avgReceivedPackets":   0,
					"totalSentBytes":       stat.BytesSent,
					"totalSentPackets":     stat.PacketsSent,
					"totalReceivedBytes":   stat.BytesRecv,
					"totalReceivedPackets": stat.PacketsRecv,
				}

				sumStat.BytesSent += stat.BytesSent
				sumStat.PacketsSent += stat.PacketsSent
				sumStat.BytesRecv += stat.BytesRecv
				sumStat.PacketsRecv += stat.PacketsRecv

				// old stat
				var foundStat net.IOCountersStat
				for _, oldStat := range this.lastStats {
					if oldStat.Name == interfaceName {
						foundStat = oldStat
						break
					}
				}

				// not found
				if len(foundStat.Name) == 0 {
					foundStat = stat
				}
				m["avgSentBytes"] = (stat.BytesSent - foundStat.BytesSent) / seconds
				m["avgSentPackets"] = (stat.PacketsSent - foundStat.PacketsSent) / seconds
				m["avgReceivedBytes"] = (stat.BytesRecv - foundStat.BytesRecv) / seconds
				m["avgReceivedPackets"] = (stat.PacketsRecv - foundStat.PacketsRecv) / seconds
				value.(map[string]interface{})[interfaceName] = m
				break
			}
		}
	}

	if len(this.lastSumStat.Name) == 0 {
		this.lastSumStat = sumStat
	}
	value.(map[string]interface{})["stat"] = map[string]interface{}{
		"avgSentBytes":         (sumStat.BytesSent - this.lastSumStat.BytesSent) / seconds,
		"avgSentPackets":       (sumStat.PacketsSent - this.lastSumStat.PacketsSent) / seconds,
		"avgReceivedBytes":     (sumStat.BytesRecv - this.lastSumStat.BytesRecv) / seconds,
		"avgReceivedPackets":   (sumStat.PacketsRecv - this.lastSumStat.PacketsRecv) / seconds,
		"totalSentBytes":       sumStat.BytesSent,
		"totalSentPackets":     sumStat.PacketsSent,
		"totalReceivedBytes":   sumStat.BytesRecv,
		"totalReceivedPackets": sumStat.PacketsRecv,
	}

	this.lastTime = now
	this.lastStats = stats
	this.lastSumStat = sumStat

	return
}

// 表单信息
func (this *NetworkSource) Form() *forms.Form {
	form := forms.NewForm(this.Code())

	{
		group := form.NewGroup()
		{
			field := forms.NewSingleValueList("网卡名称", "Interface")
			field.ValueName = "网卡名称"
			field.Comment = "可以添加只需要监控的某些网络接口名称，比如en0, em1, eth2等。默认监控全部网络接口。"
			field.Code = "interfaceNames"
			if len(this.InterfaceNames) == 0 {
				this.InterfaceNames = []string{}
			}
			field.Value = this.InterfaceNames
			group.Add(field)
		}
	}

	return form
}

// 变量
func (this *NetworkSource) Variables() []*SourceVariable {
	return []*SourceVariable{
		{
			Code:        "interfaces",
			Description: "网络接口",
		},
		{
			Code:        "interfaces.$.name",
			Description: "接口名称",
		},
		{
			Code:        "interfaces.$.addrs",
			Description: "接口地址",
		},
		{
			Code:        "interfaces.$.flags",
			Description: "接口标识",
		},
		{
			Code:        "interfaces.$.hardwareAddr",
			Description: "接口硬件地址",
		},
		{
			Code:        "interfaces.$.mtu",
			Description: "接口MTU值",
		},
		{
			Code:        "stat",
			Description: "总流量统计信息",
		},
		{
			Code:        "stat.avgReceivedBytes",
			Description: "总平均接收速率（秒）",
		},
		{
			Code:        "stat.avgReceivedPackets",
			Description: "总平均接收的数据包数量速率（秒）",
		},
		{
			Code:        "stat.avgSentBytes",
			Description: "总平均发送速率（秒）",
		},
		{
			Code:        "stat.avgSentPackets",
			Description: "总平均发送的数据包数量速率（秒）",
		},
		{
			Code:        "stat.totalReceivedBytes",
			Description: "总接收字节数",
		},
		{
			Code:        "stat.totalReceivedPackets",
			Description: "总接收数据包数量",
		},
		{
			Code:        "stat.totalSentBytes",
			Description: "总发送字节数",
		},
		{
			Code:        "stat.totalSentPackets",
			Description: "总发送数据包数量",
		},
		{
			Code:        "$",
			Description: "单个接口流量信息",
		},
		{
			Code:        "$.avgReceivedBytes",
			Description: "单个接口总平均接收速率（秒）",
		},
		{
			Code:        "$.avgReceivedPackets",
			Description: "单个接口总平均接收数据包数量速率（秒）",
		},
		{
			Code:        "$.avgSentBytes",
			Description: "单个接口总平均发送速率（秒）",
		},
		{
			Code:        "$.avgSentPackets",
			Description: "单个接口总平均发送数据包数量速率（秒）",
		},
		{
			Code:        "$.totalReceivedBytes",
			Description: "单个接口总接收字节数",
		},
		{
			Code:        "$.totalReceivedPackets",
			Description: "单个接口总接收数据包数量",
		},
		{
			Code:        "$.totalSentBytes",
			Description: "单个接口总发送字节数",
		},
		{
			Code:        "$.totalSentPackets",
			Description: "单个接口总发送数据包数量",
		},
	}
}

// 阈值
func (this *NetworkSource) Thresholds() []*Threshold {
	result := []*Threshold{}

	{
		t := NewThreshold()
		t.Param = "${stat.avgSentBytes}"
		t.Operator = ThresholdOperatorGte
		t.Value = "13107200"
		t.NoticeLevel = notices.NoticeLevelWarning
		t.NoticeMessage = "当前出口流量超过100MBit/s"
		result = append(result, t)
	}

	return result
}

// 图表
func (this *NetworkSource) Charts() []*widgets.Chart {
	charts := []*widgets.Chart{}

	// 图表
	{
		chart := widgets.NewChart()
		chart.Id = "network.usage.sent"
		chart.Name = "出口带宽（M/s）"
		chart.Columns = 2
		chart.Type = "javascript"
		chart.SupportsTimeRange = true
		chart.Options = maps.Map{
			"code": `var chart = new charts.LineChart();

var line = new charts.Line();
line.isFilled = true;

var ones = NewQuery().past(60, time.MINUTE).avg("stat.avgSentBytes");
ones.$each(function (k, v) {
	line.addValue(Math.round(v.value.stat.avgSentBytes / 1024 / 1024 * 100) / 100);
	chart.addLabel(v.label);
});
var maxValue = line.values.$max();
if (maxValue < 1) {
	chart.max = 1;
} else if (maxValue < 5) {
	chart.max = 5;
} else if (maxValue < 10) {
	chart.max = 10;
}

chart.addLine(line);
chart.render();`,
		}
		charts = append(charts, chart)
	}

	{
		chart := widgets.NewChart()
		chart.Id = "network.usage.received"
		chart.Name = "入口带宽（M/s）"
		chart.Columns = 2
		chart.Type = "javascript"
		chart.SupportsTimeRange = true
		chart.Options = maps.Map{
			"code": `var chart = new charts.LineChart();

var line = new charts.Line();
line.isFilled = true;

var ones = NewQuery().past(60, time.MINUTE).avg("stat.avgReceivedBytes");
ones.$each(function (k, v) {
	line.addValue(Math.round(v.value.stat.avgReceivedBytes / 1024 / 1024 * 100) / 100);
	chart.addLabel(v.label);
});
var maxValue = line.values.$max();
if (maxValue < 1) {
	chart.max = 1;
} else if (maxValue < 5) {
	chart.max = 5;
} else if (maxValue < 10) {
	chart.max = 10;
}

chart.addLine(line);
chart.render();`,
		}
		charts = append(charts, chart)
	}

	{
		chart := widgets.NewChart()
		chart.Id = "network.usage.sent.interfaces"
		chart.Name = "网卡出口流量（M/s）"
		chart.Columns = 2
		chart.Type = "javascript"
		chart.SupportsTimeRange = true
		chart.Options = maps.Map{
			"code": `var chart = new charts.LineChart();

var interfaceNames = []; // interface names
if (interfaceNames.length == 0) {
	var one = NewQuery().desc("_id").find();
	if (one != null) {
		interfaceNames = one.value["interfaces"].$map(function (k, v) {
			return v.name;
		});
	}
	interfaceNames.sort();
}

var fields = [];
var lines = [];
interfaceNames.$each(function (k, v) {
	fields.push(v + ".avgSentBytes");
	
	var line = new charts.Line();
	line.isFilled = true;
	line.color =colors.ARRAY[k % colors.ARRAY.length];
	lines.push(line);
	line.name = v;
	chart.addLine(line);
});
var query = NewQuery().past(60, time.MINUTE);
var ones = query.avg.apply(query, fields);

var max = 0;
ones.$each(function (k, v) {
	for (var i = 0; i < lines.length; i ++) {
		var line = lines[i];
		var sent = v.value[interfaceNames[i]].avgSentBytes;
		if (sent == null) {
			sent = 0;
		}
		if (sent > max) {
			max = sent;	
		}
		line.addValue(Math.round(sent / 1024 / 1024 * 100) / 100);
	}
	
	chart.addLabel(v.label);
});
	
max = max / 1024 / 1024;
if (max < 1) {
	chart.max = 1;	
} else if (max < 5) {
	chart.max = 5;
} else if (max < 10) {
	chart.max = 10;
} else {
	chart.max = Math.ceil(max / 10)	* 10;
}
chart.render();`,
		}
		charts = append(charts, chart)
	}

	{
		chart := widgets.NewChart()
		chart.Id = "network.usage.received.interfaces"
		chart.Name = "网卡入口流量（M/s）"
		chart.Columns = 2
		chart.Type = "javascript"
		chart.SupportsTimeRange = true
		chart.Options = maps.Map{
			"code": `var chart = new charts.LineChart();

var interfaceNames = []; // interface names
if (interfaceNames.length == 0) {
	var one = NewQuery().desc("_id").find();
	if (one != null) {
		interfaceNames = one.value["interfaces"].$map(function (k, v) {
			return v.name;
		});
	}
	interfaceNames.sort();
}

var fields = [];
var lines = [];
interfaceNames.$each(function (k, v) {
	fields.push(v + ".avgReceivedBytes");
	
	var line = new charts.Line();
	line.isFilled = true;
	line.color =colors.ARRAY[k % colors.ARRAY.length];
	lines.push(line);
	line.name = v;
	chart.addLine(line);
});
var query = NewQuery().past(60, time.MINUTE);
var ones = query.avg.apply(query, fields);

var max = 0;
ones.$each(function (k, v) {
	for (var i = 0; i < lines.length; i ++) {
		var line = lines[i];
		var received = v.value[interfaceNames[i]].avgReceivedBytes;
		if (received == null) {
			received = 0;
		}
		if (received > max) {
			max = received;	
		}
		line.addValue(Math.round(received / 1024 / 1024 * 100) / 100);
	}
	
	chart.addLabel(v.label);
});
	
max = max / 1024 / 1024;
if (max < 1) {
	chart.max = 1;	
} else if (max < 5) {
	chart.max = 5;
} else if (max < 10) {
	chart.max = 10;
} else {
	chart.max = Math.ceil(max / 10)	* 10;
}
chart.render();`,
		}
		charts = append(charts, chart)
	}

	return charts
}

// 显示信息
func (this *NetworkSource) Presentation() *forms.Presentation {
	p := forms.NewPresentation()
	p.HTML = `
<tr>
	<td>网卡名称<em>（Interface）</em></td>
	<td>
		<span v-if="source.interfaceNames == null || source.interfaceNames.length == 0">全部网卡</span>
		<div v-if="source.interfaceNames != null && source.interfaceNames.length > 0">
			<span class="ui label tiny" v-for="interfaceName in source.interfaceNames">{{interfaceName}}</span>
		</div>
	</td>
</tr>`

	return p
}
