package shared

import "regexp"

// 常用的正则表达式
var (
	RegexpDigitNumber    = regexp.MustCompile(`^\d+$`)                    // 正整数
	RegexpFloatNumber    = regexp.MustCompile(`^\d+(\.\d+)?$`)            // 正浮点数，不支持e
	RegexpAllDigitNumber = regexp.MustCompile(`^[+-]?\d+$`)               // 整数，支持正负数
	RegexpAllFloatNumber = regexp.MustCompile(`^[+-]?\d+(\.\d+)?$`)       // 浮点数，支持正负数，不支持e
	RegexpExternalURL    = regexp.MustCompile("(?i)^(http|https|ftp)://") // URL
	RegexpNamedVariable  = regexp.MustCompile("\\${[\\w.-]+}")            // 命名变量
)
