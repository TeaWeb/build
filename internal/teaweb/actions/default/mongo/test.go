package mongo

import (
	"github.com/TeaWeb/build/internal/teadb"
	"github.com/iwind/TeaGo/actions"
)

type TestAction actions.Action

// 测试Mongo连接
func (this *TestAction) Run(params struct{}) {
	err := teadb.SharedDB().Test()
	if err != nil {
		this.Fail()
	} else {
		this.Success()
	}
}
