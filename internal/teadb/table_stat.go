package teadb

// 表格统计
type TableStat struct {
	Count         int64  `yaml:"count" json:"count"`
	Size          int64  `yaml:"size" json:"size"`
	FormattedSize string `yaml:"formattedSize" json:"formattedSize"`
}
