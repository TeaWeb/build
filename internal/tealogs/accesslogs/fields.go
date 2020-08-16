package accesslogs

import "github.com/iwind/TeaGo/maps"

type AccessLogField = int

const (
	AccessLogFieldHeader       = 1
	AccessLogFieldSentHeader   = 2
	AccessLogFieldArg          = 3
	AccessLogFieldCookie       = 4
	AccessLogFieldExtend       = 5
	AccessLogFieldReferer      = 6
	AccessLogFieldUserAgent    = 7
	AccessLogFieldRequestBody  = 8
	AccessLogFieldResponseBody = 9
)

var AccessLogFieldsCodes = []int{
	AccessLogFieldHeader,
	AccessLogFieldSentHeader,
	AccessLogFieldArg,
	AccessLogFieldCookie,
	AccessLogFieldExtend,
	AccessLogFieldReferer,
	AccessLogFieldUserAgent,
	AccessLogFieldRequestBody,
	AccessLogFieldResponseBody,
}

var AccessLogDefaultFieldsCodes = []int{
	AccessLogFieldHeader,
	AccessLogFieldSentHeader,
	AccessLogFieldArg,
	AccessLogFieldCookie,
	AccessLogFieldExtend,
	AccessLogFieldReferer,
	AccessLogFieldUserAgent,
}

var AccessLogFields = []maps.Map{
	{
		"code": AccessLogFieldHeader,
		"name": "请求Header列表",
	},
	{
		"code": AccessLogFieldSentHeader,
		"name": "响应Header列表",
	},
	{
		"code": AccessLogFieldArg,
		"name": "参数列表",
	},
	{
		"code": AccessLogFieldCookie,
		"name": "Cookie列表",
	},
	{
		"code": AccessLogFieldExtend,
		"name": "位置和浏览器分析",
	},
	{
		"code": AccessLogFieldReferer,
		"name": "请求来源",
	},
	{
		"code": AccessLogFieldUserAgent,
		"name": "终端信息",
	},
	{
		"code": AccessLogFieldRequestBody,
		"name": "请求Body",
	},
	{
		"code": AccessLogFieldResponseBody,
		"name": "响应Body",
	},
}
