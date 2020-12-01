package tunnel

import (
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/rands"
)

type GenerateSecretAction actions.Action

// 生成密钥
func (this *GenerateSecretAction) RunPost(params struct{}) {
	this.Data["secret"] = rands.HexString(16)
	this.Success()
}
