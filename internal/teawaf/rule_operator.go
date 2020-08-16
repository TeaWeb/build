package teawaf

type RuleOperator = string
type RuleCaseInsensitive = string

const (
	RuleOperatorGt           RuleOperator = "gt"
	RuleOperatorGte          RuleOperator = "gte"
	RuleOperatorLt           RuleOperator = "lt"
	RuleOperatorLte          RuleOperator = "lte"
	RuleOperatorEq           RuleOperator = "eq"
	RuleOperatorNeq          RuleOperator = "neq"
	RuleOperatorEqString     RuleOperator = "eq string"
	RuleOperatorNeqString    RuleOperator = "neq string"
	RuleOperatorMatch        RuleOperator = "match"
	RuleOperatorNotMatch     RuleOperator = "not match"
	RuleOperatorContains     RuleOperator = "contains"
	RuleOperatorNotContains  RuleOperator = "not contains"
	RuleOperatorPrefix       RuleOperator = "prefix"
	RuleOperatorSuffix       RuleOperator = "suffix"
	RuleOperatorHasKey       RuleOperator = "has key" // has key in slice or map
	RuleOperatorVersionGt    RuleOperator = "version gt"
	RuleOperatorVersionLt    RuleOperator = "version lt"
	RuleOperatorVersionRange RuleOperator = "version range"

	// ip
	RuleOperatorEqIP       RuleOperator = "eq ip"
	RuleOperatorGtIP       RuleOperator = "gt ip"
	RuleOperatorGteIP      RuleOperator = "gte ip"
	RuleOperatorLtIP       RuleOperator = "lt ip"
	RuleOperatorLteIP      RuleOperator = "lte ip"
	RuleOperatorIPRange    RuleOperator = "ip range"
	RuleOperatorNotIPRange RuleOperator = "not ip range"
	RuleOperatorIPMod10    RuleOperator = "ip mod 10"
	RuleOperatorIPMod100   RuleOperator = "ip mod 100"
	RuleOperatorIPMod      RuleOperator = "ip mod"

	RuleCaseInsensitiveNone = "none"
	RuleCaseInsensitiveYes  = "yes"
	RuleCaseInsensitiveNo   = "no"
)

type RuleOperatorDefinition struct {
	Name            string
	Code            string
	Description     string
	CaseInsensitive RuleCaseInsensitive // default caseInsensitive setting
}

