package cache

import (
	"fmt"
	"github.com/TeaWeb/build/internal/teacache"
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/TeaWeb/build/internal/teaconfigs/shared"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/proxy/proxyutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/types"
)

type CreatePolicyAction actions.Action

// 缓存缓存策略
func (this *CreatePolicyAction) Run(params struct{}) {
	this.Data["types"] = teacache.AllCacheTypes()

	this.Data["skippedCacheControlValues"] = shared.DefaultSkippedResponseCacheControlValues

	// 匹配条件运算符
	this.Data["condOperators"] = shared.AllRequestOperators()
	this.Data["condVariables"] = proxyutils.DefaultRequestVariables()

	this.Show()
}

// 保存提交
func (this *CreatePolicyAction) RunPost(params struct {
	Name string
	Key  string
	Type string

	Capacity                  float64
	CapacityUnit              string
	Life                      int
	LifeUnit                  string
	StatusList                []int
	MaxSize                   float64
	MaxSizeUnit               string
	SkippedCacheControlValues []string
	SkipSetCookie             bool
	EnableRequestCachePragma  bool

	FileDir        string
	FileAutoCreate bool

	RedisNetwork  string
	RedisHost     string
	RedisPort     int
	RedisSock     string
	RedisPassword string

	LeveldbDir string

	Must *actions.Must
}) {
	params.Must.
		Field("name", params.Name).
		Require("请输入策略名称").

		Field("key", params.Key).
		Require("请输入缓存Key")

	policy := shared.NewCachePolicy()
	policy.On = true
	policy.Name = params.Name
	policy.Key = params.Key
	policy.Type = params.Type

	policy.Capacity = fmt.Sprintf("%.2f%s", params.Capacity, params.CapacityUnit)
	policy.Life = fmt.Sprintf("%d%s", params.Life, params.LifeUnit)
	for _, status := range params.StatusList {
		i := types.Int(status)
		if i >= 0 {
			policy.Status = append(policy.Status, i)
		}
	}
	policy.MaxSize = fmt.Sprintf("%.2f%s", params.MaxSize, params.MaxSizeUnit)
	policy.Status = params.StatusList
	policy.SkipResponseCacheControlValues = params.SkippedCacheControlValues
	policy.SkipResponseSetCookie = params.SkipSetCookie
	policy.EnableRequestCachePragma = params.EnableRequestCachePragma

	// 选项
	switch policy.Type {
	case "file":
		params.Must.
			Field("fileDir", params.FileDir).
			Require("请输入缓存存放目录")
		policy.Options = map[string]interface{}{
			"dir":        params.FileDir,
			"autoCreate": params.FileAutoCreate,
		}
	case "redis":
		params.Must.
			Field("redisNetwork", params.RedisNetwork).
			Require("请选择Redis连接协议").
			Field("redisHost", params.RedisHost).
			Require("请输入Redis服务器地址")
		policy.Options = map[string]interface{}{
			"network":  params.RedisNetwork,
			"host":     params.RedisHost,
			"port":     params.RedisPort,
			"password": params.RedisPassword,
			"sock":     params.RedisSock,
		}
	case "leveldb":
		params.Must.
			Field("leveldbDir", params.LeveldbDir).
			Require("请输入数据库存放目录")
		policy.Options = map[string]interface{}{
			"dir": params.LeveldbDir,
		}
	}

	// 匹配条件
	conds, breakCond, err := proxyutils.ParseRequestConds(this.Request, "request")
	if err != nil {
		this.Fail("匹配条件\"" + breakCond.Param + " " + breakCond.Operator + " " + breakCond.Value + "\"校验失败：" + err.Error())
	}
	policy.Cond = conds

	err = policy.Save()
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	config, _ := teaconfigs.SharedCacheConfig()
	config.AddPolicy(policy.Filename)
	err = config.Save()
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	this.Next("/cache", nil)
	this.Success("保存成功")
}
