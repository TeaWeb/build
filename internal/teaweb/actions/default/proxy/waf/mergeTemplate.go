package waf

import (
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/iwind/TeaGo/actions"
)

type MergeTemplateAction actions.Action

// 从模版中更新数据
func (this *MergeTemplateAction) RunPost(params struct {
	WafId string
}) {
	wafList := teaconfigs.SharedWAFList()
	waf := wafList.FindWAF(params.WafId)
	if waf == nil {
		this.Fail("找不到WAF")
	}

	waf.MergeTemplate()
	err := wafList.SaveWAF(waf)
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	this.Success()
}
