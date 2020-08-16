package tunnel

import (
	"github.com/iwind/TeaGo/actions"
	stringutil "github.com/iwind/TeaGo/utils/string"
)

type GenerateSecretAction actions.Action

// 生成密钥
func (this *GenerateSecretAction) RunPost(params struct{}) {
	this.Data["secret"] = stringutil.Rand(32)
	this.Success()
}
