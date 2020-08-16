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
			Prefix("/agents/apps").
			Helper(new(Helper)).
			Get("", new(IndexAction)).
			GetPost("/add", new(AddAction)).
			GetPost("/update", new(UpdateAction)).
			Post("/delete", new(DeleteAction)).
			Post("/move", new(MoveAction)).
			Post("/on", new(OnAction)).
			Post("/off", new(OffAction)).
			Get("/detail", new(DetailAction)).
			Get("/schedule", new(ScheduleAction)).
			Get("/boot", new(BootAction)).
			Get("/manual", new(ManualAction)).
			GetPost("/addTask", new(AddTaskAction)).
			Post("/deleteTask", new(DeleteTaskAction)).
			Get("/taskDetail", new(TaskDetailAction)).
			GetPost("/updateTask", new(UpdateTaskAction)).
			Get("/runTask", new(RunTaskAction)).
			Post("/taskOn", new(TaskOnAction)).
			Post("/taskOff", new(TaskOffAction)).
			GetPost("/taskLogs", new(TaskLogsAction)).
			GetPost("/monitor", new(MonitorAction)).
			GetPost("/addItem", new(AddItemAction)).
			Post("/deleteItem", new(DeleteItemAction)).
			Get("/itemDetail", new(ItemDetailAction)).
			Post("/itemOn", new(ItemOnAction)).
			Post("/itemOff", new(ItemOffAction)).
			Post("/execItemSource", new(ExecItemSourceAction)).
			GetPost("/updateItem", new(UpdateItemAction)).
			GetPost("/itemValues", new(ItemValuesAction)).
			GetPost("/itemCharts", new(ItemChartsAction)).
			GetPost("/addItemChart", new(AddItemChartAction)).
			Post("/deleteItemChart", new(DeleteItemChartAction)).
			GetPost("/updateItemChart", new(UpdateItemChartAction)).
			Post("/previewItemChart", new(PreviewItemChartAction)).
			Post("/clearItemValues", new(ClearItemValuesAction)).
			Post("/addDefaultCharts", new(AddDefaultChartsAction)).
			Get("/noticeReceivers", new(NoticeReceiversAction)).
			GetPost("/addNoticeReceivers", new(AddNoticeReceiversAction)).
			Post("/deleteNoticeReceivers", new(DeleteNoticeReceiversAction)).
			EndAll()
	})
}
