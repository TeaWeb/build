package teaconfigs

import "github.com/TeaWeb/build/internal/teaconfigs/widgets"

// 看板图表
type BoardChart struct {
	WidgetId string `yaml:"widgetId" json:"widgetId"`
	ChartId  string `yaml:"chartId" json:"chartId"`
}

// 查找Chart实例
func (this *BoardChart) FindChart() (widget *widgets.Widget, chart *widgets.Chart) {
	if len(this.WidgetId) == 0 || len(this.ChartId) == 0 {
		return nil, nil
	}

	widget = widgets.NewWidgetFromId(this.WidgetId)
	if widget == nil {
		return nil, nil
	}

	return widget, widget.FindChart(this.ChartId)
}
