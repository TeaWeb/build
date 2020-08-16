package teaconfigs

import (
	"github.com/TeaWeb/build/internal/teautils"
	"strconv"
	"strings"
)

// HSTS设置
// 参考： https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Strict-Transport-Security
type HSTSConfig struct {
	On                bool     `yaml:"on" json:"on"`
	MaxAge            int      `yaml:"maxAge" json:"maxAge"` // 单位秒
	IncludeSubDomains bool     `yaml:"includeSubDomains" json:"includeSubDomains"`
	Preload           bool     `yaml:"preload" json:"preload"`
	Domains           []string `yaml:"domains" json:"domains"`

	hasDomains  bool
	headerValue string
}

// 校验
func (this *HSTSConfig) Validate() error {
	this.hasDomains = len(this.Domains) > 0
	this.headerValue = this.asHeaderValue()
	return nil
}

// 判断是否匹配域名
func (this *HSTSConfig) Match(domain string) bool {
	if !this.hasDomains {
		return true
	}
	return teautils.MatchDomains(this.Domains, domain)
}

// Header Key
func (this *HSTSConfig) HeaderKey() string {
	return "Strict-Transport-Security"
}

// 取得当前的Header值
func (this *HSTSConfig) HeaderValue() string {
	return this.headerValue
}

// 转换为Header值
func (this *HSTSConfig) asHeaderValue() string {
	b := strings.Builder{}
	b.WriteString("max-age=")
	if this.MaxAge > 0 {
		b.WriteString(strconv.Itoa(this.MaxAge))
	} else {
		b.WriteString("31536000") // 1 year
	}
	if this.IncludeSubDomains {
		b.WriteString("; includeSubDomains")
	}
	if this.Preload {
		b.WriteString("; preload")
	}
	return b.String()
}
