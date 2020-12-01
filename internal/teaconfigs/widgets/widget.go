package widgets

import (
	"github.com/TeaWeb/build/internal/teaconfigs/shared"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/files"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/rands"
	"gopkg.in/yaml.v3"
	"time"
)

// Widget定义
type Widget struct {
	Id          string   `yaml:"id" json:"id"`
	On          bool     `yaml:"on" json:"on"`
	Name        string   `yaml:"name" json:"name"`
	Code        string   `yaml:"code" json:"code"`
	Author      string   `yaml:"author" json:"author"`
	Version     string   `yaml:"version" json:"version"`
	Description string   `yaml:"description" json:"description"`
	Charts      []*Chart `yaml:"charts" json:"charts"`
	CreatedAt   int64    `yaml:"createdAt" json:"createdAt"` // 添加时间，用来排序
}

// 获取新对象
func NewWidget() *Widget {
	return &Widget{
		Id:        rands.HexString(16),
		On:        true,
		CreatedAt: time.Now().Unix(),
	}
}

// 从文件中加载Widget
func NewWidgetFromId(widgetId string) *Widget {
	file := files.NewFile(Tea.ConfigFile("widgets/widget." + widgetId + ".conf"))
	if !file.Exists() {
		return nil
	}
	data, err := file.ReadAll()
	if err != nil {
		logs.Error(err)
		return nil
	}
	widget := &Widget{}
	err = yaml.Unmarshal(data, widget)
	if err != nil {
		logs.Error(err)
		return nil
	}
	return widget
}

// 获取所有Widget列表
func LoadAllWidgets() []*Widget {
	dir := files.NewFile(Tea.ConfigFile("widgets"))
	if !dir.Exists() {
		return []*Widget{}
	}
	widgets := []*Widget{}
	for _, widgetFile := range dir.List() {
		data, err := widgetFile.ReadAll()
		if err != nil {
			logs.Error(err)
			continue
		}
		widget := &Widget{}
		err = yaml.Unmarshal(data, widget)
		if err != nil {
			logs.Error(err)
			continue
		}
		widgets = append(widgets, widget)
	}

	// 排序
	if len(widgets) > 0 {
		lists.Sort(widgets, func(i int, j int) bool {
			widget1 := widgets[i]
			widget2 := widgets[j]
			return widget1.CreatedAt > widget2.CreatedAt
		})
	}

	return widgets
}

// 保存
func (this *Widget) Save() error {
	shared.Locker.Lock()
	defer shared.Locker.WriteUnlockNotify()

	dir := files.NewFile(Tea.ConfigDir() + Tea.DS + "widgets")
	if !dir.Exists() {
		err := dir.Mkdir()
		if err != nil {
			return err
		}
	}
	writer, err := files.NewWriter(Tea.ConfigDir() + Tea.DS + "widgets" + Tea.DS + "widget." + this.Id + ".conf")
	if err != nil {
		return err
	}
	defer writer.Close()
	_, err = writer.WriteYAML(this)
	return err
}

// 删除当前Widget
func (this *Widget) Delete() error {
	file := files.NewFile(Tea.ConfigDir() + Tea.DS + "widgets" + Tea.DS + "widget." + this.Id + ".conf")
	if file.Exists() {
		return file.Delete()
	}
	return nil
}

// 添加Chart
func (this *Widget) AddChart(chart *Chart) {
	this.Charts = append(this.Charts, chart)
}

// 查找Chart
func (this *Widget) FindChart(chartId string) *Chart {
	for _, chart := range this.Charts {
		if chart.Id == chartId {
			return chart
		}
	}
	return nil
}

// 删除Chart
func (this *Widget) RemoveChart(chartId string) {
	result := []*Chart{}
	for _, chart := range this.Charts {
		if chart.Id == chartId {
			continue
		}
		result = append(result, chart)
	}
	this.Charts = result
}