var AllRuleOperators = []*RuleOperatorDefinition{
	{
		Name:            "数值大于",
		Code:            RuleOperatorGt,
		Description:     "使用数值对比大于",
		CaseInsensitive: RuleCaseInsensitiveNone,
	},
	{
		Name:            "数值大于等于",
		Code:            RuleOperatorGte,
		Description:     "使用数值对比大于等于",
		CaseInsensitive: RuleCaseInsensitiveNone,
	},
	{
		Name:            "数值小于",
		Code:            RuleOperatorLt,
		Description:     "使用数值对比小于",
		CaseInsensitive: RuleCaseInsensitiveNone,
	},
	{
		Name:            "数值小于等于",
		Code:            RuleOperatorLte,
		Description:     "使用数值对比小于等于",
		CaseInsensitive: RuleCaseInsensitiveNone,
	},
	{
		Name:            "数值等于",
		Code:            RuleOperatorEq,
		Description:     "使用数值对比等于",
		CaseInsensitive: RuleCaseInsensitiveNone,
	},
	{
		Name:            "数值不等于",
		Code:            RuleOperatorNeq,
		Description:     "使用数值对比不等于",
		CaseInsensitive: RuleCaseInsensitiveNone,
	},
	{
		Name:            "字符串等于",
		Code:            RuleOperatorEqString,
		Description:     "使用字符串对比等于",
		CaseInsensitive: RuleCaseInsensitiveNo,
	},
	{
		Name:            "字符串不等于",
		Code:            RuleOperatorNeqString,
		Description:     "使用字符串对比不等于",
		CaseInsensitive: RuleCaseInsensitiveNo,
	},
	{
		Name:            "正则匹配",
		Code:            RuleOperatorMatch,
		Description:     "使用正则表达式匹配，在头部使用(?i)表示不区分大小写，<a href=\"http://teaos.cn/doc/regexp/Regexp.md\" target=\"_blank\">正则表达式语法 &raquo;</a>",
		CaseInsensitive: RuleCaseInsensitiveYes,
	},
	{
		Name:            "正则不匹配",
		Code:            RuleOperatorNotMatch,
		Description:     "使用正则表达式不匹配，在头部使用(?i)表示不区分大小写，<a href=\"http://teaos.cn/doc/regexp/Regexp.md\" target=\"_blank\">正则表达式语法 &raquo;</a>",
		CaseInsensitive: RuleCaseInsensitiveYes,
	},
	{
		Name:            "包含字符串",
		Code:            RuleOperatorContains,
		Description:     "包含某个字符串",
		CaseInsensitive: RuleCaseInsensitiveNo,
	},
	{
		Name:            "不包含字符串",
		Code:            RuleOperatorNotContains,
		Description:     "不包含某个字符串",
		CaseInsensitive: RuleCaseInsensitiveNo,
	},
	{
		Name:            "包含前缀",
		Code:            RuleOperatorPrefix,
		Description:     "包含某个前缀",
		CaseInsensitive: RuleCaseInsensitiveNo,
	},
	{
		Name:            "包含后缀",
		Code:            RuleOperatorSuffix,
		Description:     "包含某个后缀",
		CaseInsensitive: RuleCaseInsensitiveNo,
	},
	{
		Name:            "包含索引",
		Code:            RuleOperatorHasKey,
		Description:     "对于一组数据拥有某个键值或者索引",
		CaseInsensitive: RuleCaseInsensitiveNo,
	},
	{
		Name:            "版本号大于",
		Code:            RuleOperatorVersionGt,
		Description:     "对比版本号大于",
		CaseInsensitive: RuleCaseInsensitiveNo,
	},
	{
		Name:            "版本号小于",
		Code:            RuleOperatorVersionLt,
		Description:     "对比版本号小于",
		CaseInsensitive: RuleCaseInsensitiveNo,
	},
	{
		Name:            "版本号范围",
		Code:            RuleOperatorVersionRange,
		Description:     "判断版本号在某个范围内，格式为version1,version2",
		CaseInsensitive: RuleCaseInsensitiveNo,
	},
	{
		Name:            "IP等于",
		Code:            RuleOperatorEqIP,
		Description:     "将参数转换为IP进行对比",
		CaseInsensitive: RuleCaseInsensitiveNo,
	},
	{
		Name:            "IP大于",
		Code:            RuleOperatorGtIP,
		Description:     "将参数转换为IP进行对比",
		CaseInsensitive: RuleCaseInsensitiveNo,
	},
	{
		Name:            "IP大于等于",
		Code:            RuleOperatorGteIP,
		Description:     "将参数转换为IP进行对比",
		CaseInsensitive: RuleCaseInsensitiveNo,
	},
	{
		Name:            "IP小于",
		Code:            RuleOperatorLtIP,
		Description:     "将参数转换为IP进行对比",
		CaseInsensitive: RuleCaseInsensitiveNo,
	},
	{
		Name:            "IP小于等于",
		Code:            RuleOperatorLteIP,
		Description:     "将参数转换为IP进行对比",
		CaseInsensitive: RuleCaseInsensitiveNo,
	},
	{
		Name:            "IP范围",
		Code:            RuleOperatorIPRange,
		Description:     "IP在某个范围之内，范围格式可以是英文逗号分隔的ip1,ip2，或者CIDR格式的ip/bits",
		CaseInsensitive: RuleCaseInsensitiveNo,
	},
	{
		Name:            "不在IP范围",
		Code:            RuleOperatorNotIPRange,
		Description:     "IP不在某个范围之内，范围格式可以是英文逗号分隔的ip1,ip2，或者CIDR格式的ip/bits",
		CaseInsensitive: RuleCaseInsensitiveNo,
	},
	{
		Name:            "IP取模10",
		Code:            RuleOperatorIPMod10,
		Description:     "对IP参数值取模，除数为10，对比值为余数",
		CaseInsensitive: RuleCaseInsensitiveNo,
	},
	{
		Name:            "IP取模100",
		Code:            RuleOperatorIPMod100,
		Description:     "对IP参数值取模，除数为100，对比值为余数",
		CaseInsensitive: RuleCaseInsensitiveNo,
	},
	{
		Name:            "IP取模",
		Code:            RuleOperatorIPMod,
		Description:     "对IP参数值取模，对比值格式为：除数,余数，比如10,1",
		CaseInsensitive: RuleCaseInsensitiveNo,
	},
}
