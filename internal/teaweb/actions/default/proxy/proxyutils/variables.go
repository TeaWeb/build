package proxyutils

import "github.com/iwind/TeaGo/maps"

// 默认的请求变量列表
func DefaultRequestVariables() []maps.Map {
	return []maps.Map{
		{"code": "${teaVersion}", "name": "TeaWeb版本", "description": ""},
		{"code": "${remoteAddr}", "name": "客户端地址（IP）", "description": "会依次根据X-Forwarded-For、X-Real-IP、RemoteAddr获取"},
		{"code": "${rawRemoteAddr}", "name": "客户端地址（IP）", "description": "返回直接连接服务的客户端原始IP地址，从v0.1.3版本加入"},
		{"code": "${remotePort}", "name": "客户端端口", "description": ""},
		{"code": "${remoteUser}", "name": "客户端用户名", "description": ""},
		{"code": "${requestURI}", "name": "请求URI", "description": ""},
		{"code": "${requestPath}", "name": "请求路径（不包括参数）", "description": ""},
		{"code": "${requestLength}", "name": "请求内容长度", "description": ""},
		{"code": "${requestMethod}", "name": "请求方法", "description": "比如GET、POST"},
		{"code": "${requestFilename}", "name": "请求文件路径", "description": ""},
		{"code": "${scheme}", "name": "请求协议，http或https", "description": ""},
		{"code": "${proto}", "name": "包含版本的HTTP请求协议", "description:": "类似于HTTP/1.0"},
		{"code": "${timeISO8601}", "name": "ISO 8601格式的时间", "description": "比如2018-07-16T23:52:24.839+08:00"},
		{"code": "${timeLocal}", "name": "本地时间", "description": "比如17/Jul/2018:09:52:24 +0800"},
		{"code": "${msec}", "name": "带有毫秒的时间", "description": "比如1531756823.054"},
		{"code": "${timestamp}", "name": "unix时间戳，单位为秒", "description": ""},
		{"code": "${host}", "name": "主机名", "description": ""},
		{"code": "${serverName}", "name": "接收请求的服务器名", "description": ""},
		{"code": "${serverPort}", "name": "接收请求的服务器端口", "description": ""},
		{"code": "${referer}", "name": "请求来源URL", "description": ""},
		{"code": "${userAgent}", "name": "客户端信息", "description": ""},
		{"code": "${contentType}", "name": "请求头部的Content-Type", "description": ""},
		{"code": "${cookies}", "name": "所有cookie组合字符串", "description": ""},
		{"code": "${cookie.NAME}", "name": "单个cookie值", "description": ""},
		{"code": "${args}", "name": "所有参数组合字符串", "description": ""},
		{"code": "${arg.NAME}", "name": "单个参数值", "description": ""},
		{"code": "${headers}", "name": "所有Header信息组合字符串", "description": ""},
		{"code": "${header.NAME}", "name": "单个Header值", "description": ""},
	}
}
