package agents

import (
	"errors"
	"github.com/TeaWeb/build/internal/teaconfigs/forms"
	"github.com/TeaWeb/build/internal/teaconfigs/notices"
	"github.com/TeaWeb/build/internal/teaconfigs/widgets"
	"github.com/iwind/TeaGo/files"
	"github.com/iwind/TeaGo/maps"
)

// 监控文件变化
// TODO 增加权限变化监控
type FileChangeSource struct {
	Source `yaml:",inline"`
	Path   string `yaml:"path" json:"path"`

	exists        bool
	isInitialized bool
	modifiedAt    int64
	size          int64
}

// 获取新对象
func NewFileChangeSource() *FileChangeSource {
	return &FileChangeSource{}
}

// 名称
func (this *FileChangeSource) Name() string {
	return "文件变化"
}

// 代号
func (this *FileChangeSource) Code() string {
	return "file.change"
}

// 描述
func (this *FileChangeSource) Description() string {
	return "监控某个文件或目录是否有变化，主要通过文件尺寸、修改时间、是否存在等项目来判断"
}

// 执行
func (this *FileChangeSource) Execute(params map[string]string) (value interface{}, err error) {
	if len(this.Path) == 0 {
		err = errors.New("'path' should not be empty")
		return
	}

	file := files.NewFile(this.Path)
	isInitialized := this.isInitialized
	if !this.isInitialized {
		this.isInitialized = true
	}

	exists := file.Exists()
	modifiedAt := int64(0)
	size := int64(0)
	if exists {
		modifiedTime, err := file.LastModified()
		if err != nil {
			return nil, err
		}

		modifiedAt = modifiedTime.Unix()

		newSize, err := file.Size()
		if err != nil {
			return nil, err
		}
		size = newSize
	}

	if !isInitialized {
		this.exists = exists
		this.size = size
		this.modifiedAt = modifiedAt

		return maps.Map{
			"exists":     exists,
			"changes":    false,
			"size":       size,
			"modifiedAt": modifiedAt,
		}, nil
	}

	value = maps.Map{
		"exists":     exists,
		"changes":    this.exists != exists || this.modifiedAt != modifiedAt || this.size != size,
		"size":       size,
		"modifiedAt": modifiedAt,
	}

	this.exists = exists
	this.size = size
	this.modifiedAt = modifiedAt

	return
}

// 表单信息
func (this *FileChangeSource) Form() *forms.Form {
	form := forms.NewForm(this.Code())
	{
		group := form.NewGroup()
		{
			field := forms.NewTextField("文件路径", "Path")
			field.IsRequired = true
			field.Code = "path"
			field.Comment = "请输入要监控的文件或目录路径"
			field.ValidateCode = `
if (value.length == 0) {
	throw new Error("请输入要监控的文件或目录路径");
}
`
			group.Add(field)
		}
	}
	return form
}

// 变量
func (this *FileChangeSource) Variables() []*SourceVariable {
	return []*SourceVariable{
		{
			Code:        "changes",
			Description: "是否有变化",
		},
		{
			Code:        "exists",
			Description: "是否存在",
		},
		{
			Code:        "modifiedAt",
			Description: "修改时间戳",
		},
		{
			Code:        "size",
			Description: "文件尺寸",
		},
	}
}

// 阈值
func (this *FileChangeSource) Thresholds() []*Threshold {
	result := []*Threshold{}

	// 阈值
	{
		t := NewThreshold()
		t.Param = "${changes}"
		t.Operator = ThresholdOperatorEq
		t.Value = "true"
		t.NoticeLevel = notices.NoticeLevelWarning
		t.NoticeMessage = "监控的文件有变化"
		result = append(result, t)
	}

	return result
}

// 图表
func (this *FileChangeSource) Charts() []*widgets.Chart {
	charts := []*widgets.Chart{}

	{
		// chart
		chart := widgets.NewChart()
		chart.Name = "文件尺寸<em>（字节）</em>"
		chart.Columns = 2
		chart.Type = "javascript"
		chart.SupportsTimeRange = true
		chart.Options = maps.Map{
			"code": `var chart = new charts.LineChart();

var ones = NewQuery().past(60, time.MINUTE).avg("size");

var line = new charts.Line();
line.isFilled = true;

ones.$each(function (k, v) {
	line.addValue(v.value.size);
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
func (this *FileChangeSource) Presentation() *forms.Presentation {
	p := forms.NewPresentation()
	p.HTML = `
<tr>
	<td>文件路径</td>
	<td>{{source.path}}</td>
</tr>
`
	return p
}
