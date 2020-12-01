package teaconfigs

import (
	"fmt"
	"github.com/TeaWeb/build/internal/teaconfigs/shared"
	"github.com/TeaWeb/build/internal/teautils"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/rands"
	"regexp"
	"strings"
)

const (
	RewriteTargetProxy = 1
	RewriteTargetURL   = 2
)

const (
	RewriteFlagRedirect = "r" // 跳转，TODO: 实现 302, 305
	RewriteFlagProxy    = "p" // 代理
)

// 重写规则定义
//
// 参考
// - http://nginx.org/en/docs/http/ngx_http_rewrite_module.html
// - https://httpd.apache.org/docs/current/mod/mod_rewrite.html
// - https://httpd.apache.org/docs/2.4/rewrite/flags.html
type RewriteRule struct {
	shared.HeaderList `yaml:",inline"`

	On bool   `yaml:"on" json:"on"` // 是否开启
	Id string `yaml:"id" json:"id"` // ID

	// 开启的条件
	// 语法为：cond param operator value 比如：
	// - cond ${status} gte 200
	// - cond ${arg.name} eq lily
	// - cond ${requestPath} regexp .*\.png
	Cond []*shared.RequestCond `yaml:"cond" json:"cond"`

	// 规则
	// 语法为：pattern regexp 比如：
	// - pattern ^/article/(\d+).html
	Pattern string `yaml:"pattern" json:"pattern"`

	reg *regexp.Regexp

	// 要替换成的URL
	// 支持反向引用：${0}, ${1}, ...，也支持?P<NAME>语法
	// - 如果以 proxy:// 开头，表示目标为代理，首先会尝试作为代理ID请求，如果找不到，会尝试作为代理Host请求
	Replace string `yaml:"replace" json:"replace"`

	// 选项
	// @TODO 使用具体的语义化的字段代替
	Flags       []string `yaml:"flags" json:"flags"`
	FlagOptions maps.Map `yaml:"flagOptions" json:"flagOptions"` // flag => options map

	IsBreak     bool   `yaml:"isBreak" json:"isBreak"`         // 终止向下解析
	IsPermanent bool   `yaml:"isPermanent" json:"isPermanent"` // 是否持久性跳转
	ProxyHost   string `yaml:"host" json:"host"`               // 代理模式下重写后的Host

	targetType  int // RewriteTarget*
	targetURL   string
	targetProxy string
}

// 获取新对象
func NewRewriteRule() *RewriteRule {
	return &RewriteRule{
		On:          true,
		Id:          rands.HexString(16),
		FlagOptions: maps.Map{},
	}
}

// 校验
func (this *RewriteRule) Validate() error {
	reg, err := regexp.Compile(this.Pattern)
	if err != nil {
		return err
	}
	this.reg = reg

	// 替换replace中的反向引用
	if strings.HasPrefix(this.Replace, "proxy://") {
		this.targetType = RewriteTargetProxy
		url := this.Replace[len("proxy://"):]
		index := strings.Index(url, "/")
		if index >= 0 {
			this.targetProxy = url[:index]
			this.targetURL = url[index:]
		}
	} else {
		this.targetType = RewriteTargetURL
		this.targetURL = this.Replace
	}

	// 校验条件
	for _, cond := range this.Cond {
		err := cond.Validate()
		if err != nil {
			return err
		}
	}

	// 校验Header
	err = this.ValidateHeaders()
	if err != nil {
		return err
	}

	return nil
}

// 对某个请求执行规则
func (this *RewriteRule) Match(requestPath string, formatter func(source string) string) (replace string, varMapping map[string]string, matched bool) {
	if this.reg == nil {
		return "", nil, false
	}

	matches := this.reg.FindStringSubmatch(requestPath)
	if len(matches) == 0 {
		return "", nil, false
	}

	// 判断条件
	if len(this.Cond) > 0 {
		for _, cond := range this.Cond {
			if !cond.Match(formatter) {
				return "", nil, false
			}
		}
	}

	varMapping = map[string]string{}
	subNames := this.reg.SubexpNames()
	for index, match := range matches {
		varMapping[fmt.Sprintf("%d", index)] = match
		subName := subNames[index]
		if len(subName) > 0 {
			varMapping[subName] = match
		}
	}

	replace = teautils.ParseVariables(this.targetURL, func(varName string) string {
		v, ok := varMapping[varName]
		if ok {
			return v
		}
		return "${" + varName + "}"
	})

	replace = formatter(replace)

	return replace, varMapping, true
}

// 获取目标类型
func (this *RewriteRule) TargetType() int {
	return this.targetType
}

// 获取目标类型
func (this *RewriteRule) TargetProxy() string {
	return this.targetProxy
}

// 获取目标URL
func (this *RewriteRule) TargetURL() string {
	return this.targetURL
}

// 判断是否是外部URL
func (this *RewriteRule) IsExternalURL(url string) bool {
	return shared.RegexpExternalURL.MatchString(url)
}

// 添加Flag
func (this *RewriteRule) AddFlag(flag string, options maps.Map) {
	this.Flags = append(this.Flags, flag)
	if options != nil {
		this.FlagOptions[flag] = options
	}
}

// 重置模式
func (this *RewriteRule) ResetFlags() {
	this.Flags = []string{}
	this.FlagOptions = maps.Map{}
}

// 跳转模式
func (this *RewriteRule) RedirectMode() string {
	if lists.ContainsString(this.Flags, RewriteFlagProxy) {
		return RewriteFlagProxy
	}
	if lists.ContainsString(this.Flags, RewriteFlagRedirect) {
		return RewriteFlagRedirect
	}
	return RewriteFlagProxy
}

// 添加过滤条件
func (this *RewriteRule) AddCond(cond *shared.RequestCond) {
	this.Cond = append(this.Cond, cond)
}

// 是否在引用某个代理
func (this *RewriteRule) RefersProxy(proxyId string) bool {
	return strings.HasPrefix(this.Replace, "proxy://"+proxyId)
}
