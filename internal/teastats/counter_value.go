package teastats

// 数值增长型的统计值
type CounterValue struct {
	Timestamp int64                  `json:"timestamp"` // 时间戳
	Params    map[string]string      `json:"params"`    // 参数，用来区分单个统计项内的不同的项目
	Value     map[string]interface{} `json:"value"`     // 数值
}
