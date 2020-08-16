package apps

import (
	"github.com/TeaWeb/build/internal/teaweb/configs"
	"github.com/TeaWeb/build/internal/teaweb/helpers"
	"github.com/iwind/TeaGo"
)

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		server.
			Helper(&helpers.UserMustAuth{
				Grant: configs.AdminGrantAgent,
			}).
			Helper(new(Helper)).
			Prefix("/agents/board").
			GetPost("", new(IndexAction)).
			Get("/charts", new(ChartsAction)).
			Post("/addChart", new(AddChartAction)).
			Post("/removeChart", new(RemoveChartAction)).
			Post("/moveChart", new(MoveChartAction)).
			Post("/initDefaultApp", new(InitDefaultAppAction)).
			GetPost("/chart", new(ChartAction)).
			Get("/exportChartData", new(ExportChartDataAction)).
			Post("/updateChart", new(UpdateChartAction)).
			EndAll()
	})
}
