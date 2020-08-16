package agents

import (
	"errors"
	"github.com/TeaWeb/build/internal/teaconfigs/shared"
	"github.com/TeaWeb/build/internal/teaconst"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/files"
	"github.com/iwind/TeaGo/logs"
)

// 看板
type Board struct {
	TeaVersion string        `yaml:"teaVersion" json:"teaVersion"`
	Filename   string        `yaml:"filename" json:"filename"`
	Charts     []*BoardChart `yaml:"charts" json:"charts"`
}

// 取得Agent看板
func NewAgentBoard(agentId string) *Board {
	filename := "board." + agentId + ".conf"
	file := files.NewFile(Tea.ConfigFile("agents/" + filename))
	if file.Exists() {
		reader, err := file.Reader()
		if err != nil {
			logs.Error(err)
			return nil
		}
		defer func() {
			_ = reader.Close()
		}()
		board := &Board{}
		err = reader.ReadYAML(board)
		if err != nil {
			return nil
		}
		board.Filename = filename
		return board
	} else {
		return &Board{
			Filename: filename,
			Charts:   []*BoardChart{},
		}
	}
}

// 添加图表
func (this *Board) AddChart(appId, itemId, chartId string) {
	if this.FindChart(chartId) != nil {
		return
	}
	this.Charts = append(this.Charts, &BoardChart{
		AppId:   appId,
		ItemId:  itemId,
		ChartId: chartId,
	})
}

// 查找图表
func (this *Board) FindChart(chartId string) *BoardChart {
	for _, c := range this.Charts {
		if c.ChartId == chartId {
			return c
		}
	}
	return nil
}

// 查看是否有图表
func (this *Board) HasChart(chartId string) bool {
	for _, c := range this.Charts {
		if c.ChartId == chartId {
			return true
		}
	}
	return false
}

// 删除图表
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

// 删除App相关的所有图表
func (this *Board) RemoveApp(appId string) {
	result := []*BoardChart{}
	for _, c := range this.Charts {
		if c.AppId == appId {
			continue
		}
		result = append(result, c)
	}
	this.Charts = result
}

// 保存
func (this *Board) Save() error {
	shared.Locker.Lock()
	defer shared.Locker.WriteUnlockNotify()

	this.TeaVersion = teaconst.TeaVersion

	if len(this.Filename) == 0 {
		return errors.New("filename should be specified")
	}
	writer, err := files.NewWriter(Tea.ConfigFile("agents/" + this.Filename))
	if err != nil {
		return err
	}
	defer func() {
		_ = writer.Close()
	}()
	_, err = writer.WriteYAML(this)
	return err
}

// 移动图表
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
	for i := 0; i < len(this.Charts); i++ {
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
