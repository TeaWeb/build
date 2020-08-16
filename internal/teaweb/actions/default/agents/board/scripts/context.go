package scripts

import (
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/TeaWeb/build/internal/teaconfigs/agents"
)

// 脚本执行的上下文
type Context struct {
	Agent *agents.AgentConfig // Agent
	App   *agents.AppConfig   // App
	Item  *agents.Item        // Item

	TimeType string              // 时间类型
	TimePast teaconfigs.TimePast // 过去时间
	TimeUnit teaconfigs.TimeUnit // 时间单位，用于分隔X轴
	DayFrom  string              // 开始日期
	DayTo    string              // 结束日期
}
