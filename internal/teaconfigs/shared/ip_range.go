package shared

import (
	"bytes"
	"errors"
	"github.com/iwind/TeaGo/utils/string"
	"net"
	"regexp"
	"strings"
)

// IP Range类型
type IPRangeType = int

const (
	IPRangeTypeRange    IPRangeType = 1
	IPRangeTypeCIDR     IPRangeType = 2
	IPRangeTypeAll      IPRangeType = 3
	IPRangeTypeWildcard IPRangeType = 4 // 通配符，可以使用*
)

// IP Range
type IPRangeConfig struct {
	Id string `yaml:"id" json:"id"`

	Type IPRangeType `yaml:"type" json:"type"`

	Param  string `yaml:"param" json:"param"`
	CIDR   string `yaml:"cidr" json:"cidr"`
	IPFrom string `yaml:"ipFrom" json:"ipFrom"`
	IPTo   string `yaml:"ipTo" json:"ipTo"`

	cidr   *net.IPNet
	ipFrom net.IP
	ipTo   net.IP
	reg    *regexp.Regexp
}

// 获取新对象
func NewIPRangeConfig() *IPRangeConfig {
	return &IPRangeConfig{
		Id: stringutil.Rand(16),
	}
}

// 从字符串中分析
func ParseIPRange(s string) (*IPRangeConfig, error) {
	if len(s) == 0 {
		return nil, errors.New("invalid ip range")
	}

	ipRange := &IPRangeConfig{}

	if s == "*" || s == "all" || s == "ALL" || s == "0.0.0.0" {
		ipRange.Type = IPRangeTypeAll
		return ipRange, nil
	}

	if strings.Contains(s, "/") {
		ipRange.Type = IPRangeTypeCIDR
		ipRange.CIDR = strings.Replace(s, " ", "", -1)
	} else if strings.Contains(s, "-") {
		ipRange.Type = IPRangeTypeRange
		pieces := strings.SplitN(s, "-", 2)
		ipRange.IPFrom = strings.TrimSpace(pieces[0])
		ipRange.IPTo = strings.TrimSpace(pieces[1])
	} else if strings.Contains(s, ",") {
		ipRange.Type = IPRangeTypeRange
		pieces := strings.SplitN(s, ",", 2)
		ipRange.IPFrom = strings.TrimSpace(pieces[0])
		ipRange.IPTo = strings.TrimSpace(pieces[1])
	} else if strings.Contains(s, "*") {
		ipRange.Type = IPRangeTypeWildcard
		s = "^" + strings.Replace(regexp.QuoteMeta(s), `\*`, `\d+`, -1) + "$"
		ipRange.reg = regexp.MustCompile(s)
	} else {
		ipRange.Type = IPRangeTypeRange
		ipRange.IPFrom = s
		ipRange.IPTo = s
	}

	err := ipRange.Validate()
	if err != nil {
		return nil, err
	}
	return ipRange, nil
}

// 校验
func (this *IPRangeConfig) Validate() error {
	if this.Type == IPRangeTypeCIDR {
		if len(this.CIDR) == 0 {
			return errors.New("cidr should not be empty")
		}

		_, cidr, err := net.ParseCIDR(this.CIDR)
		if err != nil {
			return err
		}
		this.cidr = cidr
	}

	if this.Type == IPRangeTypeRange {
		this.ipFrom = net.ParseIP(this.IPFrom)
		this.ipTo = net.ParseIP(this.IPTo)

		if this.ipFrom.To4() == nil && this.ipFrom.To16() == nil {
			return errors.New("from ip should in IPv4 or IPV6 format")
		}

		if this.ipTo.To4() == nil && this.ipTo.To16() == nil {
			return errors.New("to ip should in IPv4 or IPV6 format")
		}
	}

	return nil
}

// 是否包含某个IP
func (this *IPRangeConfig) Contains(ipString string) bool {
	ip := net.ParseIP(ipString)
	if ip.To4() == nil {
		return false
	}
	if this.Type == IPRangeTypeCIDR {
		if this.cidr == nil {
			return false
		}
		return this.cidr.Contains(ip)
	}
	if this.Type == IPRangeTypeRange {
		if this.ipFrom == nil || this.ipTo == nil {
			return false
		}
		return bytes.Compare(ip, this.ipFrom) >= 0 && bytes.Compare(ip, this.ipTo) <= 0
	}
	if this.Type == IPRangeTypeWildcard {
		if this.reg == nil {
			return false
		}
		return this.reg.MatchString(ipString)
	}
	if this.Type == IPRangeTypeAll {
		return true
	}
	return false
}
