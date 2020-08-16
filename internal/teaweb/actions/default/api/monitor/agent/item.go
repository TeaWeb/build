package agent

import (
	"github.com/TeaWeb/build/internal/teadb"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/api/apiutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/logs"
)

type ItemAction actions.Action

// 监控项
func (this *ItemAction) Run(params struct {
	AgentId string
	ItemId  string
}) {
	apiutils.ValidateUser(this)

	// 获取数据
	v, err := teadb.AgentValueDAO().FindLatestItemValue(params.AgentId, "", params.ItemId)
	if err != nil {
		logs.Error(err)
		apiutils.Fail(this, err.Error())
	}

	if v == nil {
		apiutils.Fail(this, "item value not found")
	}

	apiutils.Success(this, v.Value)
}
