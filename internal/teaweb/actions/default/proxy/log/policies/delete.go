package policies

import (
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/TeaWeb/build/internal/tealogs"
	"github.com/iwind/TeaGo/actions"
)

type DeleteAction actions.Action

// 删除策略
func (this *DeleteAction) RunPost(params struct {
	PolicyId string
}) {
	if len(params.PolicyId) == 0 {
		this.Fail("请输入策略ID")
	}

	policy := teaconfigs.NewAccessLogStoragePolicyFromId(params.PolicyId)
	if policy == nil {
		this.Fail("找不到要删除的策略")
	}

	// 判断是否正在使用
	if len(FindAllUsingPolicy(policy.Id)) > 0 {
		this.Fail("此策略正在被使用，不能删除，点击“详情”查看使用此策略的项目")
	}

	policyList := teaconfigs.SharedAccessLogStoragePolicyList()
	policyList.RemoveId(params.PolicyId)
	err := policyList.Save()
	if err != nil {
		this.Fail("删除失败：" + err.Error())
	}

	err = policy.Delete()
	if err != nil {
		this.Fail("删除失败：" + err.Error())
	}

	// 重置缓存策略
	tealogs.ResetPolicyStorage(policy.Id)

	this.Success()
}
