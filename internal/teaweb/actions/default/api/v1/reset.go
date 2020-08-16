package v1

import (
	"github.com/TeaWeb/build/internal/teaweb/actions/default/api/apiutils"
	"github.com/TeaWeb/build/internal/teaweb/configs"
	"github.com/iwind/TeaGo/actions"
)

type ResetAction actions.Action

// 重置系统配置状态
func (this *ResetAction) RunGet(params struct{}) {
	configs.SharedAdminConfig().Reset()
	apiutils.SuccessOK(this)
}
