package search

import (
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/TeaWeb/build/internal/teaconfigs/agents"
	"github.com/TeaWeb/build/internal/teaconfigs/db"
	"github.com/TeaWeb/build/internal/teautils"
	"github.com/TeaWeb/build/internal/teaweb/utils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type IndexAction actions.Action

// 全局搜索页
func (this *IndexAction) RunGet(params struct {
	From string
}) {
	this.Data["teaMenu"] = "search"
	this.Data["from"] = params.From

	// 操作按钮
	menuGroup := utils.NewMenuGroup()
	{
		menu := menuGroup.FindMenu("operations", "[操作]")
		menu.AlwaysActive = true
		menuGroup.AlwaysMenu = menu
		menu.Index = 10000
		menu.Add("搜索", "", "/search", true)
	}

	menuGroup.Sort()
	utils.SetSubMenu(this, menuGroup)

	this.Show()
}

// 执行搜索
func (this *IndexAction) RunPost(params struct {
	Keyword string
	Must    *actions.Must
}) {
	results := []maps.Map{} // [ type, name, tags, link ]
	if len(params.Keyword) == 0 {
		this.Data["results"] = results
		this.Success()
	}

	// proxy
	serverList, err := teaconfigs.SharedServerList()
	if err == nil {
		for _, server := range serverList.FindAllServers() {
			if matched, name, tags := server.MatchKeyword(params.Keyword); matched {
				results = append(results, maps.Map{
					"type": "代理",
					"name": name,
					"tags": tags,
					"link": "/proxy/board?serverId=" + server.Id,
				})
			}

			for _, location := range server.Locations {
				if matched, name, tags := location.MatchKeyword(params.Keyword); matched {
					if len(server.Description) > 0 {
						tags = append([]string{"代理：" + server.Description}, tags...)
					}
					results = append(results, maps.Map{
						"type": "路径规则",
						"name": name,
						"tags": tags,
						"link": "/proxy/locations/detail?serverId=" + server.Id + "&locationId=" + location.Id,
					})
				}
			}
		}
	}

	// 缓存策略
	cacheConfig, err := teaconfigs.SharedCacheConfig()
	if err == nil {
		for _, policy := range cacheConfig.FindAllPolicies() {
			if matched, name, tags := policy.MatchKeyword(params.Keyword); matched {
				results = append(results, maps.Map{
					"type": "缓存策略",
					"name": name,
					"tags": tags,
					"link": "/cache/policy?filename=" + policy.Filename,
				})
			}
		}
	}

	// WAF策略
	wafList := teaconfigs.SharedWAFList()
	for _, waf := range wafList.FindAllConfigs() {
		if matched, name, tags := waf.MatchKeyword(params.Keyword); matched {
			results = append(results, maps.Map{
				"type": "WAF",
				"name": name,
				"tags": tags,
				"link": "/proxy/waf/detail?wafId=" + waf.Id,
			})
		}
	}

	// 日志策略
	for _, policy := range teaconfigs.SharedAccessLogStoragePolicyList().FindAllPolicies() {
		if matched, name, tags := policy.MatchKeyword(params.Keyword); matched {
			results = append(results, maps.Map{
				"type": "日志策略",
				"name": name,
				"tags": tags,
				"link": "/proxy/log/policies/policy?policyId=" + policy.Id,
			})
		}
	}

	// 证书
	for _, cert := range teaconfigs.SharedSSLCertList().Certs {
		if matched, name, tags := cert.MatchKeyword(params.Keyword); matched {
			results = append(results, maps.Map{
				"type": "证书",
				"name": name,
				"tags": tags,
				"link": "/proxy/certs/detail?certId=" + cert.Id,
			})
		}
	}

	// agent
	allAgents := []*agents.AgentConfig{agents.NewAgentConfigFromId("local")}

	agentList, err := agents.SharedAgentList()
	if err == nil {
		allAgents = append(allAgents, agentList.FindAllAgents()...)
	}

	for _, agent := range allAgents {
		if matched, name, tags := agent.MatchKeyword(params.Keyword); matched {
			results = append(results, maps.Map{
				"type": "主机",
				"name": name,
				"tags": tags,
				"link": "/agents/board?agentId=" + agent.Id,
			})
		}

		for _, app := range agent.Apps {
			if matched, name, tags := app.MatchKeyword(params.Keyword); matched {
				tags = append([]string{"主机：" + agent.Name}, tags...)
				results = append(results, maps.Map{
					"type": "App",
					"name": name,
					"tags": tags,
					"link": "/agents/apps/detail?agentId=" + agent.Id + "&appId=" + app.Id,
				})
			}

			// 监控项
			for _, item := range app.Items {
				if matched, name, tags := item.MatchKeyword(params.Keyword); matched {
					tags = append([]string{"主机：" + agent.Name, "App：" + app.Name}, tags...)
					results = append(results, maps.Map{
						"type": "监控项",
						"name": name,
						"tags": tags,
						"link": "/agents/apps/itemDetail?agentId=" + agent.Id + "&appId=" + app.Id + "&itemId=" + item.Id,
					})
				}
			}

			// 任务
			for _, task := range app.Tasks {
				if matched, name, tags := task.MatchKeyword(params.Keyword); matched {
					tags = append([]string{"主机：" + agent.Name, "App：" + app.Name}, tags...)
					tabbar := "schedule"
					description := "定时任务"
					if task.IsBooting {
						tabbar = "boot"
						description = "启动任务"
					} else if task.IsManual {
						tabbar = "manual"
						description = "手动任务"
					}
					tags = append(tags, description)
					results = append(results, maps.Map{
						"type": "任务",
						"name": name,
						"tags": tags,
						"link": "/agents/apps/taskDetail?agentId=" + agent.Id + "&appId=" + app.Id + "&taskId=" + task.Id + "&tabbar=" + tabbar,
					})
				}
			}
		}
	}

	// agent分组
	for _, group := range agents.SharedGroupList().FindAllGroups() {
		if matched, name, tags := group.MatchKeyword(params.Keyword); matched {
			results = append(results, maps.Map{
				"type": "主机分组",
				"name": name,
				"tags": tags,
				"link": "/agents/groups/detail?groupId=" + group.Id,
			})
		}
	}

	// 通知
	if teautils.MatchKeyword("通知", params.Keyword) {
		results = append(results, maps.Map{
			"type": "通知",
			"name": "通知",
			"tags": []string{},
			"link": "/notices",
		})
	}

	// 设置
	if teautils.MatchKeyword("个人资料设置", params.Keyword) {
		results = append(results, maps.Map{
			"type": "设置",
			"name": "个人资料",
			"tags": []string{},
			"link": "/settings/profile",
		})
	}
	if teautils.MatchKeyword("登录设置", params.Keyword) {
		results = append(results, maps.Map{
			"type": "设置",
			"name": "登录设置",
			"tags": []string{},
			"link": "/settings/login",
		})
	}

	if teautils.MatchKeyword("数据库设置", params.Keyword) {
		results = append(results, maps.Map{
			"type": "设置",
			"name": "数据库设置",
			"tags": []string{},
			"link": "/settings/database",
		})
	}

	dbType := db.SharedDBConfig().Type
	if dbType == db.DBTypeMongo && teautils.MatchKeyword("MongoDB设置", params.Keyword) {
		results = append(results, maps.Map{
			"type": "设置",
			"name": "MongoDB设置",
			"tags": []string{},
			"link": "/settings/mongo",
		})
	}
	if dbType == db.DBTypeMySQL && teautils.MatchKeyword("MySQL设置", params.Keyword) {
		results = append(results, maps.Map{
			"type": "设置",
			"name": "MySQL配置",
			"tags": []string{},
			"link": "/settings/mysql",
		})
	}
	if dbType == db.DBTypePostgres && teautils.MatchKeyword("Postgres设置", params.Keyword) {
		results = append(results, maps.Map{
			"type": "设置",
			"name": "PostgreSQL设置",
			"tags": []string{},
			"link": "/settings/postgres",
		})
	}
	if teautils.MatchKeyword("备份Backup设置", params.Keyword) {
		results = append(results, maps.Map{
			"type": "设置",
			"name": "备份",
			"tags": []string{},
			"link": "/settings/backup",
		})
	}
	if teautils.MatchKeyword("版本更新设置", params.Keyword) {
		results = append(results, maps.Map{
			"type": "设置",
			"name": "检查版本更新",
			"tags": []string{},
			"link": "/settings/update",
		})
	}
	if teautils.MatchKeyword("ClusterNode集群设置", params.Keyword) {
		results = append(results, maps.Map{
			"type": "设置",
			"name": "集群",
			"tags": []string{},
			"link": "/settings/cluster",
		})
	}

	this.Data["results"] = results
	this.Success()
}
