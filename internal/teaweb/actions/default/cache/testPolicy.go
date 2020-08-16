package cache

import (
	"github.com/TeaWeb/build/internal/teacache"
	"github.com/TeaWeb/build/internal/teaconfigs/shared"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/actionutils"
)

type TestPolicyAction struct {
	actionutils.ParentAction
}

// 测试缓存策略
func (this *TestPolicyAction) Run(params struct {
	Filename string
}) {
	this.SecondMenu("test")

	policy := shared.NewCachePolicyFromFile(params.Filename)
	if policy == nil {
		this.Fail("找不到Policy")
	}

	this.Data["policy"] = policy

	// 类型
	this.Data["typeName"] = teacache.FindTypeName(policy.Type)

	this.Show()
}

// 提交测试
func (this *TestPolicyAction) RunPost(params struct {
	Filename string
	Action   string
	Key      string
	Value    string
}) {
	this.Data["result"] = ""

	policy := shared.NewCachePolicyFromFile(params.Filename)
	if policy == nil {
		this.Data["result"] = "找不到Policy"
		this.Fail()
	}

	manager := teacache.FindCachePolicyManager(params.Filename)
	if manager == nil {
		manager = teacache.NewManagerFromConfig(policy)

		if policy.Type != "memory" {
			defer manager.Close()
		}
	}

	if manager == nil {
		this.Data["result"] = "无法获取对应的缓存管理器"
		this.Fail()
	}

	if params.Action == "write" {
		err := manager.Write(params.Key, []byte(params.Value))
		if err != nil {
			this.Data["result"] = "写入失败：" + err.Error()
			this.Fail()
		}

		this.Data["result"] = "写入成功"
		this.Success()
	} else if params.Action == "read" {
		value, err := manager.Read(params.Key)
		if err != nil {
			if err == teacache.ErrNotFound {
				this.Data["result"] = "找不到对应的值"
				this.Fail()
			} else {
				this.Data["result"] = "读取失败：" + err.Error()
				this.Fail()
			}
		}

		this.Data["result"] = "读取成功，返回值：" + string(value)
		this.Success()
	}
}
