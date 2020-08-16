package teaconfigs

import "github.com/TeaWeb/build/internal/teaconfigs/shared"

// IP Range列表
type IPRangeList struct {
	IPRanges []*shared.IPRangeConfig
}

// 校验
func (this *IPRangeList) Validate() error {
	for _, ipRange := range this.IPRanges {
		err := ipRange.Validate()
		if err != nil {
			return err
		}
	}
	return nil
}

// 添加
func (this *IPRangeList) AddIPRange(ipRange *shared.IPRangeConfig) {
	this.IPRanges = append(this.IPRanges, ipRange)
}

// 删除
func (this *IPRangeList) RemoveIPRange(ipRangeId string) {
	result := []*shared.IPRangeConfig{}
	for _, r := range this.IPRanges {
		if r.Id == ipRangeId {
			continue
		}
		result = append(result, r)
	}
	this.IPRanges = result
}
