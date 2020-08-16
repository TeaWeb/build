package board

import (
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/iwind/TeaGo/actions"
)

type MoveChartAction actions.Action

// 移动图表位置
func (this *MoveChartAction) Run(params struct {
	ServerId string
	Type     string
	OldIndex int
	NewIndex int
}) {
	server := teaconfigs.NewServerConfigFromId(params.ServerId)
	if server == nil {
		this.Fail("找不到Server")
	}

	switch params.Type {
	case "realtime":
		server.RealtimeBoard.MoveChart(params.OldIndex, params.NewIndex)
	case "stat":
		server.StatBoard.MoveChart(params.OldIndex, params.NewIndex)
	}

	err := server.Save()
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	this.Success()
}
