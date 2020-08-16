package database

import (
	"github.com/TeaWeb/build/internal/teadb"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/maps"
	"regexp"
)

type TablesAction actions.Action

// 集合列表
func (this *TablesAction) Run(params struct{}) {
	tables, err := teadb.SharedDB().ListTables()
	if err != nil {
		this.Fail("数据库查询错误：" + err.Error())
	}

	// 排序
	result := []maps.Map{}
	recognizedNames := []string{}

	// 日志
	{
		reg := regexp.MustCompile("^(?:logs\\.|teaweb_logs_)(\\d{8})$")
		for _, name := range tables {
			matches := reg.FindStringSubmatch(name)
			if len(matches) == 0 {
				continue
			}
			recognizedNames = append(recognizedNames, name)
			result = append(result, maps.Map{
				"name":      name,
				"type":      "代理访问日志",
				"canDelete": true,
				"subName":   matches[1][:4] + "-" + matches[1][4:6] + "-" + matches[1][6:],
			})
		}
	}

	// 统计
	{
		reg := regexp.MustCompile("^(values\\.server\\.|teaweb_values_server_)")
		for _, name := range tables {
			if !reg.MatchString(name) {
				continue
			}
			recognizedNames = append(recognizedNames, name)
			result = append(result, maps.Map{
				"name":      name,
				"type":      "代理统计数据",
				"canDelete": true,
			})
		}
	}

	// 监控数据
	{
		reg := regexp.MustCompile("^(values\\.agent\\.|teaweb_values_agent_)")
		for _, name := range tables {
			if !reg.MatchString(name) {
				continue
			}
			recognizedNames = append(recognizedNames, name)
			result = append(result, maps.Map{
				"name":      name,
				"type":      "主机监控数据",
				"canDelete": true,
			})
		}
	}

	// 监控日志
	{
		reg := regexp.MustCompile("^(logs\\.agent\\.|teaweb_logs_agent_)")
		for _, name := range tables {
			if !reg.MatchString(name) {
				continue
			}
			recognizedNames = append(recognizedNames, name)
			result = append(result, maps.Map{
				"name":      name,
				"type":      "主机任务运行日志",
				"canDelete": true,
			})
		}
	}

	// 通知
	{
		reg := regexp.MustCompile("^(notices|teaweb_notices)$")
		for _, name := range tables {
			if !reg.MatchString(name) {
				continue
			}
			recognizedNames = append(recognizedNames, name)
			result = append(result, maps.Map{
				"name":      name,
				"type":      "通知提醒",
				"canDelete": true,
			})
		}
	}

	// 审计日志
	{
		reg := regexp.MustCompile("^(logs\\.audit|teaweb_logs_audit)$")
		for _, name := range tables {
			if !reg.MatchString(name) {
				continue
			}
			recognizedNames = append(recognizedNames, name)
			result = append(result, maps.Map{
				"name":      name,
				"type":      "审计日志（操作日志）",
				"canDelete": true,
			})
		}
	}

	// 旧的统计数据
	{
		reg := regexp.MustCompile("^(stats\\.|teaweb_stats_)")
		for _, name := range tables {
			if !reg.MatchString(name) {
				continue
			}
			recognizedNames = append(recognizedNames, name)
			result = append(result, maps.Map{
				"name":      name,
				"type":      "旧的统计数据",
				"canDelete": true,
				"warning":   true,
			})
		}
	}

	// 测试
	{
		reg := regexp.MustCompile("^(test\\.|teaweb_test)")
		for _, name := range tables {
			if !reg.MatchString(name) {
				continue
			}
			recognizedNames = append(recognizedNames, name)
			result = append(result, maps.Map{
				"name":      name,
				"type":      "测试表",
				"canDelete": true,
				"warning":   true,
			})
		}
	}

	// 其他
	for _, name := range tables {
		if lists.ContainsString(recognizedNames, name) {
			continue
		}
		result = append(result, maps.Map{
			"name":      name,
			"type":      "无法识别，请报告官方",
			"canDelete": true,
			"warning":   true,
		})
	}

	this.Data["tables"] = result

	this.Success()
}
