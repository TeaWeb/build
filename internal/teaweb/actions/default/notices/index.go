package notices

import (
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/TeaWeb/build/internal/teaconfigs/agents"
	"github.com/TeaWeb/build/internal/teaconfigs/notices"
	"github.com/TeaWeb/build/internal/teadb"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/utils/time"
	"math"
	"time"
)

type IndexAction actions.Action

// 通知
func (this *IndexAction) Run(params struct {
	Read int
	Page int
}) {
	this.Data["isRead"] = params.Read > 0

	count := 0
	countUnread, err := teadb.NoticeDAO().CountAllUnreadNotices()
	if err != nil {
		logs.Error(err)
	}
	if params.Read == 0 {
		count = countUnread
	} else {
		count, err = teadb.NoticeDAO().CountAllReadNotices()
		if err != nil {
			logs.Error(err)
		}
	}

	this.Data["countUnread"] = countUnread
	this.Data["count"] = count
	this.Data["soundOn"] = notices.SharedNoticeSetting().SoundOn

	// 分页
	if params.Page < 1 {
		params.Page = 1
	}
	pageSize := 10
	this.Data["page"] = params.Page
	if count > 0 {
		this.Data["countPages"] = int(math.Ceil(float64(count) / float64(pageSize)))
	} else {
		this.Data["countPages"] = 0
	}

	// 读取数据
	ones, err := teadb.NoticeDAO().ListNotices(params.Read == 1, (params.Page-1)*pageSize, pageSize)
	if err != nil {
		logs.Error(err)
		this.Data["notices"] = []maps.Map{}
	} else {
		this.Data["notices"] = lists.Map(ones, func(k int, v interface{}) interface{} {
			notice := v.(*notices.Notice)
			isAgent := len(notice.Agent.AgentId) > 0
			isProxy := len(notice.Proxy.ServerId) > 0
			m := maps.Map{
				"id":       notice.Id,
				"isAgent":  isAgent,
				"isProxy":  isProxy,
				"isRead":   notice.IsRead,
				"message":  notice.Message,
				"datetime": timeutil.Format("Y-m-d H:i:s", time.Unix(notice.Timestamp, 0)),
			}

			// Agent
			if isAgent {
				m["level"] = notices.FindNoticeLevel(notice.Agent.Level)

				links := []maps.Map{}
				agent := agents.NewAgentConfigFromId(notice.Agent.AgentId)
				if agent != nil {
					links = append(links, maps.Map{
						"name": agent.Name,
						"url":  "/agents/board?agentId=" + agent.Id,
					})

					app := agent.FindApp(notice.Agent.AppId)
					if app != nil {
						links = append(links, maps.Map{
							"name": app.Name,
							"url":  "/agents/apps/detail?agentId=" + agent.Id + "&appId=" + app.Id,
						})

						// item
						if len(notice.Agent.ItemId) > 0 {
							item := app.FindItem(notice.Agent.ItemId)
							if item != nil {
								links = append(links, maps.Map{
									"name": item.Name,
									"url":  "/agents/apps/itemDetail?agentId=" + agent.Id + "&appId=" + app.Id + "&itemId=" + item.Id,
								})
							}
						}

						// task
						if len(notice.Agent.TaskId) > 0 {
							task := app.FindTask(notice.Agent.TaskId)
							if task != nil {
								links = append(links, maps.Map{
									"name": task.Name,
									"url":  "/agents/apps/itemDetail?agentId=" + agent.Id + "&appId=" + app.Id + "&taskId=" + task.Id,
								})
							}
						}
					}
				}

				m["links"] = links
			}

			// Proxy
			if isProxy {
				m["level"] = notices.FindNoticeLevel(notice.Proxy.Level)

				links := []maps.Map{}
				server := teaconfigs.NewServerConfigFromId(notice.Proxy.ServerId)
				if server != nil {
					links = append(links, maps.Map{
						"name": server.Description,
						"url":  "/proxy/board?serverId=" + server.Id,
					})
				}

				if len(notice.Proxy.BackendId) > 0 {
					if len(notice.Proxy.LocationId) > 0 {
						location := server.FindLocation(notice.Proxy.LocationId)
						if location != nil {
							links = append(links, maps.Map{
								"name": location.Pattern,
								"url":  "/proxy/locations/detail?serverId=" + server.Id + "&locationId=" + notice.Proxy.LocationId,
							})
							if notice.Proxy.Websocket {
								links = append(links, maps.Map{
									"name": "Websocket",
									"url":  "/proxy/locations/websocket?serverId=" + server.Id + "&locationId=" + notice.Proxy.LocationId,
								})
								links = append(links, maps.Map{
									"name": "后端服务器",
									"url":  "/proxy/locations/websocket?serverId=" + server.Id + "&locationId=" + notice.Proxy.LocationId,
								})
							} else {
								links = append(links, maps.Map{
									"name": "后端服务器",
									"url":  "/proxy/locations/backends?serverId=" + server.Id + "&locationId=" + notice.Proxy.LocationId,
								})
							}
						}
					} else {
						links = append(links, maps.Map{
							"name": "后端服务器",
							"url":  "/proxy/backend?serverId=" + server.Id,
						})
					}
				}

				m["links"] = links
			}

			return m
		})
	}

	this.Show()
}
