package shared

import (
	"github.com/iwind/TeaGo/rands"
)

// 客户端配置
type ClientConfig struct {
	Id          string `yaml:"id" json:"id"`                   // ID
	On          bool   `yaml:"on" json:"on"`                   // 是否开启
	IP          string `yaml:"ip" json:"ip"`                   // IP
	Description string `yaml:"description" json:"description"` // 描述

	ipRange *IPRangeConfig
}

// 取得新配置对象
func NewClientConfig() *ClientConfig {
	return &ClientConfig{
		Id: rands.HexString(16),
		On: true,
	}
}

// 校验
func (this *ClientConfig) Validate() error {
	if len(this.IP) > 0 {
		ipRange, err := ParseIPRange(this.IP)
		if err != nil {
			return err
		}
		this.ipRange = ipRange
	}
	return nil
}

// 判断是否匹配某个IP
func (this *ClientConfig) Match(ip string) bool {
	if len(ip) == 0 || this.ipRange == nil {
		return false
	}
	return this.ipRange.Contains(ip)
}
