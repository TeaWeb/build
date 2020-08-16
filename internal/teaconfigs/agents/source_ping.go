package agents

import (
	"bytes"
	"errors"
	"github.com/TeaWeb/build/internal/teaconfigs/forms"
	"github.com/TeaWeb/build/internal/teaconfigs/notices"
	"github.com/TeaWeb/build/internal/teaconfigs/widgets"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/types"
	"github.com/tatsushid/go-fastping"
	"net"
	"os/exec"
	"regexp"
	"runtime"
	"time"
)

// Ping
type PingSource struct {
	Source `yaml:",inline"`

	Host string `yaml:"host" json:"host"`
}

// 获取新对象
func NewPingSource() *PingSource {
	return &PingSource{}
}

// 名称
func (this *PingSource) Name() string {
	return "Ping"
}

// 代号
func (this *PingSource) Code() string {
	return "ping"
}

// 描述
func (this *PingSource) Description() string {
	return "通过Ping获取主机响应时间"
}

// 执行
func (this *PingSource) Execute(params map[string]string) (value interface{}, err error) {
	host := this.Host

	// 去除http|https|ftp
	host = regexp.MustCompile(`^(?i)(http|https|ftp)://`).ReplaceAllLiteralString(host, "")

	if len(host) == 0 {
		err = errors.New("'host' should not be empty")
		return maps.Map{
			"rtt": -1,
		}, err
	}

	if runtime.GOOS == "linux" { // Linux
		value, err = this.pingLinux(host)
		if err == nil {
			return
		}
	} else if runtime.GOOS == "freebsd" {
		value, err = this.pingFreebsd(host)
		if err == nil {
			return
		}
	} else if runtime.GOOS == "windows" { // windows
		value, err = this.pingWindows(host)
		if err == nil {
			return
		}
	}

	p := fastping.NewPinger()
	if runtime.GOOS == "darwin" {
		_, err = p.Network("udp")
	} else {
		_, err = p.Network("ip")
	}
	if err != nil {
		return maps.Map{
			"rtt": -1,
		}, err
	}

	ra, err := net.ResolveIPAddr("ip4:icmp", host)
	if err != nil {
		return maps.Map{
			"rtt": -1,
		}, err
	}
	p.AddIPAddr(ra)

	p.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {
		value = maps.Map{
			"rtt": rtt.Seconds() * 1000,
		}
	}
	p.OnIdle = func() {
		if value == nil {
			err = errors.New("ping timeout")
		}
	}

	runningErr := p.Run()
	if runningErr != nil {
		return maps.Map{
			"rtt": -1,
		}, runningErr
	}

	if err != nil {
		return maps.Map{
			"rtt": -1,
		}, err
	}

	return
}

// 表单信息
func (this *PingSource) Form() *forms.Form {
	form := forms.NewForm(this.Code())
	{
		group := form.NewGroup()
		{
			field := forms.NewTextField("主机地址", "Host")
			field.IsRequired = true
			field.Code = "host"
			field.ValidateCode = `
if (value.length == 0) {
	throw new Error("请输入主机地址");
}
`
			field.Comment = "要Ping的主机地址，可以是一个域名或一个IP"
			group.Add(field)
		}
	}
	return form
}

// 变量
func (this *PingSource) Variables() []*SourceVariable {
	return []*SourceVariable{
		{
			Code:        "rtt",
			Description: "响应时间（单位ms）",
		},
	}
}

// 阈值
func (this *PingSource) Thresholds() []*Threshold {
	result := []*Threshold{}

	{
		t := NewThreshold()
		t.Param = "${rtt}"
		t.Operator = ThresholdOperatorEq
		t.Value = "-1"
		t.NoticeLevel = notices.NoticeLevelWarning
		t.NoticeMessage = "Ping超时"
		t.MaxFails = 5
		result = append(result, t)
	}

	return result
}

