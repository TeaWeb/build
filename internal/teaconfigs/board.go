package teaconfigs

// 看板
type Board struct {
	Charts []*BoardChart `yaml:"charts" json:"charts"`
}

// 获取新对象
func NewBoard() *Board {
	return &Board{}
}

// 判断是否在使用某个Chart
func (this *Board) HasChart(chartId string) bool {
	for _, c := range this.Charts {
		if c.ChartId == chartId {
			return true
		}
	}

	return false
}

// 移除某个Chart
func (this *Board) RemoveChart(chartId string) {
	result := []*BoardChart{}
	for _, c := range this.Charts {
		if c.ChartId == chartId {
			continue
		}
		result = append(result, c)
	}
	this.Charts = result
}

// 添加Chart
func (this *Board) AddChart(widgetId string, chartId string) {
	if this.HasChart(chartId) {
		return
	}
	this.Charts = append(this.Charts, &BoardChart{
		WidgetId: widgetId,
		ChartId:  chartId,
	})
}

// 移动Chart
func (this *Board) MoveChart(fromIndex int, toIndex int) {
	if fromIndex < 0 || fromIndex >= len(this.Charts) {
		return
	}
	if toIndex < 0 || toIndex >= len(this.Charts) {
		return
	}
	if fromIndex == toIndex {
		return
	}

	chart := this.Charts[fromIndex]
	newList := []*BoardChart{}
	for i := 0; i < len(this.Charts); i ++ {
		if i == fromIndex {
			continue
		}
		if fromIndex > toIndex && i == toIndex {
			newList = append(newList, chart)
		}
		newList = append(newList, this.Charts[i])
		if fromIndex < toIndex && i == toIndex {
			newList = append(newList, chart)
		}
	}

	this.Charts = newList
}
