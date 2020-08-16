package scheduling

// 候选对象接口
type CandidateInterface interface {
	// 权重
	CandidateWeight() uint

	// 代号
	CandidateCodes() []string
}