// 图表
func (this *PingSource) Charts() []*widgets.Chart {
	charts := []*widgets.Chart{}

	{
		// chart
		chart := widgets.NewChart()
		chart.Name = "Ping"
		chart.Columns = 2
		chart.Type = "javascript"
		chart.SupportsTimeRange = true
		chart.Options = maps.Map{
			"code": `var chart = new charts.LineChart();

var ones = NewQuery().past(60, time.MINUTE).avg("rtt");

var line = new charts.Line();
line.isFilled = true;

ones.$each(function (k, v) {
	line.addValue(v.value.rtt);
	chart.addLabel(v.label);
});

chart.addLine(line);
chart.render();`,
		}

		charts = append(charts, chart)
	}

	return charts
}

// 显示信息
func (this *PingSource) Presentation() *forms.Presentation {
	p := forms.NewPresentation()
	p.HTML = `
<tr>
	<td>主机地址</td>
	<td>{{source.host}}</td>
</tr>
`
	return p
}

func (this *PingSource) pingLinux(host string) (value interface{}, err error) {
	value = maps.Map{
		"rtt": -1,
	}

	pingExe, err := exec.LookPath("ping")
	if err != nil {
		return
	}
	stdout := bytes.NewBuffer([]byte{})
	cmd := exec.Command(pingExe, "-c", "3", "-W", "3", host)
	cmd.Stdout = stdout
	err = cmd.Start()
	if err != nil {
		return
	}
	err = cmd.Wait()
	if err != nil {
		return
	}

	// 匹配 time=x ms
	results := regexp.MustCompile(`time=([0-9.]+)\s+ms`).FindAllStringSubmatch(string(stdout.Bytes()), -1)
	if len(results) == 0 {
		value = -1
		err = errors.New("timeout")
		return
	}

	total := float32(0)
	for _, result := range results {
		total += types.Float32(result[1])
	}
	value = maps.Map{
		"rtt": total / float32(len(results)),
	}

	return
}

func (this *PingSource) pingFreebsd(host string) (value interface{}, err error) {
	value = maps.Map{
		"rtt": -1,
	}

	pingExe, err := exec.LookPath("ping")
	if err != nil {
		return
	}
	stdout := bytes.NewBuffer([]byte{})
	cmd := exec.Command(pingExe, "-c", "3", "-W", "3000", host) // -W 单位是ms
	cmd.Stdout = stdout
	err = cmd.Start()
	if err != nil {
		return
	}
	err = cmd.Wait()
	if err != nil {
		return
	}

	// 匹配 time=x ms
	results := regexp.MustCompile(`time=([0-9.]+)\s+ms`).FindAllStringSubmatch(string(stdout.Bytes()), -1)
	if len(results) == 0 {
		value = -1
		err = errors.New("timeout")
		return
	}

	total := float32(0)
	for _, result := range results {
		total += types.Float32(result[1])
	}
	value = maps.Map{
		"rtt": total / float32(len(results)),
	}

	return
}

func (this *PingSource) pingWindows(host string) (value interface{}, err error) {
	logs.Println("ping:", host)
	value = maps.Map{
		"rtt": -1,
	}

	pingExe, err := exec.LookPath("ping")
	if err != nil {
		return
	}
	stdout := bytes.NewBuffer([]byte{})
	cmd := exec.Command(pingExe, "-n", "3", host)
	cmd.Stdout = stdout
	err = cmd.Start()
	if err != nil {
		return
	}
	err = cmd.Wait()
	if err != nil {
		return
	}

	// 匹配 time=[x]ms
	results := regexp.MustCompile(`[=<]([0-9.]+)ms\s+TTL=`).FindAllStringSubmatch(string(stdout.Bytes()), -1)
	if len(results) == 0 {
		value = -1
		err = errors.New("timeout")
		return
	}

	total := float32(0)
	for _, result := range results {
		total += types.Float32(result[1])
	}
	value = maps.Map{
		"rtt": total / float32(len(results)),
	}

	return
}
