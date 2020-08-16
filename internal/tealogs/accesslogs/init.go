package accesslogs

import (
	"github.com/TeaWeb/build/internal/teautils"
	"github.com/TeaWeb/uaparser"
	"github.com/iwind/TeaGo"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/logs"
	"reflect"
)

var accessLogVars = map[string]string{}
var userAgentParser *uaparser.Parser

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		// 初始化UserAgent分析器
		logs.Println("[proxy]start user-agent parser")
		var err error
		userAgentParser, err = uaparser.NewParser(teautils.WebRoot() + Tea.DS + "resources" + Tea.DS + "regexes.yaml")
		if err != nil {
			logs.Error(err)
		}
	})

	// 初始化访问日志变量
	reflectType := reflect.TypeOf(AccessLog{})
	countField := reflectType.NumField()
	for i := 0; i < countField; i++ {
		field := reflectType.Field(i)
		value := field.Tag.Get("var")
		if len(value) > 0 {
			accessLogVars[value] = field.Name
		}
	}
}
