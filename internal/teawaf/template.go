package teawaf

// 感谢以下规则来源：
// - Janusec: https://www.janusec.com/
func Template() *WAF {
	waf := NewWAF()
	waf.Id = "template"
	waf.On = true

	// black list
	{
		group := NewRuleGroup()
		group.On = false
		group.IsInbound = true
		group.Name = "白名单"
		group.Code = "whiteList"
		group.Description = "在此名单中的IP地址可以直接跳过防火墙设置"

		{

			set := NewRuleSet()
			set.On = true
			set.Name = "IP白名单"
			set.Code = "9001"
			set.Connector = RuleConnectorOr
			set.Action = ActionAllow
			set.AddRule(&Rule{
				Param:             "${remoteAddr}",
				Operator:          RuleOperatorMatch,
				Value:             `127\.0\.0\.1|0\.0\.0\.0`,
				IsCaseInsensitive: false,
			})
			group.AddRuleSet(set)
		}

		waf.AddRuleGroup(group)
	}

	// black list
	{
		group := NewRuleGroup()
		group.On = false
		group.IsInbound = true
		group.Name = "黑名单"
		group.Code = "blackList"
		group.Description = "在此名单中的IP地址直接阻止"

		{

			set := NewRuleSet()
			set.On = true
			set.Name = "IP黑名单"
			set.Code = "10001"
			set.Connector = RuleConnectorOr
			set.Action = ActionBlock
			set.AddRule(&Rule{
				Param:             "${remoteAddr}",
				Operator:          RuleOperatorMatch,
				Value:             `1\.1\.1\.1|2\.2\.2\.2`,
				IsCaseInsensitive: false,
			})
			group.AddRuleSet(set)
		}

		waf.AddRuleGroup(group)
	}

	// xss
	{
		group := NewRuleGroup()
		group.On = true
		group.IsInbound = true
		group.Name = "XSS"
		group.Code = "xss"
		group.Description = "防跨站脚本攻击（Cross Site Scripting）"

		{
			set := NewRuleSet()
			set.On = true
			set.Name = "Javascript事件"
			set.Code = "1001"
			set.Connector = RuleConnectorOr
			set.Action = ActionBlock
			set.AddRule(&Rule{
				Param:             "${requestURI}",
				Operator:          RuleOperatorMatch,
				Value:             `(onmouseover|onmousemove|onmousedown|onmouseup|onerror|onload|onclick|ondblclick|onkeydown|onkeyup|onkeypress)\s*=`, // TODO more keywords here
				IsCaseInsensitive: true,
			})
			group.AddRuleSet(set)
		}

		{
			set := NewRuleSet()
			set.On = true
			set.Name = "Javascript函数"
			set.Code = "1002"
			set.Connector = RuleConnectorOr
			set.Action = ActionBlock
			set.AddRule(&Rule{
				Param:             "${requestURI}",
				Operator:          RuleOperatorMatch,
				Value:             `(alert|eval|prompt|confirm)\s*\(`, // TODO more keywords here
				IsCaseInsensitive: true,
			})
			group.AddRuleSet(set)
		}

		{
			set := NewRuleSet()
			set.On = true
			set.Name = "HTML标签"
			set.Code = "1003"
			set.Connector = RuleConnectorOr
			set.Action = ActionBlock
			set.AddRule(&Rule{
				Param:             "${requestURI}",
				Operator:          RuleOperatorMatch,
				Value:             `<(script|iframe|link)`, // TODO more keywords here
				IsCaseInsensitive: true,
			})
			group.AddRuleSet(set)
		}

		waf.AddRuleGroup(group)
	}

	// upload
	{
		group := NewRuleGroup()
		group.On = true
		group.IsInbound = true
		group.Name = "文件上传"
		group.Code = "upload"
		group.Description = "防止上传可执行脚本文件到服务器"

		{
			set := NewRuleSet()
			set.On = true
			set.Name = "上传文件扩展名"
			set.Code = "2001"
			set.Connector = RuleConnectorOr
			set.Action = ActionBlock
			set.AddRule(&Rule{
				Param:             "${requestUpload.ext}",
				Operator:          RuleOperatorMatch,
				Value:             `\.(php|jsp|aspx|asp|exe|asa|rb|py)\b`, // TODO more keywords here
				IsCaseInsensitive: true,
			})
			group.AddRuleSet(set)
		}

		waf.AddRuleGroup(group)
	}

	// web shell
	{
		group := NewRuleGroup()
		group.On = true
		group.IsInbound = true
		group.Name = "Web Shell"
		group.Code = "webShell"
		group.Description = "防止远程执行服务器命令"

		{
			set := NewRuleSet()
			set.On = true
			set.Name = "Web Shell"
			set.Code = "3001"
			set.Connector = RuleConnectorOr
			set.Action = ActionBlock
			set.AddRule(&Rule{
				Param:             "${requestAll}",
				Operator:          RuleOperatorMatch,
				Value:             `\b(eval|system|exec|execute|passthru|shell_exec|phpinfo)\s*\(`, // TODO more keywords here
				IsCaseInsensitive: true,
			})
			group.AddRuleSet(set)
		}

		waf.AddRuleGroup(group)
	}

	// command injection
	{
		group := NewRuleGroup()
		group.On = true
		group.IsInbound = true
		group.Name = "命令注入"
		group.Code = "commandInjection"

		{
			set := NewRuleSet()
			set.On = true
			set.Name = "命令注入"
			set.Code = "4001"
			set.Connector = RuleConnectorOr
			set.Action = ActionBlock
			set.AddRule(&Rule{
				Param:             "${requestURI}",
				Operator:          RuleOperatorMatch,
				Value:             `\b(pwd|ls|ll|whoami|id|net\s+user)\s*$`, // TODO more keywords here
				IsCaseInsensitive: false,
			})
			set.AddRule(&Rule{
				Param:             "${requestBody}",
				Operator:          RuleOperatorMatch,
				Value:             `\b(pwd|ls|ll|whoami|id|net\s+user)\s*$`, // TODO more keywords here
				IsCaseInsensitive: false,
			})
			group.AddRuleSet(set)
		}

		waf.AddRuleGroup(group)
	}

	// path traversal
	{
		group := NewRuleGroup()
		group.On = true
		group.IsInbound = true
		group.Name = "路径穿越"
		group.Code = "pathTraversal"
		group.Description = "防止读取网站目录之外的其他系统文件"

		{
			set := NewRuleSet()
			set.On = true
			set.Name = "路径穿越"
			set.Code = "5001"
			set.Connector = RuleConnectorOr
			set.Action = ActionBlock
			set.AddRule(&Rule{
				Param:             "${requestURI}",
				Operator:          RuleOperatorMatch,
				Value:             `((\.+)(/+)){2,}`, // TODO more keywords here
				IsCaseInsensitive: false,
			})
			group.AddRuleSet(set)
		}

		waf.AddRuleGroup(group)
	}

	// special dirs
	{
		group := NewRuleGroup()
		group.On = true
		group.IsInbound = true
		group.Name = "特殊目录"
		group.Code = "denyDirs"
		group.Description = "防止通过Web访问到一些特殊目录"

		{
			set := NewRuleSet()
			set.On = true
			set.Name = "特殊目录"
			set.Code = "6001"
			set.Connector = RuleConnectorOr
			set.Action = ActionBlock
			set.AddRule(&Rule{
				Param:             "${requestPath}",
				Operator:          RuleOperatorMatch,
				Value:             `/\.(git|svn|htaccess|idea)\b`, // TODO more keywords here
				IsCaseInsensitive: true,
			})
			group.AddRuleSet(set)
		}

		waf.AddRuleGroup(group)
	}

	// sql injection
	{
		group := NewRuleGroup()
		group.On = true
		group.IsInbound = true
		group.Name = "SQL注入"
		group.Code = "sqlInjection"
		group.Description = "防止SQL注入漏洞"

		{
			set := NewRuleSet()
			set.On = true
			set.Name = "Union SQL Injection"
			set.Code = "7001"
			set.Connector = RuleConnectorOr
			set.Action = ActionBlock

			set.AddRule(&Rule{
				Param:             "${requestAll}",
				Operator:          RuleOperatorMatch,
				Value:             `union[\s/\*]+select`,
				IsCaseInsensitive: true,
			})

			group.AddRuleSet(set)
		}

		{
			set := NewRuleSet()
			set.On = true
			set.Name = "SQL注释"
			set.Code = "7002"
			set.Connector = RuleConnectorOr
			set.Action = ActionBlock

			set.AddRule(&Rule{
				Param:             "${requestAll}",
				Operator:          RuleOperatorMatch,
				Value:             `/\*(!|\x00)`,
				IsCaseInsensitive: true,
			})

			group.AddRuleSet(set)
		}

		{
			set := NewRuleSet()
			set.On = true
			set.Name = "SQL条件"
			set.Code = "7003"
			set.Connector = RuleConnectorOr
			set.Action = ActionBlock

			set.AddRule(&Rule{
				Param:             "${requestAll}",
				Operator:          RuleOperatorMatch,
				Value:             `\s(and|or|rlike)\s+(if|updatexml)\s*\(`,
				IsCaseInsensitive: true,
			})
			set.AddRule(&Rule{
				Param:             "${requestAll}",
				Operator:          RuleOperatorMatch,
				Value:             `\s+(and|or|rlike)\s+(select|case)\s+`,
				IsCaseInsensitive: true,
			})
			set.AddRule(&Rule{
				Param:             "${requestAll}",
				Operator:          RuleOperatorMatch,
				Value:             `\s+(and|or|procedure)\s+[\w\p{L}]+\s*=\s*[\w\p{L}]+(\s|$|--|#)`,
				IsCaseInsensitive: true,
			})
			set.AddRule(&Rule{
				Param:             "${requestAll}",
				Operator:          RuleOperatorMatch,
				Value:             `\(\s*case\s+when\s+[\w\p{L}]+\s*=\s*[\w\p{L}]+\s+then\s+`,
				IsCaseInsensitive: true,
			})

			group.AddRuleSet(set)
		}

		{
			set := NewRuleSet()
			set.On = true
			set.Name = "SQL函数"
			set.Code = "7004"
			set.Connector = RuleConnectorOr
			set.Action = ActionBlock

			set.AddRule(&Rule{
				Param:             "${requestAll}",
				Operator:          RuleOperatorMatch,
				Value:             `(updatexml|extractvalue|ascii|ord|char|chr|count|concat|rand|floor|substr|length|len|user|database|benchmark|analyse)\s*\(`,
				IsCaseInsensitive: true,
			})

			group.AddRuleSet(set)
		}

		{
			set := NewRuleSet()
			set.On = true
			set.Name = "SQL附加语句"
			set.Code = "7005"
			set.Connector = RuleConnectorOr
			set.Action = ActionBlock

			set.AddRule(&Rule{
				Param:             "${requestAll}",
				Operator:          RuleOperatorMatch,
				Value:             `;\s*(declare|use|drop|create|exec|delete|update|insert)\s`,
				IsCaseInsensitive: true,
			})

			group.AddRuleSet(set)
		}

		waf.AddRuleGroup(group)
	}

	// bot
	{
		group := NewRuleGroup()
		group.On = false
		group.IsInbound = true
		group.Name = "网络爬虫"
		group.Code = "bot"
		group.Description = "禁止一些网络爬虫"

		{
			set := NewRuleSet()
			set.On = true
			set.Name = "常见网络爬虫"
			set.Code = "20001"
			set.Connector = RuleConnectorOr
			set.Action = ActionBlock

			set.AddRule(&Rule{
				Param:             "${userAgent}",
				Operator:          RuleOperatorMatch,
				Value:             `Googlebot|AdsBot|bingbot|BingPreview|facebookexternalhit|Slurp|Sogou|proximic|Baiduspider|yandex|twitterbot|spider|python`,
				IsCaseInsensitive: true,
			})

			group.AddRuleSet(set)
		}

		waf.AddRuleGroup(group)
	}

	// cc
	{
		group := NewRuleGroup()
		group.On = false
		group.IsInbound = true
		group.Name = "CC攻击"
		group.Description = "Challenge Collapsar，防止短时间大量请求涌入，请谨慎开启和设置"
		group.Code = "cc"

		{
			set := NewRuleSet()
			set.On = true
			set.Name = "CC请求数"
			set.Description = "限制单IP在一定时间内的请求数"
			set.Code = "8001"
			set.Connector = RuleConnectorAnd
			set.Action = ActionBlock
			set.AddRule(&Rule{
				Param:    "${cc.requests}",
				Operator: RuleOperatorGt,
				Value:    "1000",
				CheckpointOptions: map[string]string{
					"period": "60",
				},
				IsCaseInsensitive: false,
			})
			set.AddRule(&Rule{
				Param:             "${remoteAddr}",
				Operator:          RuleOperatorNotIPRange,
				Value:             `127.0.0.1/8`,
				IsCaseInsensitive: false,
			})
			set.AddRule(&Rule{
				Param:             "${remoteAddr}",
				Operator:          RuleOperatorNotIPRange,
				Value:             `192.168.0.1/16`,
				IsCaseInsensitive: false,
			})
			set.AddRule(&Rule{
				Param:             "${remoteAddr}",
				Operator:          RuleOperatorNotIPRange,
				Value:             `10.0.0.1/8`,
				IsCaseInsensitive: false,
			})
			set.AddRule(&Rule{
				Param:             "${remoteAddr}",
				Operator:          RuleOperatorNotIPRange,
				Value:             `172.16.0.1/12`,
				IsCaseInsensitive: false,
			})

			group.AddRuleSet(set)
		}

		waf.AddRuleGroup(group)
	}

	// custom
	{
		group := NewRuleGroup()
		group.On = true
		group.IsInbound = true
		group.Name = "自定义规则分组"
		group.Description = "我的自定义规则分组，可以将自定义的规则放在这个分组下"
		group.Code = "custom"
		waf.AddRuleGroup(group)
	}

	return waf
}
