package teaplugins

type Plugin struct {
	IsExternal bool // 是否第三方开发的

	Name        string // 名称
	Code        string // 代号
	Version     string // 版本
	Date        string // 发布日期
	Site        string // 网站链接
	Developer   string // 开发者
	Description string // 插件简介

	HasRequestFilter  bool
	HasResponseFilter bool
}

func NewPlugin() *Plugin {
	return &Plugin{
	}
}
