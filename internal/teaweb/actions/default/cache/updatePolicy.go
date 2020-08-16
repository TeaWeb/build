package cache

import (
	"fmt"
	"github.com/TeaWeb/build/internal/teacache"
	"github.com/TeaWeb/build/internal/teaconfigs/shared"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/actionutils"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/cache/cacheutils"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/proxy/proxyutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/types"
)

type UpdatePolicyAction struct {
	actionutils.ParentAction
}

// 修改缓存策略
func (this *UpdatePolicyAction) Run(params struct {
	Filename string
}) {
	this.SecondMenu("policy")

	policy := shared.NewCachePolicyFromFile(params.Filename)
	if policy == nil {
		this.Fail("找不到要修改的缓存策略")
	}

	this.Data["types"] = teacache.AllCacheTypes()

	err := policy.Validate()
	if err != nil {
		logs.Error(err)
	}

	this.Data["policy"] = maps.Map{
		"filename":                 policy.Filename,
		"name":                     policy.Name,
		"key":                      policy.Key,
		"type":                     policy.Type,
		"options":                  policy.Options,
		"life":                     policy.Life,
		"status":                   policy.Status,
		"maxSize":                  policy.MaxSize,
		"capacity":                 policy.Capacity,
		"skipSetCookie":            policy.SkipResponseSetCookie,
		"enableRequestCachePragma": policy.EnableRequestCachePragma,
		"cond":                     policy.Cond,
	}

	if len(policy.SkipResponseCacheControlValues) == 0 {
		policy.SkipResponseCacheControlValues = []string{}
	}
	this.Data["skippedCacheControlValues"] = policy.SkipResponseCacheControlValues

	// 匹配条件运算符
	this.Data["condOperators"] = shared.AllRequestOperators()
	this.Data["condVariables"] = proxyutils.DefaultRequestVariables()

	this.Show()
}

// 提交保存
func (this *UpdatePolicyAction) RunPost(params struct {
	Filename string
	Name     string
	Key      string
	Type     string

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
	policy := shared.NewCachePolicyFromFile(params.Filename)
	if policy == nil {
		this.Fail("找不到要修改的缓存策略")
	}

	params.Must.
		Field("name", params.Name).
		Require("请输入策略名称").

		Field("key", params.Key).
		Require("请输入缓存Key")

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

	// 重置缓存策略实例
	teacache.ResetCachePolicyManager(policy.Filename)

	if cacheutils.IsPolicyUsed(params.Filename) {
		proxyutils.NotifyChange()
	}

	this.Success("保存成功")
}
