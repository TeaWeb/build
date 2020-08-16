package cache

import (
	"fmt"
	"github.com/TeaWeb/build/internal/teacache"
	"github.com/TeaWeb/build/internal/teaconfigs/shared"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/actionutils"
	"github.com/iwind/TeaGo/maps"
)

type StatPolicyAction struct {
	actionutils.ParentAction
}

// 统计
func (this *StatPolicyAction) Run(params struct {
	Filename string
}) {
	this.SecondMenu("stat")

	policy := shared.NewCachePolicyFromFile(params.Filename)
	if policy == nil {
		this.Fail("找不到Policy")
	}

	this.Data["policy"] = policy

	// 类型
	this.Data["typeName"] = teacache.FindTypeName(policy.Type)

	this.Show()
}

// 获取统计数据
func (this *StatPolicyAction) RunPost(params struct {
	Filename string
}) {
	policy := shared.NewCachePolicyFromFile(params.Filename)
	if policy == nil {
		this.Data["result"] = "找不到Policy"
		this.Fail()
	}

	manager := teacache.FindCachePolicyManager(params.Filename)
	if manager == nil {
		manager = teacache.NewManagerFromConfig(policy)
		if manager == nil {
			this.Fail("找不到Policy")
		}
		defer manager.Close()
	}

	size, countKeys, err := manager.Stat()
	if err != nil {
		this.Fail("发生错误：" + err.Error())
	}

	humanSize := ""
	if size < 1024 {
		humanSize = fmt.Sprintf("%dB", size)
	} else if size < 1024*1024 {
		humanSize = fmt.Sprintf("%.2fKB", float64(size)/1024)
	} else if size < 1024*1024*1024 {
		humanSize = fmt.Sprintf("%.2fMB", float64(size)/1024/1024)
	} else {
		humanSize = fmt.Sprintf("%.2fGB", float64(size)/1024/1024/1024)
	}

	avgSize := ""
	if countKeys == 0 {
		avgSize = "0KB"
	} else {
		avgSize = fmt.Sprintf("%.2fKB", float64(size)/1024/float64(countKeys))
	}

	this.Data["result"] = maps.Map{
		"size":      humanSize,
		"countKeys": countKeys,
		"avgSize":   avgSize,
	}
	this.Success()
}
