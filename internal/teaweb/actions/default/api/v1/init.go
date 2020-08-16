package v1

import (
	"github.com/TeaWeb/build/internal/teaweb/actions/default/api/v1/agent"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/api/v1/agent/app"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/api/v1/agent/item"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/api/v1/agent/task"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/api/v1/backup"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/api/v1/cluster"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/api/v1/notice/media"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/api/v1/proxy"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/api/v1/proxy/accesslog"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/api/v1/proxy/stat"
	"github.com/iwind/TeaGo"
)

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		server.Prefix("/api/v1").
			Helper(new(Helper)).
			Get("/status", new(StatusAction)).
			Get("/reload", new(ReloadAction)).
			Get("/reset", new(ResetAction)).
			Get("/stop", new(StopAction)).
			Get("/proxy/servers", new(proxy.ServersAction)).
			Get("/proxy/:serverId", new(proxy.ServerAction)).
			Get("/proxy/:serverId/accesslog/latest", new(accesslog.LatestAction)).
			Get("/proxy/:serverId/accesslog/next/:lastId", new(accesslog.NextAction)).
			Get("/proxy/:serverId/accesslog/list/:size", new(accesslog.ListAction)).
			Get("/proxy/:serverId/accesslog/next/:lastId/list/:size", new(accesslog.NextListAction)).
			Get("/proxy/:serverId/stat", new(stat.IndexAction)).
			Get("/agents", new(agent.AgentsAction)).
			Get("/agent/:agentId", new(agent.AgentAction)).
			Get("/agent/:agentId/start", new(agent.StartAction)).
			Get("/agent/:agentId/stop", new(agent.StopAction)).
			Get("/agent/:agentId/delete", new(agent.DeleteAction)).
			Get("/agent/:agentId/apps", new(app.AppsAction)).
			Get("/agent/:agentId/app/:appId", new(app.AppAction)).
			Get("/agent/:agentId/app/:appId/item/:itemId", new(item.ItemAction)).
			Get("/agent/:agentId/app/:appId/item/:itemId/latest", new(item.LatestAction)).
			Get("/agent/:agentId/app/:appId/item/:itemId/execute", new(item.ExecuteAction)).
			Get("/agent/:agentId/app/:appId/task/:taskId", new(task.TaskAction)).
			Get("/agent/:agentId/app/:appId/task/:taskId/run", new(task.RunAction)).
			Get("/notice/medias", new(media.MediasAction)).
			Get("/notice/media/:mediaId", new(media.MediaAction)).
			Post("/notice/media/:mediaId/send", new(media.SendAction)).
			Get("/cluster/node", new(cluster.NodeAction)).
			Get("/cluster/push", new(cluster.PushAction)).
			Get("/cluster/pull", new(cluster.PullAction)).
			Get("/backup/files", new(backup.FilesAction)).
			Get("/backup/latest", new(backup.LatestAction)).
			Get("/backup/file/:filename", new(backup.FileAction)).
			Get("/backup/file/:filename/restore", new(backup.RestoreAction)).
			Get("/backup/file/:filename/delete", new(backup.DeleteAction)).
			EndAll()
	})
}
