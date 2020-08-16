package agentutils

import (
	"fmt"
	"github.com/TeaWeb/build/internal/teaconfigs/agents"
	"github.com/TeaWeb/build/internal/teadb"
	"github.com/TeaWeb/build/internal/teaweb/utils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/logs"
	"net/http"
)

func AddTabbar(actionWrapper actions.ActionWrapper) {
	if actionWrapper.Object().Request.Method != http.MethodGet {
		return
	}

	action := actionWrapper.Object()
	action.Data["teaMenu"] = "agents"

	// 子菜单
	menuGroup := utils.NewMenuGroup()

	isIndex := !action.HasPrefix("/agents/addAgent", "/agents/cluster/add", "/agents/groups")

	if isIndex {
		agentId := action.ParamString("agentId")
		if len(agentId) == 0 {
			agentId = "local"
		}

		actionCode := "board"
		if action.HasPrefix("/agents/apps") {
			actionCode = "apps"
		} else if action.HasPrefix("/agents/settings") {
			actionCode = "settings"
		} else if action.HasPrefix("/agents/delete") {
			actionCode = "delete"
		} else if action.HasPrefix("/agents/notices") {
			actionCode = "notices"
		}

		topSubName := ""
		if lists.ContainsAny([]string{"/agents/board", "/agents/menu"}, action.Request.URL.Path) {
			topSubName = ""
		}
		defaultGroup := agents.SharedGroupList().FindGroup("default")

		state := FindAgentState("local")
		menu := menuGroup.FindMenu("default", defaultGroup.Name+topSubName)
		if state.IsActive {
			subName := "已连接"
			if state != nil && len(state.OsName) > 0 {
				subName = state.OsName
			}

			item := menu.Add("本地", subName, "/agents/"+actionCode+"?agentId=local", agentId == "local" && !action.HasPrefix("/agents/addAgent", "/agents/cluster/add", "/agents/groups"))
			item.SubColor = "green"
		} else {
			menu.Add("本地", "", "/agents/"+actionCode+"?agentId=local", agentId == "local" && !action.HasPrefix("/agents/addAgent", "/agents/cluster/add", "/agents/groups"))
		}

		// agent列表
		allAgents := agents.SharedAgents()

		counterMapping := map[string]int{} // groupId => count
		maxCount := 50
		for _, agent := range allAgents {
			state := FindAgentState(agent.Id)

			var menu *utils.Menu = nil
			if len(agent.GroupIds) > 0 {
				group := agents.SharedGroupList().FindGroup(agent.GroupIds[0])
				if group == nil {
					menu = menuGroup.FindMenu("default", defaultGroup.Name+topSubName)
				} else {
					menu = menuGroup.FindMenu(group.Id, group.Name)
					menu.Index = group.Index

					// 计算数量
					_, found := counterMapping[group.Id]
					if found {
						counterMapping[group.Id] ++
					} else {
						counterMapping[group.Id] = 1
					}
					if counterMapping[group.Id] > maxCount {
						if counterMapping[group.Id] == maxCount+1 {
							menu.AddSpecial("[更多主机]", "", "/agents/groups/detail?groupId="+group.Id, false)
						}
						continue
					}
				}
			} else {
				menu = menuGroup.FindMenu("default", defaultGroup.Name+topSubName)

				// 计算数量
				groupId := ""
				_, found := counterMapping[groupId]
				if found {
					counterMapping[groupId] ++
				} else {
					counterMapping[groupId] = 1
				}
				if counterMapping[groupId] > maxCount {
					if counterMapping[groupId] == maxCount+1 {
						menu.AddSpecial("[更多主机]", "", "/agents/groups/detail?groupId="+groupId, false)
					}
					continue
				}
			}

			if state.IsActive {
				subName := "已连接"
				if state != nil && len(state.OsName) > 0 {
					subName = state.OsName
				}
				item := menu.Add(agent.Name, subName, "/agents/"+actionCode+"?agentId="+agent.Id, agentId == agent.Id)
				item.Id = agent.Id
				item.IsSortable = true
				item.SubColor = "green"
			} else if !agent.On {
				item := menu.Add(agent.Name, "未启用", "/agents/"+actionCode+"?agentId="+agent.Id, agentId == agent.Id)
				item.Id = agent.Id
				item.IsSortable = true
			} else {
				item := menu.Add(agent.Name, "已断开连接", "/agents/"+actionCode+"?agentId="+agent.Id, agentId == agent.Id)
				item.Id = agent.Id
				item.IsSortable = true
				item.SubColor = "red"
			}
		}

		// Tabbar
		if !action.HasPrefix("/agents/addAgent", "/agents/cluster/add", "/agents/groups") {
			agent := agents.NewAgentConfigFromId(agentId)
			if agent != nil {
				tabbar := utils.NewTabbar()

				// 看板和Apps
				tabbar.Add("看板", "", "/agents/board?agentId="+agentId, "dashboard", action.HasPrefix("/agents/board"))
				tabbar.Add("Apps", fmt.Sprintf("%d", len(agent.Apps)), "/agents/apps?agentId="+agentId, "gem outline", action.HasPrefix("/agents/apps"))

				// 通知
				countUnreadNotices, err := teadb.NoticeDAO().CountUnreadNoticesForAgent(agentId)
				if err != nil && err != teadb.ErrorDBUnavailable {
					logs.Error(err)
				}
				if countUnreadNotices > 0 {
					tabbar.Add("通知", fmt.Sprintf("%d", countUnreadNotices), "/agents/notices?agentId="+agentId, "bell blink orange", action.HasPrefix("/agents/notices"))
				} else {
					tabbar.Add("通知", fmt.Sprintf("%d", countUnreadNotices), "/agents/notices?agentId="+agentId, "bell", action.HasPrefix("/agents/notices"))
				}

				// 设置和删除
				tabbar.Add("设置", "", "/agents/settings?agentId="+agentId, "setting", action.HasPrefix("/agents/settings"))
				if agentId != "local" {
					tabbar.Add("删除", "", "/agents/delete?agentId="+agentId, "trash", action.HasPrefix("/agents/delete"))
				}
				utils.SetTabbar(actionWrapper, tabbar)
			}
		}
	}

	// 操作按钮
	{
		menu := menuGroup.FindMenu("operations", "[操作]")
		menu.AlwaysActive = true
		menuGroup.AlwaysMenu = menu
		menu.Index = 10000
		menu.Add("主机", "", "/agents", isIndex)
		menu.Add("+添加新主机", "", "/agents/addAgent", action.HasPrefix("/agents/addAgent", "/agents/cluster/add"))
		menu.Add("分组管理", "", "/agents/groups", action.HasPrefix("/agents/groups"))
	}

	menuGroup.Sort()
	utils.SetSubMenu(action, menuGroup)
}
