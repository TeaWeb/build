package fastcgi

import (
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/maps"
)

type DataAction actions.Action

// Fastcgi数据
func (this *DataAction) Run(params struct {
	ServerId   string
	LocationId string
}) {
	server := teaconfigs.NewServerConfigFromId(params.ServerId)
	if server == nil {
		this.Fail("找不到Server")
	}

	fastcgiList, err := server.FindFastcgiList(params.LocationId)
	if err != nil {
		this.Fail(err.Error())
	}
	this.Data["fastcgiList"] = lists.Map(fastcgiList.AllFastcgi(), func(k int, v interface{}) interface{} {
		f := v.(*teaconfigs.FastcgiConfig)
		return maps.Map{
			"id":             f.Id,
			"on":             f.On,
			"index":          f.Index,
			"poolSize":       f.PoolSize,
			"pass":           f.Pass,
			"readTimeout":    f.ReadTimeout,
			"scriptFilename": f.Param("SCRIPT_FILENAME"),
		}
	})

	this.Success()
}
