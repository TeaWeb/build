package scheduling

import (
	"github.com/TeaWeb/build/internal/teaconfigs/shared"
	"github.com/iwind/TeaGo/maps"
)

// 调度算法接口
type SchedulingInterface interface {
	// 是否有候选对象
	HasCandidates() bool

	// 添加候选对象
	Add(candidate ...CandidateInterface)

	// 启动
	Start()

	// 查找下一个候选对象
	Next(call *shared.RequestCall) CandidateInterface

	// 获取简要信息
	Summary() maps.Map
}

// 调度算法基础类
type Scheduling struct {
	Candidates []CandidateInterface
}

// 判断是否有候选对象
func (this *Scheduling) HasCandidates() bool {
	return len(this.Candidates) > 0
}

// 添加候选对象
func (this *Scheduling) Add(candidate ...CandidateInterface) {
	this.Candidates = append(this.Candidates, candidate ...)
}
