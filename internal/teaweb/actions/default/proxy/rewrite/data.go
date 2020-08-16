package rewrite

import (
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/maps"
)

type DataAction actions.Action

// 重写规则数据
func (this *DataAction) Run(params struct {
	ServerId   string
	LocationId string
}) {
	server := teaconfigs.NewServerConfigFromId(params.ServerId)
	if server == nil {
		this.Fail("找不到Server")
	}

	rewriteList, err := server.FindRewriteList(params.LocationId)
	if err != nil {
		this.Fail(err.Error())
	}

	servers := teaconfigs.LoadServerConfigsFromDir(Tea.ConfigDir())
	serversMap := map[string]*teaconfigs.ServerConfig{}
	for _, s := range servers {
		serversMap[s.Id] = s
	}

	this.Data["rewriteList"] = lists.Map(rewriteList.AllRewriteRules(), func(k int, v interface{}) interface{} {
		r := v.(*teaconfigs.RewriteRule)

		targetProxyId := r.TargetProxy()
		targetProxyName := ""
		targetProxyFilename := ""
		if len(targetProxyId) > 0 {
			targetServer, ok := serversMap[targetProxyId]
			if !ok {
				targetProxyName = "代理[已失效]"
			} else {
				targetProxyName = "代理[" + targetServer.Description + "]"
				targetProxyFilename = targetServer.Filename
			}
		}

		return maps.Map{
			"on":                  r.On,
			"id":                  r.Id,
			"flags":               r.Flags,
			"pattern":             r.Pattern,
			"replace":             r.Replace,
			"targetProxyName":     targetProxyName,
			"targetProxyFilename": targetProxyFilename,
			"targetURL":           r.TargetURL(),
			"conds":               r.Cond,
			"isBreak":             r.IsBreak,
			"isPermanent":         r.IsPermanent,
			"proxyHost":           r.ProxyHost,
		}
	})

	this.Success()
}
